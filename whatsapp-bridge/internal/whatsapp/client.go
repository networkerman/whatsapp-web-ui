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

type Client struct {
	client     *whatsmeow.Client
	eventStore EventStore
	qrCode     []byte
	qrMutex    sync.RWMutex
	logger     waLog.Logger
}

type EventStore interface {
	StoreMessage(id, chatJID, sender, content string, timestamp time.Time, isFromMe bool) error
	StoreChat(jid, name string, lastMessageTime time.Time) error
}

type ConnectionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

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
	}

	client.AddEventHandler(wa.handleEvent)
	return wa, nil
}

func (c *Client) Connect(ctx context.Context) error {
	if c.client.Store.ID == nil {
		qrChan, _ := c.client.GetQRChannel(ctx)
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
		case <-ctx.Done():
			return ctx.Err()
		}
	} else {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
	}
	return nil
}

func (c *Client) GetQRCode(ctx context.Context) ([]byte, error) {
	log.Printf("GetQRCode called - Client state: connected=%v, hasQR=%v",
		c.client.IsConnected(), c.qrCode != nil)

	c.qrMutex.RLock()
	if c.qrCode != nil {
		log.Printf("Returning existing QR code (size: %d bytes)", len(c.qrCode))
		defer c.qrMutex.RUnlock()
		return c.qrCode, nil
	}
	c.qrMutex.RUnlock() // Release read lock before acquiring write lock

	c.qrMutex.Lock()
	defer c.qrMutex.Unlock()

	// Double-check after acquiring write lock
	if c.qrCode != nil {
		log.Printf("QR code was generated while waiting for lock (size: %d bytes)", len(c.qrCode))
		return c.qrCode, nil
	}

	if c.client.IsConnected() {
		log.Printf("Cannot generate QR code - client is already connected")
		return nil, fmt.Errorf("already connected")
	}

	log.Printf("Starting QR code generation process")
	qrChan, err := c.client.GetQRChannel(ctx)
	if err != nil {
		log.Printf("Failed to get QR channel: %v", err)
		return nil, fmt.Errorf("failed to get QR channel: %v", err)
	}

	log.Printf("Attempting to connect client")
	err = c.client.Connect()
	if err != nil {
		log.Printf("Failed to connect for QR generation: %v", err)
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	log.Printf("Waiting for QR code event")
	select {
	case evt := <-qrChan:
		if evt.Event == "code" {
			log.Printf("Received QR code event, generating PNG")
			png, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
			if err != nil {
				log.Printf("Failed to encode QR code: %v", err)
				return nil, fmt.Errorf("failed to generate QR code: %v", err)
			}
			c.qrCode = png
			log.Printf("Successfully generated QR code (size: %d bytes)", len(png))
			return png, nil
		}
		log.Printf("Received unexpected QR event: %s", evt.Event)
		return nil, fmt.Errorf("unexpected QR event: %s", evt.Event)
	case <-ctx.Done():
		log.Printf("Context cancelled while waiting for QR code")
		return nil, ctx.Err()
	case <-time.After(30 * time.Second):
		log.Printf("Timeout waiting for QR code")
		return nil, fmt.Errorf("timeout waiting for QR code")
	}
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

func (c *Client) Close() error {
	if c.client != nil {
		c.client.Disconnect()
	}
	return nil
}

func (c *Client) handleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		log.Printf("WhatsApp client connected")
		c.qrMutex.Lock()
		c.qrCode = nil // Clear QR code on successful connection
		c.qrMutex.Unlock()
	case *events.Disconnected:
		log.Printf("WhatsApp client disconnected")
	case *events.LoggedOut:
		log.Printf("WhatsApp client logged out")
		c.qrMutex.Lock()
		c.qrCode = nil
		c.qrMutex.Unlock()
	case *events.Message:
		if err := c.eventStore.StoreMessage(
			v.Info.ID,
			v.Info.Chat.String(),
			v.Info.Sender.String(),
			v.Message.GetConversation(),
			v.Info.Timestamp,
			false,
		); err != nil {
			c.logger.Errorf("Error storing message: %v", err)
		}

		if err := c.eventStore.StoreChat(
			v.Info.Chat.String(),
			v.Info.PushName,
			v.Info.Timestamp,
		); err != nil {
			c.logger.Errorf("Error storing chat: %v", err)
		}
	}
}
