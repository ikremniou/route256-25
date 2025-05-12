package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/IBM/sarama"
)

type LogHandler struct{}

func NewLogHandler() sarama.ConsumerGroupHandler {
	return &LogHandler{}
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (l *LogHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// Setup implements sarama.ConsumerGroupHandler.
func (l *LogHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (l *LogHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				slog.Error("claim messages channel closed")
				return nil
			}

			if err := logOrderState(message); err != nil {
				slog.Error("failed to log order state", "err", err)
				return fmt.Errorf("failed to log order state: %w", err)
			}

			session.MarkMessage(message, "notifier")
		case <-session.Context().Done():
			return nil
		}
	}
}

func logOrderState(message *sarama.ConsumerMessage) error {
	var orderStateMessage OrderStateMessage
	if err := json.Unmarshal(message.Value, &orderStateMessage); err != nil {
		return fmt.Errorf("failed to unmarshal order state message: %w", err)
	}

	userIdString := string(message.Key)
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to int64: %w", err)
	}

	slog.Info("Order state changed",
		"topic", message.Topic,
		"partition", message.Partition,
		"offset", message.Offset,
		"at", message.Timestamp,
		"user_id", userId,
		"order_id", orderStateMessage.OrderId,
		"from_status", orderStateMessage.FromStatus,
		"to_status", orderStateMessage.Status,
	)

	return nil
}
