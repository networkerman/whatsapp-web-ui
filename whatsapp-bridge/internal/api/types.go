package api

// ConnectionStatus represents the current connection state of the WhatsApp client
type ConnectionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
