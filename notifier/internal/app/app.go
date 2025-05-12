package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"route256/notifier/internal/app/handlers"
	"route256/notifier/internal/notifier_config"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type App struct {
	config      *notifier_config.Config
	kafkaConfig *sarama.Config
	cg          sarama.ConsumerGroup
}

func NewApp(configPath string) (*App, error) {
	config, err := notifier_config.LoadNotifierConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read notifier config, %w", err)
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cg, err := connectToConsumerGroup(timeoutCtx, config, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to kafka consumer group, %w", err)
	}

	err = kafkaConfig.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate kafka config, %w", err)
	}

	app := &App{
		config:      config,
		kafkaConfig: kafkaConfig,
		cg:          cg,
	}

	return app, nil
}

func (app *App) ListenAndServe() error {
	slog.Info("Starting consuming", "group", app.config.Kafka.ConsumerGroupID, "brokers", app.config.Kafka.BrokersPath)
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := app.cg.Consume(ctx, []string{app.config.Kafka.OrderTopic}, handlers.NewLogHandler()); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					slog.Info("Consumer group closed")
					return
				}

				slog.Error("Error consuming messages", "error", err)
			}
		}
	}()

	wg.Wait()
	return nil
}

func (app *App) Shutdown(ctx context.Context) error {
	if app.cg == nil {
		return nil
	}

	if err := app.cg.Close(); err != nil {
		return err
	}

	return nil
}

func connectToConsumerGroup(ctx context.Context, config *notifier_config.Config, kafkaConfig *sarama.Config) (sarama.ConsumerGroup, error) {
	var err error

	maxRetries := 30
	backoff := time.Second

	for retries := range maxRetries {
		if retries > 0 {
			log.Printf("Retrying to create consumer group (attempt %d/%d)...", retries+1, maxRetries)
			time.Sleep(backoff * time.Duration(retries))
		}

		cg, err := sarama.NewConsumerGroup([]string{config.Kafka.BrokersPath}, config.Kafka.ConsumerGroupID, kafkaConfig)
		if err == nil {
			return cg, nil
		}

		if ctx.Err() != nil {
			log.Fatalf("context cancelled while creating producer: %v", ctx.Err())
		}

		log.Printf("Failed to create consumer group: %v, path: %v", err, config.Kafka.BrokersPath)
	}

	return nil, fmt.Errorf("failed to create producer after %d attempts: %w", maxRetries, err)
}
