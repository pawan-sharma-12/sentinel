package main
import (
	"net/http"
	"time"
	"context"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"github.com/pawan-sharma-12/sentinel/internal/models"
	"github.com/pawan-sharma-12/sentinel/internal/queue"
)
func main(){
	kafkaWriter := queue.NewKafkaWriter("localhost:9092", "raw-events")
	defer kafkaWriter.Close()
	r := gin.Default()
	//health check endpoint 
	r.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{"status":"alive"})
	})
	//ingestion endpoint
	r.POST("/v1/alerts", func(c *gin.Context){
		var event models.Event 
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid request payload"})
			return   
		}
		event.ID = uuid.New().String()
		event.CreatedAt = time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel();
		err := queue.PublishEvent(ctx, kafkaWriter, event);
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to queue event "})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"status" : "accepted",
			"event_id" : event.ID,
			"message" : "Event Streamed to Kafka successfully",
		})
	})
	//Start the server on port 8080
	r.Run(":8080")
}