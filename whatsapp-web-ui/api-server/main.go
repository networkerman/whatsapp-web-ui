package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
	log.Printf("Starting API server...")
	log.Printf("FRONTEND_URL: %s", os.Getenv("FRONTEND_URL"))
	log.Printf("PORT: %s", os.Getenv("PORT"))
	log.Printf("MCP_SERVER_PATH: %s", os.Getenv("MCP_SERVER_PATH"))
	log.Printf("All environment variables: %v", os.Environ())

	r := mux.NewRouter()

	// Enable CORS with specific origins
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://messageai.netlify.app",
			"http://localhost:3000",
			os.Getenv("FRONTEND_URL"),
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		Debug:            true,
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Routes
	r.HandleFunc("/", healthCheck).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/chats", getChats).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/messages/{chatId}", getMessages).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/messages/{chatId}", sendMessage).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/test-connection", testConnection).Methods("GET", "OPTIONS")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, c.Handler(r)); err != nil {
		log.Fatal(err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health check request received")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().String(),
	})
}

func getChats(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for /api/chats")
	log.Printf("Request headers: %v", r.Header)
	log.Printf("Request method: %s", r.Method)
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Request origin: %s", r.Header.Get("Origin"))

	// Set cache control headers
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Return mock data for testing
	log.Printf("Returning mock data")
	mockChats := []Chat{
		{ID: "1", Name: "Test Chat 1", LastMessage: "Hello!", Timestamp: time.Now().Unix()},
		{ID: "2", Name: "Test Chat 2", LastMessage: "Hi there!", Timestamp: time.Now().Unix()},
		{ID: "3", Name: "Test Chat 3", LastMessage: "How are you?", Timestamp: time.Now().Unix()},
	}

	if err := json.NewEncoder(w).Encode(mockChats); err != nil {
		log.Printf("Error encoding mock data: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully sent mock data")
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chatId"]
	log.Printf("Received request for /api/messages/%s", chatID)
	log.Printf("Request headers: %v", r.Header)
	log.Printf("Request method: %s", r.Method)
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Request origin: %s", r.Header.Get("Origin"))

	// Get the WhatsApp bridge URL from environment
	mcpServerPath := os.Getenv("MCP_SERVER_PATH")
	if mcpServerPath == "" {
		log.Printf("MCP_SERVER_PATH not set, using mock data")
		// Return mock data as fallback
		mockMessages := []Message{
			{ID: chatID + "-1", Content: "Hello!", Timestamp: time.Now().Unix(), Sender: "user"},
			{ID: chatID + "-2", Content: "Hi there!", Timestamp: time.Now().Unix(), Sender: "bot"},
		}
		json.NewEncoder(w).Encode(mockMessages)
		return
	}

	// Fetch messages from WhatsApp bridge
	resp, err := http.Get(mcpServerPath + "/api/messages/" + chatID)
	if err != nil {
		log.Printf("Error fetching messages from WhatsApp bridge: %v", err)
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("WhatsApp bridge returned status: %d", resp.StatusCode)
		http.Error(w, "Failed to fetch messages", resp.StatusCode)
		return
	}

	// Forward the response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Copy the response body
	io.Copy(w, resp.Body)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chatId"]
	log.Printf("Received request to send message to chat %s", chatID)
	log.Printf("Request headers: %v", r.Header)
	log.Printf("Request method: %s", r.Method)
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Request origin: %s", r.Header.Get("Origin"))

	var message struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		log.Printf("Error decoding message: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Message content: %s", message.Content)

	// Get the WhatsApp bridge URL from environment
	mcpServerPath := os.Getenv("MCP_SERVER_PATH")
	if mcpServerPath == "" {
		log.Printf("MCP_SERVER_PATH not set, using mock response")
		// Return mock response as fallback
		response := map[string]interface{}{
			"success":   true,
			"chatId":    chatID,
			"message":   message.Content,
			"timestamp": time.Now().Unix(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Forward message to WhatsApp bridge
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error encoding message: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(mcpServerPath+"/api/messages/"+chatID, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending message to WhatsApp bridge: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("WhatsApp bridge returned status: %d", resp.StatusCode)
		http.Error(w, "Failed to send message", resp.StatusCode)
		return
	}

	// Forward the response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Copy the response body
	io.Copy(w, resp.Body)
}

func testConnection(w http.ResponseWriter, r *http.Request) {
	log.Printf("Testing connection to WhatsApp bridge...")

	// Get the WhatsApp bridge URL from environment
	mcpServerPath := os.Getenv("MCP_SERVER_PATH")
	log.Printf("MCP_SERVER_PATH: %s", mcpServerPath)

	if mcpServerPath == "" {
		log.Printf("MCP_SERVER_PATH not set")
		http.Error(w, "MCP_SERVER_PATH not set", http.StatusInternalServerError)
		return
	}

	// Try to connect to the WhatsApp bridge
	url := mcpServerPath + "/api/chats"
	log.Printf("Attempting to connect to: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error connecting to WhatsApp bridge: %v", err)
		http.Error(w, "Failed to connect to WhatsApp bridge: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		http.Error(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)
}
