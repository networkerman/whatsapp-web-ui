package whatsapp

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// Client represents a WhatsApp client
type Client struct {
	client     *whatsmeow.Client
	eventStore EventStore
	qrCode     []byte
	qrMutex    sync.RWMutex
	logger     waLog.Logger
	mu         sync.Mutex
	qrChan     chan string
	connecting bool
	store      *MessageStore
}

type EventStore interface {
	StoreMessage(id, chatJID, sender, content string, timestamp time.Time, isFromMe bool) error
	StoreChat(jid, name string, lastMessageTime time.Time) error
}

type ConnectionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// NewClient creates a new WhatsApp client
func NewClient(storePath string, eventStore EventStore) (*Client, error) {
	logger := waLog.Stdout("Client", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", fmt.Sprintf("file:%s/device.db?_foreign_keys=on", storePath), logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create device store: %v", err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	client := whatsmeow.NewClient(deviceStore, logger)
	wa := &Client{
		client:     client,
		eventStore: eventStore,
		logger:     logger,
		store:      NewMessageStore(),
	}

	client.AddEventHandler(wa.handleEvent)
	return wa, nil
}

// Start starts the WhatsApp client
func (c *Client) Start() error {
	if c.client.Store.ID == nil {
		qrChan, _ := c.client.GetQRChannel(context.Background())
		err := c.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}

		select {
		case evt := <-qrChan:
			if evt.Event == "code" {
				png, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
				if err != nil {
					return fmt.Errorf("failed to generate QR code: %v", err)
				}

				c.qrMutex.Lock()
				c.qrCode = png
				c.qrMutex.Unlock()
			}
		case <-context.Background().Done():
			return context.Background().Err()
		}
	} else {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
	}
	return nil
}

// Stop stops the WhatsApp client
func (c *Client) Stop() error {
	if c.client != nil {
		c.client.Disconnect()
	}
	return nil
}

// GetChats returns all chats
func (c *Client) GetChats() []Chat {
	return c.store.GetChats()
}

// GetMessages returns all messages for a chat
func (c *Client) GetMessages(chatID string) []Message {
	return c.store.GetMessages(chatID)
}

// StoreMessage stores a message for a chat
func (c *Client) StoreMessage(chatID string, msg Message) {
	c.store.StoreMessage("", chatID, msg.Sender, msg.Content, msg.Time, msg.IsFromMe)
}

// StoreChat stores a chat
func (c *Client) StoreChat(chat Chat) {
	timestamp := time.Unix(0, chat.Timestamp)
	c.store.StoreChat(chat.ID, chat.Name, timestamp)
}

func (c *Client) GetQRCode(ctx context.Context) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client.IsConnected() {
		return nil, fmt.Errorf("already connected")
	}

	// Get URL query parameters from context
	params := ctx.Value("params")
	
	// Check if refresh is requested or we don't have a QR code yet
	freshQR := false
	if params != nil {
		if p, ok := params.(map[string]string); ok {
			freshQR = p["refresh"] == "true"
		}
	}
	
	c.qrMutex.RLock()
	hasExistingQR := c.qrCode != nil
	c.qrMutex.RUnlock()
	
	// If we have an existing QR and refresh wasn't requested, return it
	if hasExistingQR && !freshQR {
		c.qrMutex.RLock()
		qrCode := c.qrCode
		c.qrMutex.RUnlock()
		log.Printf("Returning existing QR code (fresh=%v)", freshQR)
		return qrCode, nil
	}
	
	// If refresh was requested, clear existing QR code
	if freshQR {
		c.qrMutex.Lock()
		c.qrCode = nil
		c.qrMutex.Unlock()
		log.Printf("Cleared existing QR code for refresh")
	}

	// Create a new QR channel if needed
	if c.qrChan == nil {
		log.Printf("Creating new QR channel")
		c.qrChan = make(chan string, 1)
	}

	// Start connection if not already started
	if !c.connecting {
		log.Printf("Starting new WhatsApp connection")
		if err := c.startConnection(); err != nil {
			return nil, fmt.Errorf("failed to start connection: %v", err)
		}
	} else {
		log.Printf("Connection already in progress")
	}

	// Wait for QR code with timeout
	deadline, _ := ctx.Deadline()
	log.Printf("Waiting for QR code with timeout: %v", deadline)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case qrCode := <-c.qrChan:
		if qrCode == "" {
			return nil, fmt.Errorf("empty QR code received")
		}
		log.Printf("Received QR code, encoding as image")
		// Generate QR code image with higher quality for better scanning
		qr, err := qrcode.Encode(qrCode, qrcode.High, 256)
		if err != nil {
			return nil, fmt.Errorf("failed to encode QR code: %v", err)
		}
		c.qrMutex.Lock()
		c.qrCode = qr
		c.qrMutex.Unlock()
		return qr, nil
	}
}

func (c *Client) startConnection() error {
	if c.connecting {
		return nil
	}
	c.connecting = true

	// Clear any existing client
	if c.client != nil {
		c.client.Disconnect()
		c.client = nil
	}

	// Initialize new client
	client := whatsmeow.NewClient(c.client.Store, nil)
	c.client = client

	// Set up event handler
	client.AddEventHandler(c.handleEvent)

	// Connect in background
	go func() {
		err := client.Connect()
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			c.mu.Lock()
			c.connecting = false
			c.mu.Unlock()
		}
	}()

	return nil
}

func (c *Client) SendMessage(ctx context.Context, recipient, message string) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("not connected to WhatsApp")
	}

	recipientJID, err := types.ParseJID(recipient)
	if err != nil {
		return fmt.Errorf("invalid recipient: %v", err)
	}

	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	_, err = c.client.SendMessage(ctx, recipientJID, msg)
	return err
}

func (c *Client) GetStatus() ConnectionStatus {
	log.Printf("GetStatus called - Client state: connected=%v, hasQR=%v",
		c.client != nil && c.client.IsConnected(), c.qrCode != nil)

	if c.client == nil {
		log.Printf("GetStatus: client is nil")
		return ConnectionStatus{
			Status:  "disconnected",
			Message: "WhatsApp client not initialized",
		}
	}

	c.qrMutex.RLock()
	hasQR := c.qrCode != nil
	c.qrMutex.RUnlock()

	if !c.client.IsConnected() {
		if hasQR {
			log.Printf("GetStatus: not connected, has QR code")
			return ConnectionStatus{
				Status:  "waiting_for_qr",
				Message: "Please scan the QR code to connect",
			}
		}
		log.Printf("GetStatus: not connected, no QR code")
		return ConnectionStatus{
			Status:  "disconnected",
			Message: "WhatsApp client not connected",
		}
	}

	log.Printf("GetStatus: client is connected")
	return ConnectionStatus{
		Status: "connected",
	}
}

func (c *Client) handleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.QR:
		log.Printf("QR code event received")
		if len(v.Codes) > 0 {
			log.Printf("Processing QR code [%d codes available]", len(v.Codes))
			c.mu.Lock()
			if c.qrChan != nil {
				// Use the first QR code in the array
				c.qrChan <- v.Codes[0]
			}
			c.mu.Unlock()
		} else {
			log.Printf("QR code event received but no codes available")
		}

	case *events.Connected:
		c.mu.Lock()
		c.connecting = false
		// Clear QR code and channel on successful connection
		c.qrCode = nil
		if c.qrChan != nil {
			close(c.qrChan)
			c.qrChan = nil
		}
		c.mu.Unlock()

	case *events.Disconnected:
		c.mu.Lock()
		c.connecting = false
		c.mu.Unlock()
	}
}
