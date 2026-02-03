package api

import (
	"cligram/internal/domain"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for now
	},
}

// ClientConnection represents a connected WebSocket client
type ClientConnection struct {
	UserID string
	Conn   *websocket.Conn
}

// WSManager keeps track of connected clients and their chat subscriptions
type WSManager struct {
	Clients     map[string]*ClientConnection            // userID -> client
	ChatClients map[string]map[string]*ClientConnection // chatID -> userID -> client
	Mutex       sync.RWMutex
}

// NewWSManager creates a new manager
func NewWSManager() *WSManager {
	return &WSManager{
		Clients:     make(map[string]*ClientConnection),
		ChatClients: make(map[string]map[string]*ClientConnection),
	}
}

// RegisterClient adds a new client
func (m *WSManager) RegisterClient(client *ClientConnection) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Clients[client.UserID] = client
}

// UnregisterClient removes a client from all chats
func (m *WSManager) UnregisterClient(client *ClientConnection) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	delete(m.Clients, client.UserID)
	for _, clients := range m.ChatClients {
		delete(clients, client.UserID)
	}
}

// SubscribeClientToChat subscribes a client to a chat
func (m *WSManager) SubscribeClientToChat(userID, chatID string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	client, ok := m.Clients[userID]
	if !ok {
		return
	}

	if m.ChatClients[chatID] == nil {
		m.ChatClients[chatID] = make(map[string]*ClientConnection)
	}

	m.ChatClients[chatID][userID] = client
}

// UnsubscribeClientFromChat removes a client from a specific chat
func (m *WSManager) UnsubscribeClientFromChat(userID, chatID string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if clients := m.ChatClients[chatID]; clients != nil {
		delete(clients, userID)
		if len(clients) == 0 {
			delete(m.ChatClients, chatID)
		}
	}
}

// BroadcastMessage sends a message to all clients in a chat
func (m *WSManager) BroadcastMessage(msg domain.Message) {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()

	clients := m.ChatClients[msg.ChatID]
	for _, client := range clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Printf("Error sending message to %s: %v", client.UserID, err)
		}
	}
}

func (s *Server) HandleWS(manager *WSManager, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &ClientConnection{
		UserID: userID,
		Conn:   conn,
	}

	manager.RegisterClient(client)
	log.Printf("User %s connected via WebSocket", userID)

	defer func() {
		manager.UnregisterClient(client)
		conn.Close()
		log.Printf("User %s disconnected from WebSocket", userID)
	}()

	for {
		var msg struct {
			Type   string `json:"type"`
			ChatID string `json:"chat_id"`
			Text   string `json:"text"`
		}

		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message from %s: %v", userID, err)
			break
		}

		switch msg.Type {
		case "subscribe":
			manager.SubscribeClientToChat(userID, msg.ChatID)
			log.Printf("User %s subscribed to chat %s", userID, msg.ChatID)

		case "unsubscribe":
			manager.UnsubscribeClientFromChat(userID, msg.ChatID)
			log.Printf("User %s unsubscribed from chat %s", userID, msg.ChatID)

		case "message":
			if err := s.Service.SendMessage(userID, msg.ChatID, msg.Text); err != nil {
				log.Printf("Failed to save message from %s: %v", userID, err)
				continue
			}

			savedMsgs, err := s.Service.ListMessages(userID, msg.ChatID)
			if err != nil || len(savedMsgs) == 0 {
				log.Printf("Failed to retrieve saved message for broadcast: %v", err)
				continue
			}
			lastMsg := savedMsgs[len(savedMsgs)-1]

			// Broadcast the saved message to all subscribed clients
			manager.BroadcastMessage(lastMsg)
			log.Printf("Broadcasted message from %s to chat %s", userID, msg.ChatID)

		default:
			log.Printf("Unknown message type from %s: %s", userID, msg.Type)
		}
	}
}
