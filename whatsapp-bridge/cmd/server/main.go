package main

import (
	"log"
	"net/http"

	"github.com/networkerman/whatsapp-web-ui/whatsapp-bridge/internal/api"
	"github.com/networkerman/whatsapp-web-ui/whatsapp-bridge/internal/config"
	"github.com/networkerman/whatsapp-web-ui/whatsapp-bridge/internal/store"
	"github.com/networkerman/whatsapp-web-ui/whatsapp-bridge/internal/whatsapp"
)

func main() {
	cfg := config.New()

	messageStore, err := store.NewMessageStore(cfg.StorePath)
	if err != nil {
		log.Fatalf("Failed to create message store: %v", err)
	}
	defer messageStore.Close()

	whatsappClient, err := whatsapp.NewClient(cfg.StorePath, messageStore)
	if err != nil {
		log.Fatalf("Failed to create WhatsApp client: %v", err)
	}
	defer whatsappClient.Close()

	apiConfig := &api.Config{
		AllowedOrigins: cfg.AllowedOrigins,
	}

	handler := api.NewHandler(whatsappClient, messageStore, apiConfig)

	http.HandleFunc("/api/status", handler.HandleStatus)
	http.HandleFunc("/api/qr", handler.HandleQR)
	// Add other route handlers here...

	log.Printf("Starting server on port %s with allowed origins: %v", cfg.ServerPort, cfg.AllowedOrigins)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
