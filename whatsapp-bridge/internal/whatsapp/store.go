package whatsapp

import (
	"sync"
	"time"
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

// StoreMessage implements EventStore.StoreMessage
func (s *MessageStore) StoreMessage(id, chatJID, sender, content string, timestamp time.Time, isFromMe bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	msg := Message{
		Time:     timestamp,
		Sender:   sender,
		Content:  content,
		IsFromMe: isFromMe,
	}
	s.messages[chatJID] = append(s.messages[chatJID], msg)
	return nil
}

// GetMessages returns all messages for a chat
func (s *MessageStore) GetMessages(chatID string) []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.messages[chatID]
}

// StoreChat implements EventStore.StoreChat
func (s *MessageStore) StoreChat(jid, name string, lastMessageTime time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chats[jid] = Chat{
		ID:        jid,
		Name:      name,
		Timestamp: lastMessageTime.UnixNano(),
	}
	return nil
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
