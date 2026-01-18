package queue
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/pawan-sharma-12/sentinel/internal/models"
)
// NewKafkaWriter initializes a new Kafka writer instance
func NewKafkaWriter(brokerURL, topic string) *kafka.Writer{
	return &kafka.Writer{
		Addr: kafka.TCP(brokerURL),
		Topic: topic, 
		Balancer : &kafka.LeastBytes{},
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