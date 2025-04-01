package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Chat struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LastMessage string `json:"lastMessage,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Sender    string `json:"sender"`
}

func main() {
	r := mux.NewRouter()

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("FRONTEND_URL")},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Routes
	r.HandleFunc("/api/chats", getChats).Methods("GET")
	r.HandleFunc("/api/messages/{chatId}", getMessages).Methods("GET")
	r.HandleFunc("/api/messages/{chatId}", sendMessage).Methods("POST")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, c.Handler(r)); err != nil {
		log.Fatal(err)
	}
}

func getChats(w http.ResponseWriter, r *http.Request) {
	chats := []Chat{
		{
			ID:        "1",
			Name:      "Test Chat",
			LastMessage: "Hello!",
			Timestamp: 1709123456,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	messages := []Message{
		{
			ID:        "1",
			Content:   "Hello!",
			Timestamp: 1709123456,
			Sender:    "user",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Echo back the message
	message.ID = "2"
	message.Sender = "bot"
	message.Timestamp = 1709123457

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
} 