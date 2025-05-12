package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"route256/cart/internal/app/handlers/checkout_handler"
	"route256/cart/internal/app/handlers/create_cart_item_handler"
	"route256/cart/internal/app/handlers/delete_cart_handler"
	"route256/cart/internal/app/handlers/delete_cart_item_handler"
	"route256/cart/internal/app/handlers/get_cart_items_handler"
	"route256/cart/internal/domain/cart/repository"
	"route256/cart/internal/domain/cart/service"
	"route256/cart/internal/domain/loms"
	"route256/cart/internal/domain/products"
	"route256/cart/internal/infra/cart_config"
	"route256/cart/internal/infra/logger"
	"route256/cart/internal/infra/sre"
	"route256/cart/pkg/route_http/middleware"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type App struct {
	config *cart_config.Config
	server http.Server
}

func NewApp(ctx context.Context, configPath string) (*App, error) {
	configImpl, err := cart_config.LoadCartConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	err = setupTracing(ctx, configImpl)
	if err != nil {
		logger.Fatal("Failed to setup tracing", "error", err)
	}

	app := &App{
		config: configImpl,
	}

	app.server.Handler = bootstrapHandler(ctx, configImpl)

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	logger.Info("Serving carts", "path", fmt.Sprintf("http://%s", listener.Addr().String()))
	return app.server.Serve(listener)
}

func (app *App) Shutdown(context context.Context) error {
	return app.server.Shutdown(context)
}

func bootstrapHandler(ctx context.Context, config *cart_config.Config) http.Handler {
	ctx, span := otel.Tracer("initialize").Start(ctx, "bootstrap")
	defer span.End()

	productClient := products.NewProductsClient(config)
	orderClient := loms.NewOrderClient(config)
	cartRepository := repository.NewCartRepository()
	cartService := service.NewCartService(cartRepository, productClient, orderClient)
	mux := http.NewServeMux()

	go repository.StartCollectingRepositoryStats(ctx, cartRepository)

	mux.Handle("GET /user/{user_id}/cart", get_cart_items_handler.New(cartService))
	mux.Handle("DELETE /user/{user_id}/cart", delete_cart_handler.New(cartService))

	mux.Handle("POST /user/{user_id}/cart/{sku_id}", create_cart_item_handler.New(cartService))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", delete_cart_item_handler.New(cartService))

	mux.Handle("POST /checkout/{user_id}", checkout_handler.New(cartService))

	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	mux.Handle("GET /metrics", promhttp.Handler())
	registerPprofHandlers(mux, config)

	handler := otelhttp.NewHandler(mux, "http_request")
	handler = sre.NewHandler(handler)
	handler = middleware.NewRequestLoggerMiddleware(handler)
	return middleware.NewGlobalRequestErrorMiddleware(handler)
}

func setupTracing(ctx context.Context, config *cart_config.Config) error {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tracingResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("cart"),
		),
	)

	if err != nil {
		return fmt.Errorf("failed to create jaeger resource: %w", err)
	}

	logger.Info("Connecting to tracer", "jaeger.host", config.Jaeger.Host, "jaeger.port", config.Jaeger.Port)
	hostname := fmt.Sprintf("http://%s:%s", config.Jaeger.Host, config.Jaeger.Port)
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(hostname))
	if err != nil {
		return fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(tracingResource),
	)

	otel.SetTracerProvider(traceProvider)

	return nil
}

func registerPprofHandlers(mux *http.ServeMux, _ *cart_config.Config) {
	mux.Handle("GET /debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("GET /debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("GET /debug/pprof/block", pprof.Handler("block"))
	mux.Handle("GET /debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
