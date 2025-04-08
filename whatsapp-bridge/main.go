package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
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
func (store *MessageStore) GetChats() ([]Chat, error) {
	rows, err := store.db.Query("SELECT jid, name, last_message_time FROM chats ORDER BY last_message_time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		var lastMessageTime time.Time
		err := rows.Scan(&chat.ID, &chat.Name, &lastMessageTime)
		if err != nil {
			return nil, err
		}
		chat.Timestamp = lastMessageTime.Unix()

		// Get last message from database
		var lastMessage string
		err = store.db.QueryRow(`
			SELECT content 
			FROM messages 
			WHERE chat_jid = ? 
			ORDER BY timestamp DESC 
			LIMIT 1
		`, chat.ID).Scan(&lastMessage)
		if err != nil {
			lastMessage = "" // No last message found
		}
		chat.LastMessage = lastMessage

		chats = append(chats, chat)
	}

	return chats, nil
}

// WhatsAppClient represents our WhatsApp client
type WhatsAppClient struct {
	client  *whatsmeow.Client
	store   *MessageStore
	qrCode  []byte
	qrMutex sync.Mutex
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
				wa.qrCode = []byte(evt.Code)
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
		if c.qrCode != nil {
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
		// Add CORS headers for the image
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Cache-Control", "no-store, must-revalidate")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// If already connected, return error
		if client.client.IsConnected() {
			http.Error(w, "Already connected", http.StatusBadRequest)
			return
		}

		client.qrMutex.Lock()
		defer client.qrMutex.Unlock()

		// If we have a stored QR code and not connected, return it
		if client.qrCode != nil && !client.client.IsConnected() {
			w.Header().Set("Content-Type", "image/png")
			w.Write(client.qrCode)
			return
		}

		// Clear any existing QR code
		client.qrCode = nil

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Get new QR channel
		qrChan, _ := client.client.GetQRChannel(ctx)

		// Only try to connect if not already connecting
		if !client.client.IsConnected() {
			err := client.client.Connect()
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to connect: %v", err), http.StatusInternalServerError)
				return
			}
		}

		select {
		case evt := <-qrChan:
			if evt.Event == "code" {
				// Generate QR code image
				png, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
				if err != nil {
					http.Error(w, fmt.Sprintf("Failed to generate QR code: %v", err), http.StatusInternalServerError)
					return
				}

				// Store the QR code
				client.qrCode = png

				// Set headers and send QR code
				w.Header().Set("Content-Type", "image/png")
				w.Write(png)
				return
			}
		case <-ctx.Done():
			http.Error(w, "Timeout waiting for QR code", http.StatusRequestTimeout)
			return
		}
	})

	// Health check endpoint
	router.HandleFunc("/api/chats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !client.client.IsConnected() {
			// Return empty array when not connected
			json.NewEncoder(w).Encode([]Chat{})
			return
		}

		// Get all chats from the database
		chats := make([]Chat, 0)
		rows, err := store.db.Query(`
			SELECT c.jid, c.name, c.last_message_time, m.content 
			FROM chats c 
			LEFT JOIN messages m ON m.chat_jid = c.jid 
			WHERE m.timestamp = c.last_message_time OR m.timestamp IS NULL
			ORDER BY c.last_message_time DESC
		`)
		if err != nil {
			http.Error(w, "Failed to get chats", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var chat Chat
			var lastMessage sql.NullString
			var lastMessageTime sql.NullTime
			err := rows.Scan(&chat.ID, &chat.Name, &lastMessageTime, &lastMessage)
			if err != nil {
				continue
			}

			if lastMessageTime.Valid {
				chat.Timestamp = lastMessageTime.Time.Unix()
			}
			if lastMessage.Valid {
				chat.LastMessage = lastMessage.String
			}

			// Try to get contact info if it's not a group chat
			if !strings.HasSuffix(chat.ID, "@g.us") {
				jid, _ := types.ParseJID(chat.ID)
				if contact, err := client.client.Store.Contacts.GetContact(jid); err == nil {
					if contact.PushName != "" {
						chat.Name = contact.PushName
					} else if contact.BusinessName != "" {
						chat.Name = contact.BusinessName
					}
				}
			}

			chats = append(chats, chat)
		}

		json.NewEncoder(w).Encode(chats)
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
		w.Header().Set("Content-Type", "application/json")
		if !client.client.IsConnected() {
			if client.client.Store.ID == nil {
				// No ID stored, need to login
				json.NewEncoder(w).Encode(map[string]string{
					"status":  "waiting_for_qr",
					"message": "Please scan the QR code to connect",
				})
			} else {
				// Not connected but has ID
				json.NewEncoder(w).Encode(map[string]string{
					"status":  "disconnected",
					"message": "WhatsApp is disconnected. Please wait while we try to reconnect...",
				})
			}
		} else {
			// Connected
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "connected",
				"message": "Connected! Loading your chats...",
			})
		}
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
