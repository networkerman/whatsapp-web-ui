package whatsapp

import "time"

// Message represents a chat message
type Message struct {
	Time     time.Time
	Sender   string
	Content  string
	IsFromMe bool
}

// Chat represents a WhatsApp chat
type Chat struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LastMessage string `json:"lastMessage,omitempty"`
	Timestamp   int64  `json:"timestamp"`
}
