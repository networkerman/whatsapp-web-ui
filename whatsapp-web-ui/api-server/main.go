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
	ID          string `json:"id"`
	Name        string `json:"name"`
	LastMessage string `json:"lastMessage,omitempty"`
	Timestamp   int64  `json:"timestamp"`
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
	mcpServerPath := os.Getenv("MCP_SERVER_PATH")
	if mcpServerPath == "" {
		log.Printf("MCP_SERVER_PATH not set")
		http.Error(w, "MCP server path not configured", http.StatusInternalServerError)
		return
	}

	// Call the MCP server to get chats
	resp, err := http.Get(mcpServerPath + "/api/chats")
	if err != nil {
		log.Printf("Error fetching chats: %v", err)
		http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("MCP server returned status: %d", resp.StatusCode)
		http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
		return
	}

	var chats []Chat
	if err := json.NewDecoder(resp.Body).Decode(&chats); err != nil {
		log.Printf("Error decoding chats: %v", err)
		http.Error(w, "Failed to decode chats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chatId"]

	mcpServerPath := os.Getenv("MCP_SERVER_PATH")
	if mcpServerPath == "" {
		log.Printf("MCP_SERVER_PATH not set")
		http.Error(w, "MCP server path not configured", http.StatusInternalServerError)
		return
	}

	// Call the MCP server to get messages
	resp, err := http.Get(mcpServerPath + "/api/messages/" + chatID)
	if err != nil {
		log.Printf("Error fetching messages: %v", err)
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("MCP server returned status: %d", resp.StatusCode)
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		log.Printf("Error decoding messages: %v", err)
		http.Error(w, "Failed to decode messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chatId"]

	var message struct {
		Content string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mcpServerPath := os.Getenv("MCP_SERVER_PATH")
	if mcpServerPath == "" {
		log.Printf("MCP_SERVER_PATH not set")
		http.Error(w, "MCP server path not configured", http.StatusInternalServerError)
		return
	}

	// Call the MCP server to send message
	resp, err := http.Post(mcpServerPath+"/api/messages/"+chatID, "application/json", nil)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("MCP server returned status: %d", resp.StatusCode)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	var sentMessage Message
	if err := json.NewDecoder(resp.Body).Decode(&sentMessage); err != nil {
		log.Printf("Error decoding sent message: %v", err)
		http.Error(w, "Failed to decode sent message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sentMessage)
}
