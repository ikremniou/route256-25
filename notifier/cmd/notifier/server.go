package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"route256/notifier/internal/app"
	"syscall"
	"time"
)

func main() {
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		log.Fatal("CONFIG_FILE is required to run the application")
	}

	app, err := app.NewApp(configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create application %w", err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = app.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("failed to listen and serve %w", err))
		}
	}()
	<-quit

	slog.Info("Gracefully shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Errorf("failed to shutdown application %w", err))
	}
}
