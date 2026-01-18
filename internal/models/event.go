package models
import "time"

type Event struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	Type string `json:"type"`
	Payload map[string]interface{} `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}
type Notification struct {
	EventID string `json:"event_id"`
	UserID string `json:"user_id"`
	Message string `json:"message"`
	Channel string `json:"channel"` 
	Status string `json:"status"`
}