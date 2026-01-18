package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"github.com/pawan-sharma-12/sentinel/internal/cache"
	"github.com/pawan-sharma-12/sentinel/internal/models"
)

func main() {
	// 1. Initialize Kafka Reader
	// GroupID "notification-group" allows you to scale by running multiple workers
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "raw-events",
		GroupID: "notification-group",
	})
	defer reader.Close()

	// 2. Initialize Redis Client for Rate Limiting
	redisClient := cache.NewRedisClient("localhost:6379")

	log.Println("🚀 Sentinel Worker started. Listening for events...")

	// 3. Graceful Shutdown handling (Senior Engineer best practice)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that we can cancel when the program exits
	ctx, cancel := context.WithCancel(context.Background())
	
	go func() {
		<-sigChan
		log.Println("Shutting down worker gracefully...")
		cancel()
	}()

	// 4. Main Processing Loop
	for {
		// Fetch message from Kafka
		// We use the cancellable ctx so the worker stops immediately on Ctrl+C
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			// If context was cancelled, break the loop
			if ctx.Err() != nil {
				break
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		// 5. Unmarshal JSON into our Event Model
		var event models.Event
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		// 6. Check Rate Limit in Redis
		limited, err := redisClient.IsRateLimited(ctx, event.UserID)
		if err != nil {
			log.Printf("Redis error for user %s: %v", event.UserID, err)
			// In a real system, we might allow the alert if Redis is down (Fail Open)
		}

		if limited {
			log.Printf("⚠️  THROTTLED: Dropping alert for User %s (ID: %s) - Limit exceeded", event.UserID, event.ID)
			continue
		}

		// 7. Success - Simulate sending the notification
		fmt.Printf("✅ ALERT SENT | User: %s | Type: %s | ID: %s | Message: %v\n", 
			event.UserID, event.Type, event.ID, event.Payload["message"])
	}

	log.Println("Worker stopped.")
}