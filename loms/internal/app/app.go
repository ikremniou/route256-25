package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"route256/loms/internal/infra/logger"
	"route256/loms/internal/infra/loms_config"
	"route256/loms/internal/infra/sre"
	"route256/loms/internal/mw"
	orders_v1 "route256/loms/pkg/api/orders/v1"
	stocks_v1 "route256/loms/pkg/api/stocks/v1"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"net/http/pprof"
)

type App struct {
	config     *loms_config.Config
	httpServer *http.Server
	grpcServer *grpc.Server
	deps       *Deps
}

func NewApp(ctx context.Context, configPath string) (*App, error) {
	configImpl, err := loms_config.LoadLomsConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	err = setupSre(ctx, configImpl)
	if err != nil {
		return nil, fmt.Errorf("setup SRE failed: %w", err)
	}

	app := &App{
		config: configImpl,
	}

	bootstrapApp(ctx, app, configImpl)
	return app, nil
}

func (app *App) ListenAndServe() error {
	grpcAddress := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.GrpcPort)
	httpAddress := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.HttpPort)

	grpcListener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return err
	}

	httpListener, err := net.Listen("tcp", httpAddress)
	if err != nil {
		return err
	}

	logger.Info("Starting loms service", "grpc_address", grpcAddress, "http_address", httpAddress)

	go func() {
		if err := app.grpcServer.Serve(grpcListener); err != nil {
			logger.Error("LomsService gRPC server error", "error", err)
		}
	}()

	return app.httpServer.Serve(httpListener)
}

func (app *App) Shutdown(context context.Context) error {
	var wg sync.WaitGroup
	var grpcErr, httpErr, notifierErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcErr = shutdownGrpcServer(context, app.grpcServer)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpErr = app.httpServer.Shutdown(context)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		notifierErr = app.deps.notifier.Close()
	}()

	wg.Wait()

	if grpcErr != nil || httpErr != nil || notifierErr != nil {
		return fmt.Errorf("failed to shutdown servers: grpc: %w, http: %w, notifier %w", grpcErr, httpErr, notifierErr)
	}

	return nil
}

func shutdownGrpcServer(ctx context.Context, grpcServer *grpc.Server) error {
	grpcShutdown := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(grpcShutdown)
	}()

	select {
	case <-ctx.Done():
		grpcServer.Stop()
		return ctx.Err()
	case <-grpcShutdown:
	}

	return nil
}

func bootstrapApp(ctx context.Context, app *App, config *loms_config.Config) {
	ctx, span := otel.GetTracerProvider().Tracer("initialize").Start(ctx, "bootstrap")
	defer span.End()

	otelgrpc.NewServerHandler()

	var grpcServer = grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			sre.GrpcMw,
			mw.Logger,
			mw.Panic,
			mw.Validate,
		))
	reflection.Register(grpcServer)
	app.grpcServer = grpcServer

	app.deps = InitializeDeps(ctx, grpcServer, config)

	bootstrapHttpGateway(app)
}

func bootstrapHttpGateway(app *App) {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	mux := runtime.NewServeMux()

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.GrpcPort)

	err := orders_v1.RegisterOrdersServiceHandlerFromEndpoint(context.Background(), mux, address, options)
	if err != nil {
		logger.Fatal("Failed to register orders gateway", "error", err)
	}

	err = stocks_v1.RegisterStocksServiceHandlerFromEndpoint(context.Background(), mux, address, options)
	if err != nil {
		logger.Fatal("Failed to register stocks gateway", "error", err)
	}

	bootstrapHandler(mux, app)
}

func bootstrapHandler(mux *runtime.ServeMux, app *App) {
	var defaultPromHandler = promhttp.Handler()
	mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		defaultPromHandler.ServeHTTP(w, r)
	})

	mux.HandlePath("GET", "/health", func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
	})

	registerPprofHandlers(mux, app.config)
	handler := enableCors(mux, app.config)
	handler = sre.NewHandler(handler)

	app.httpServer = &http.Server{
		Handler: handler,
	}
}

func registerPprofHandlers(mux *runtime.ServeMux, _ *loms_config.Config) {
	var handlerWrapper = func(handler http.Handler) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			r.URL.Path = r.URL.Path[len("/debug/pprof"):]
			handler.ServeHTTP(w, r)
		}
	}

	mux.HandlePath("GET", "/debug/pprof/heap", handlerWrapper(pprof.Handler("heap")))
	mux.HandlePath("GET", "/debug/pprof/goroutine", handlerWrapper(pprof.Handler("goroutine")))
	mux.HandlePath("GET", "/debug/pprof/block", handlerWrapper(pprof.Handler("block")))
	mux.HandlePath("GET", "/debug/pprof/threadcreate", handlerWrapper(pprof.Handler("threadcreate")))
	mux.HandlePath("GET", "/debug/pprof//", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		pprof.Index(w, r)
	})
	mux.HandlePath("GET", "/debug/pprof/cmdline", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		pprof.Cmdline(w, r)
	})
	mux.HandlePath("GET", "/debug/pprof/profile", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		pprof.Profile(w, r)
	})
	mux.HandlePath("GET", "/debug/pprof/symbol", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		pprof.Symbol(w, r)
	})
	mux.HandlePath("GET", "/debug/pprof/trace", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		pprof.Trace(w, r)
	})
}
