package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pawan-sharma-12/sentinel/internal/models"
	"github.com/segmentio/kafka-go"
)

// NewKafkaWriter initializes a new Kafka writer instance
func NewKafkaWriter(brokerURL, topic string) *kafka.Writer{
	return &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "raw-events",
		Balancer:     &kafka.LeastBytes{},
		// Add these three lines:
		BatchTimeout: 10 * time.Millisecond, // Send every 10ms instead of 1s
		BatchSize:    100,                  // Or when 100 messages are ready
		Async:        true,                 // Don't wait for Kafka to acknowledge before responding to GIN
	}
}
// PublishEvent sends the event to the Kafka topic
func PublishEvent(ctx context.Context, w *kafka.Writer, event models.Event) error{
	messageBytes, err := json.Marshal(event)
	if(err != nil){
		return fmt.Errorf("failed to marshal event : %w", err)
	}
	err = w.WriteMessages(ctx, kafka.Message {
		Key: []byte(event.ID),
		Value: messageBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to write message to kafka : %w", err);
	}
	return nil
}