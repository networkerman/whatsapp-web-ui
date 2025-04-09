package whatsapp

import (
	"sync"
)

// MessageStore stores messages and chats
type MessageStore struct {
	messages map[string][]Message
	chats    map[string]Chat
	mu       sync.RWMutex
}

// NewMessageStore creates a new MessageStore
func NewMessageStore() *MessageStore {
	return &MessageStore{
		messages: make(map[string][]Message),
		chats:    make(map[string]Chat),
	}
}

// StoreMessage stores a message for a chat
func (s *MessageStore) StoreMessage(chatID string, msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[chatID] = append(s.messages[chatID], msg)
}

// GetMessages returns all messages for a chat
func (s *MessageStore) GetMessages(chatID string) []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.messages[chatID]
}

// StoreChat stores a chat
func (s *MessageStore) StoreChat(chat Chat) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chats[chat.ID] = chat
}

// GetChats returns all stored chats
func (s *MessageStore) GetChats() []Chat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	chats := make([]Chat, 0, len(s.chats))
	for _, chat := range s.chats {
		chats = append(chats, chat)
	}
	return chats
}
