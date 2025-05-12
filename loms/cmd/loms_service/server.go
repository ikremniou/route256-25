package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"route256/loms/internal/app"
	"route256/loms/internal/infra/logger"
	"syscall"
	"time"
)

func main() {
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		logger.Fatal("CONFIG_FILE is required to run the application")
	}

	bootContext := context.Background()
	bootContext, cancel := context.WithTimeout(bootContext, 30*time.Second)
	defer cancel()

	app, err := app.NewApp(bootContext, configPath)
	if err != nil {
		logger.Fatal("Failed to create application", "error", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = app.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()
	<-quit

	logger.Info("Received shutdown signal, shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Failed to shutdown application", "error", err)
	}
}
