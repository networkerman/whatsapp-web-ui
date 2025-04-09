package main

import (
	"log"
	"net/http"

	"github.com/networkerman/whatsapp-web-ui/whatsapp-bridge/internal/api"
	"github.com/networkerman/whatsapp-web-ui/whatsapp-bridge/internal/config"
	"github.com/networkerman/whatsapp-web-ui/whatsapp-bridge/internal/whatsapp"
)

func main() {
	cfg := config.New()

	messageStore := whatsapp.NewMessageStore()

	whatsappClient, err := whatsapp.NewClient(cfg.StorePath, messageStore)
	if err != nil {
		log.Fatalf("Failed to create WhatsApp client: %v", err)
	}
	defer whatsappClient.Stop()

	// Start the WhatsApp client
	if err := whatsappClient.Start(); err != nil {
		log.Printf("Warning: Failed to start WhatsApp client: %v", err)
		// Continue even if client fails to start, as QR code might need to be generated
	}

	apiConfig := &api.Config{
		AllowedOrigins: cfg.AllowedOrigins,
	}

	handler := api.NewHandler(whatsappClient, messageStore, apiConfig)

	http.HandleFunc("/api/status", handler.HandleStatus)
	http.HandleFunc("/api/qr", handler.HandleQR)
	// Add other route handlers here...

	log.Printf("Starting server on port %s with allowed origins: %v", cfg.ServerPort, cfg.AllowedOrigins)

	// Use DefaultServeMux for better logging
	server := &http.Server{
		Addr:    ":"+cfg.ServerPort,
		Handler: http.DefaultServeMux,
	}

	log.Printf("Server is listening on port %s", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
