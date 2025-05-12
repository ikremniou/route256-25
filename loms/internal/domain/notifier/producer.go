package notifier

import (
	"context"
	"fmt"
	"route256/loms/internal/infra/logger"
	"route256/loms/internal/infra/loms_config"
	"route256/loms/internal/infra/sre"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
)

type NotifierProducer struct {
	cfg      *loms_config.Config
	producer sarama.SyncProducer
	stopChan chan struct{}
	wg       *sync.WaitGroup

	outbox *OutboxRepository
}

func NewNotifierProducer(ctx context.Context, cfg *loms_config.Config, outbox *OutboxRepository) *NotifierProducer {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true

	err := kafkaConfig.Validate()
	if err != nil {
		logger.Fatal("Failed to validate kafka config", "error", err)
	}

	ctx, span := otel.GetTracerProvider().Tracer("initialize").Start(ctx, "kafka.producer")
	producer, err := connectToProducer(ctx, cfg, kafkaConfig)
	if err != nil {
		logger.Fatal("Failed to connect to kafka producer", "error", err)
	}
	span.End()

	notifierProducer := &NotifierProducer{
		producer: producer,
		cfg:      cfg,
		wg:       &sync.WaitGroup{},
		stopChan: make(chan struct{}),
		outbox:   outbox,
	}

	notifierProducer.runOutboxPoller()

	logger.Info("Notifier Producer is ready to roll", "brokers", cfg.Kafka.Brokers)
	return notifierProducer
}

func (p *NotifierProducer) Close() error {
	logger.Info("Closing notifier producer...")

	close(p.stopChan)
	err := p.producer.Close()
	p.wg.Wait()

	if err != nil {
		return fmt.Errorf("failed to close producer: %v", err)
	}

	return nil
}

func (p *NotifierProducer) runOutboxPoller() {
	ctx := context.Background()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		ticker := time.NewTicker(time.Duration(p.cfg.Kafka.PollMs) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-p.stopChan:
				return
			case <-ticker.C:
				if err := p.processPendingMessages(ctx); err != nil {
					logger.Warn("Failed to process pending outbox messages", "error", err)
				}
			}
		}
	}()
}

func (p *NotifierProducer) processPendingMessages(ctx context.Context) error {
	return p.outbox.ProcessPendingMessagesFn(ctx, 10, func(row []OutboxEntity) error {
		if len(row) == 0 {
			return nil
		}

		messages := make([]*sarama.ProducerMessage, 0, len(row))
		for _, entity := range row {
			message := &sarama.ProducerMessage{
				Topic: entity.Topic,
				Key:   sarama.StringEncoder(entity.Key),
				Value: sarama.ByteEncoder(entity.Payload),
			}

			messages = append(messages, message)
		}

		startTime := time.Now()
		err := p.producer.SendMessages(messages)
		sre.TrackExternalRequest("kafka_send_messages", err, startTime)
		if err != nil {
			return fmt.Errorf("failed to send messages to producer: %w", err)
		}

		logger.Info("Sent messages to kafka producer", "topic", row[0].Topic, "message_count", len(messages))
		return nil
	})
}

func connectToProducer(
	ctx context.Context,
	cfg *loms_config.Config,
	kafkaConfig *sarama.Config,
) (sarama.SyncProducer, error) {
	var err error

	maxRetries := 30
	backoff := time.Second

	for retries := range maxRetries {
		if retries > 0 {
			logger.Info("Retrying to create producer", "attempt", retries+1, "max_retries", maxRetries)
			time.Sleep(backoff * time.Duration(retries))
		}

		producer, err := sarama.NewSyncProducer([]string{cfg.Kafka.Brokers}, kafkaConfig)
		if err == nil {
			return producer, nil
		}

		if ctx.Err() != nil {
			logger.Warn("Context cancelled while creating producer", "error", ctx.Err())

			return nil, fmt.Errorf("context cancelled while creating producer: %w", ctx.Err())
		}

		logger.Info("Failed to create producer", "error", err, "brokers", cfg.Kafka.Brokers)
	}

	return nil, fmt.Errorf("failed to create producer after %d attempts: %w", maxRetries, err)
}
