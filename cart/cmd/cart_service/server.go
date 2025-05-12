package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"route256/cart/internal/app"
	"route256/cart/internal/infra/logger"
	"syscall"
	"time"
)

func main() {
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		logger.Fatal("CONFIG_FILE is required to run the application")
	}

	ctx, cancel := context.WithCancel(context.Background())
	app, err := app.NewApp(ctx, configPath)
	if err != nil {
		logger.Fatal("failed to create application", "error", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = app.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("failed to listen and serve", "error", err)
		}
	}()
	<-quit
	cancel()

	logger.Info("Gracefully shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown application", "error", err)
	}
}
