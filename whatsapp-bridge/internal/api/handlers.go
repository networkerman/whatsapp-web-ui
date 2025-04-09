package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/networkerman/whatsapp-web-ui/whatsapp-bridge/internal/whatsapp"
)

type Handler struct {
	client WhatsAppClient
	store  MessageStore
	config *Config
}

type Config struct {
	AllowedOrigins []string
}

type WhatsAppClient interface {
	GetStatus() whatsapp.ConnectionStatus
	GetQRCode(ctx context.Context) ([]byte, error)
	Start() error
	Stop() error
}

type MessageStore interface {
	GetChats() []whatsapp.Chat
	GetMessages(chatID string) []whatsapp.Message
}

func NewHandler(client WhatsAppClient, store MessageStore, config *Config) *Handler {
	return &Handler{
		client: client,
		store:  store,
		config: config,
	}
}

func (h *Handler) setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// Check if origin is from allowed list
	allowed := false
	for _, allowedOrigin := range h.config.AllowedOrigins {
		if origin == allowedOrigin {
			allowed = true
			break
		}
	}

	// If not in allowed list, check if it's a Netlify domain
	if !allowed && (strings.HasSuffix(origin, ".netlify.app") ||
		strings.HasSuffix(origin, "messageai.netlify.app")) {
		allowed = true
	}

	if allowed {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Cache-Control, Pragma")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
	}

	// Always set no-cache headers for dynamic content
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func (h *Handler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("Status request from %s - Method: %s", r.Header.Get("Origin"), r.Method)
	h.setCORSHeaders(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	status := h.client.GetStatus()
	log.Printf("Status response: %+v", status)

	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Failed to encode status response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleQR(w http.ResponseWriter, r *http.Request) {
	log.Printf("QR code request from %s - Method: %s, Accept: %s",
		r.Header.Get("Origin"), r.Method, r.Header.Get("Accept"))
	h.setCORSHeaders(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check if client accepts image/png
	accept := r.Header.Get("Accept")
	if accept != "*/*" && !strings.Contains(accept, "image/png") {
		log.Printf("Client doesn't accept image/png: %s", accept)
		http.Error(w, "Client must accept image/png", http.StatusNotAcceptable)
		return
	}

	// Check connection status first
	status := h.client.GetStatus()
	if status.Status == "connected" {
		log.Printf("QR code request rejected - client already connected")
		http.Error(w, "Already connected", http.StatusBadRequest)
		return
	}

	// Create a context with timeout for QR code generation
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	log.Printf("Requesting QR code from client")
	qrCode, err := h.client.GetQRCode(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Printf("QR code generation timed out")
			http.Error(w, "Timeout waiting for QR code", http.StatusRequestTimeout)
			return
		}
		if err.Error() == "already connected" {
			log.Printf("QR code generation failed - client connected during request")
			http.Error(w, "Already connected", http.StatusBadRequest)
			return
		}
		log.Printf("Failed to get QR code: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get QR code: %v", err), http.StatusInternalServerError)
		return
	}

	if qrCode == nil {
		log.Printf("No QR code available from client")
		http.Error(w, "No QR code available", http.StatusNotFound)
		return
	}

	log.Printf("Successfully got QR code (size: %d bytes)", len(qrCode))
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(qrCode)))
	w.Write(qrCode)
}
