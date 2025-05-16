package kafka

import (
	"context"
	"encoding/json"
	"log"

	"notify/internal/domain"

	"github.com/Shopify/sarama"
)

// Consumer represents the Kafka consumer implementation
type Consumer struct {
	consumer sarama.ConsumerGroup
	config   ConsumerConfig
	handler  domain.NotificationService
}

// ConsumerConfig holds the Kafka consumer configuration
type ConsumerConfig struct {
	Brokers []string
	GroupID string
	Topics  Topics
}

// Topics holds the Kafka topic configuration
type Topics struct {
	Notifications string
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg ConsumerConfig, handler domain.NotificationService) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}

	return &Consumer{
		consumer: consumer,
		config:   cfg,
		handler:  handler,
	}
}

// Start begins consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	topics := []string{c.config.Topics.Notifications}

	for {
		err := c.consumer.Consume(ctx, topics, &consumerGroupHandler{
			handler: c.handler,
		})
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type consumerGroupHandler struct {
	handler domain.NotificationService
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var notification domain.Notification
		if err := json.Unmarshal(message.Value, &notification); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		if err := h.handler.ProcessNotification(&notification); err != nil {
			log.Printf("Error processing notification: %v", err)
			continue
		}

		session.MarkMessage(message, "")
	}
	return nil
}
