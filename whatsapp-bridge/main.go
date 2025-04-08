package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// Message represents a chat message for our client
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

// Database handler for storing message history
type MessageStore struct {
	db *sql.DB
}

// Initialize message store
func NewMessageStore() (*MessageStore, error) {
	storePath := os.Getenv("STORE_PATH")
	if storePath == "" {
		storePath = "store"
	}

	// Create directory for database if it doesn't exist
	if err := os.MkdirAll(storePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create store directory: %v", err)
	}

	// Open SQLite database for messages
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s/messages.db?_foreign_keys=on", storePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open message database: %v", err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chats (
			jid TEXT PRIMARY KEY,
			name TEXT,
			last_message_time TIMESTAMP
		);
		
		CREATE TABLE IF NOT EXISTS messages (
			id TEXT,
			chat_jid TEXT,
			sender TEXT,
			content TEXT,
			timestamp TIMESTAMP,
			is_from_me BOOLEAN,
			PRIMARY KEY (id, chat_jid),
			FOREIGN KEY (chat_jid) REFERENCES chats(jid)
		);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return &MessageStore{db: db}, nil
}

// Close the database connection
func (store *MessageStore) Close() error {
	return store.db.Close()
}

// Store a chat in the database
func (store *MessageStore) StoreChat(jid, name string, lastMessageTime time.Time) error {
	_, err := store.db.Exec(
		"INSERT OR REPLACE INTO chats (jid, name, last_message_time) VALUES (?, ?, ?)",
		jid, name, lastMessageTime,
	)
	return err
}

// Store a message in the database
func (store *MessageStore) StoreMessage(id, chatJID, sender, content string, timestamp time.Time, isFromMe bool) error {
	// Only store if there's actual content
	if content == "" {
		return nil
	}

	_, err := store.db.Exec(
		"INSERT OR REPLACE INTO messages (id, chat_jid, sender, content, timestamp, is_from_me) VALUES (?, ?, ?, ?, ?, ?)",
		id, chatJID, sender, content, timestamp, isFromMe,
	)
	return err
}

// Get messages from a chat
func (store *MessageStore) GetMessages(chatJID string, limit int) ([]Message, error) {
	rows, err := store.db.Query(
		"SELECT sender, content, timestamp, is_from_me FROM messages WHERE chat_jid = ? ORDER BY timestamp DESC LIMIT ?",
		chatJID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var timestamp time.Time
		err := rows.Scan(&msg.Sender, &msg.Content, &timestamp, &msg.IsFromMe)
		if err != nil {
			return nil, err
		}
		msg.Time = timestamp
		messages = append(messages, msg)
	}

	return messages, nil
}

// Get all chats
func (store *MessageStore) GetChats() (map[string]time.Time, error) {
	rows, err := store.db.Query("SELECT jid, last_message_time FROM chats ORDER BY last_message_time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make(map[string]time.Time)
	for rows.Next() {
		var jid string
		var lastMessageTime time.Time
		err := rows.Scan(&jid, &lastMessageTime)
		if err != nil {
			return nil, err
		}
		chats[jid] = lastMessageTime
	}

	return chats, nil
}

// WhatsAppClient represents our WhatsApp client
type WhatsAppClient struct {
	client *whatsmeow.Client
	store  *MessageStore
	qrCode string
}

// Create a new WhatsApp client
func NewWhatsAppClient(store *MessageStore) (*WhatsAppClient, error) {
	storePath := os.Getenv("STORE_PATH")
	if storePath == "" {
		storePath = "store"
	}

	// Create device store
	container, err := sqlstore.New("sqlite3", fmt.Sprintf("file:%s/device.db?_foreign_keys=on", storePath), waLog.Stdout("Database", "DEBUG", true))
	if err != nil {
		return nil, fmt.Errorf("failed to create device store: %v", err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	// Create WhatsApp client
	client := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "DEBUG", true))

	wa := &WhatsAppClient{
		client: client,
		store:  store,
	}

	// Set up event handlers
	client.AddEventHandler(wa.handleEvent)

	return wa, nil
}

// Handle WhatsApp events
func (wa *WhatsAppClient) handleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Store incoming message
		err := wa.store.StoreMessage(
			v.Info.ID,
			v.Info.Chat.String(),
			v.Info.Sender.String(),
			v.Message.GetConversation(),
			v.Info.Timestamp,
			false,
		)
		if err != nil {
			fmt.Printf("Error storing message: %v\n", err)
		}

		// Update chat info
		err = wa.store.StoreChat(
			v.Info.Chat.String(),
			v.Info.PushName,
			v.Info.Timestamp,
		)
		if err != nil {
			fmt.Printf("Error storing chat: %v\n", err)
		}
	}
}

// Connect to WhatsApp
func (wa *WhatsAppClient) Connect() error {
	if wa.client.Store.ID == nil {
		// No ID stored, need to login
		qrChan, _ := wa.client.GetQRChannel(context.Background())
		err := wa.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Store the QR code
				wa.qrCode = evt.Code
				fmt.Printf("\nQR code received. Please scan it with WhatsApp to log in\n\n")
				fmt.Println("Waiting for connection...")
			}
		}
	} else {
		// Already logged in, just connect
		err := wa.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
	}
	return nil
}

// Send a message
func (wa *WhatsAppClient) SendMessage(recipient string, message string) error {
	if !wa.client.IsConnected() {
		return fmt.Errorf("not connected to WhatsApp")
	}

	recipientJID, err := types.ParseJID(recipient)
	if err != nil {
		return fmt.Errorf("invalid recipient: %v", err)
	}

	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	_, err = wa.client.SendMessage(context.Background(), recipientJID, msg)
	return err
}

// Close the connection
func (wa *WhatsAppClient) Close() error {
	wa.client.Disconnect()
	return nil
}

// SendMessageResponse represents the response for the send message API
type SendMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SendMessageRequest represents the request body for the send message API
type SendMessageRequest struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

type ConnectionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (c *WhatsAppClient) GetStatus() ConnectionStatus {
	if c.client == nil {
		return ConnectionStatus{
			Status:  "disconnected",
			Message: "WhatsApp client not initialized",
		}
	}

	if !c.client.IsConnected() {
		if c.qrCode != "" {
			return ConnectionStatus{
				Status:  "waiting_for_qr",
				Message: "Please scan QR code to connect",
			}
		}
		return ConnectionStatus{
			Status:  "disconnected",
			Message: "WhatsApp client not connected",
		}
	}

	return ConnectionStatus{
		Status: "connected",
	}
}

func main() {
	// Create message store
	store, err := NewMessageStore()
	if err != nil {
		fmt.Printf("Failed to create message store: %v\n", err)
		return
	}
	defer store.Close()

	// Create WhatsApp client
	client, err := NewWhatsAppClient(store)
	if err != nil {
		fmt.Printf("Failed to create WhatsApp client: %v\n", err)
		return
	}
	defer client.Close()

	// Set up REST API server
	router := mux.NewRouter()

	// Add CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://messageai.netlify.app",
			"http://localhost:3000",
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})
	router.Use(corsMiddleware.Handler)

	// QR code endpoint
	router.HandleFunc("/api/qr", func(w http.ResponseWriter, r *http.Request) {
		if client.client.Store.ID != nil {
			http.Error(w, "Already connected to WhatsApp", http.StatusBadRequest)
			return
		}

		if client.qrCode == "" {
			http.Error(w, "No QR code available", http.StatusNotFound)
			return
		}

		// Generate QR code as PNG
		qr, err := qrcode.New(client.qrCode, qrcode.Medium)
		if err != nil {
			http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
			return
		}

		// Convert QR code to PNG
		png, err := qr.PNG(256)
		if err != nil {
			http.Error(w, "Failed to generate QR code image", http.StatusInternalServerError)
			return
		}

		// Serve the QR code
		w.Header().Set("Content-Type", "image/png")
		w.Write(png)
	})

	// Health check endpoint
	router.HandleFunc("/api/chats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !client.client.IsConnected() {
			// Return empty array when not connected
			json.NewEncoder(w).Encode([]Chat{})
			return
		}

		chats, err := store.GetChats()
		if err != nil {
			http.Error(w, "Failed to get chats", http.StatusInternalServerError)
			return
		}

		// Convert map to array of Chat objects
		var chatList []Chat
		for jid, timestamp := range chats {
			// Get chat name from database
			var name string
			err := store.db.QueryRow("SELECT name FROM chats WHERE jid = ?", jid).Scan(&name)
			if err != nil {
				name = jid // Fallback to JID if name not found
			}

			// Get last message from database
			var lastMessage string
			err = store.db.QueryRow(`
				SELECT content 
				FROM messages 
				WHERE chat_jid = ? 
				ORDER BY timestamp DESC 
				LIMIT 1
			`, jid).Scan(&lastMessage)
			if err != nil {
				lastMessage = "" // No last message found
			}

			chatList = append(chatList, Chat{
				ID:          jid,
				Name:        name,
				LastMessage: lastMessage,
				Timestamp:   timestamp.Unix(),
			})
		}

		// Sort chats by timestamp in descending order
		sort.Slice(chatList, func(i, j int) bool {
			return chatList[i].Timestamp > chatList[j].Timestamp
		})

		// Debug logging
		fmt.Printf("Returning chat list: %+v\n", chatList)

		json.NewEncoder(w).Encode(chatList)
	})

	// Messages endpoint
	router.HandleFunc("/api/messages/{chatID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		chatID := vars["chatID"]

		if !client.client.IsConnected() {
			fmt.Println("WhatsApp client not connected, returning mock data")
			// Return mock data if not connected
			mockMessages := []Message{
				{
					Time:     time.Now(),
					Sender:   "1234567890@s.whatsapp.net",
					Content:  "Mock message",
					IsFromMe: false,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockMessages)
			return
		}

		messages, err := store.GetMessages(chatID, 50)
		if err != nil {
			http.Error(w, "Failed to get messages", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	})

	// Send message endpoint
	router.HandleFunc("/api/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		if req.Recipient == "" || req.Message == "" {
			http.Error(w, "Recipient and message are required", http.StatusBadRequest)
			return
		}

		if !client.client.IsConnected() {
			response := SendMessageResponse{
				Success: false,
				Message: "WhatsApp client not connected",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		err := client.SendMessage(req.Recipient, req.Message)
		response := SendMessageResponse{
			Success: err == nil,
			Message: "Message sent successfully",
		}
		if err != nil {
			response.Message = fmt.Sprintf("Error sending message: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Add status endpoint
	router.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		status := client.GetStatus()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	// Start server
	fmt.Println("Starting REST API server on :8080...")
	go http.ListenAndServe(":8080", router)

	// Connect to WhatsApp in a separate goroutine
	go func() {
		fmt.Println("Connecting to WhatsApp...")
		err = client.Connect()
		if err != nil {
			fmt.Printf("Failed to connect to WhatsApp: %v\n", err)
			return
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
}
