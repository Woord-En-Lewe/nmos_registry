package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn           *websocket.Conn
	subscriptionID string
	send           chan []byte
}

type Manager struct {
	clients         map[*Client]bool
	subscriptions   map[string]*Client
	mu              sync.RWMutex
	broadcast       chan *BroadcastMessage
	register        chan *Client
	unregister      chan *Client
	subscriptionEng *registry.SubscriptionEngine
}

type BroadcastMessage struct {
	SubscriptionID string
	Data           []byte
}

func NewManager(subscriptionEng *registry.SubscriptionEngine) *Manager {
	return &Manager{
		clients:         make(map[*Client]bool),
		subscriptions:   make(map[string]*Client),
		broadcast:       make(chan *BroadcastMessage),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		subscriptionEng: subscriptionEng,
	}
}

func (m *Manager) SetSubscriptionEngine(eng *registry.SubscriptionEngine) {
	m.subscriptionEng = eng
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			m.mu.Unlock()

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}
			m.mu.Unlock()

		case message := <-m.broadcast:
			m.mu.RLock()
			if client, ok := m.subscriptions[message.SubscriptionID]; ok {
				select {
				case client.send <- message.Data:
				default:
					close(client.send)
					delete(m.subscriptions, message.SubscriptionID)
				}
			}
			m.mu.RUnlock()
		}
	}
}

func (m *Manager) HandleWebSocket(w http.ResponseWriter, r *http.Request, subscriptionID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		conn:           conn,
		subscriptionID: subscriptionID,
		send:           make(chan []byte, 256),
	}

	m.register <- client

	go m.writePump(client)
	go m.readPump(client)
}

func (m *Manager) writePump(client *Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (m *Manager) readPump(client *Client) {
	defer func() {
		m.unregister <- client
		client.conn.Close()
	}()

	client.conn.SetReadLimit(512)
	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

func (m *Manager) Notify(ctx context.Context, resourceType registry.ResourceType, action string, data interface{}) error {
	msg := WebSocketMessage{
		SubscriptionID: "",
		ResourceType:   resourceType,
		Action:         action,
		Data:           data,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	resourcePath := "/" + string(resourceType) + "s"

	for client := range m.clients {
		sub, err := m.subscriptionEng.GetSubscription(ctx, client.subscriptionID)
		if err != nil {
			continue
		}

		if m.matchesSubscription(resourcePath, sub.ResourcePath) {
			select {
			case client.send <- payload:
			default:
				close(client.send)
				delete(m.clients, client)
			}
		}
	}

	return nil
}

func (m *Manager) matchesSubscription(resourcePath, subscriptionPath string) bool {
	if subscriptionPath == resourcePath {
		return true
	}

	if subscriptionPath == "/"+string(registry.ResourceTypeNode)+"s" ||
		subscriptionPath == "/"+string(registry.ResourceTypeDevice)+"s" ||
		subscriptionPath == "/"+string(registry.ResourceTypeSource)+"s" ||
		subscriptionPath == "/"+string(registry.ResourceTypeFlow)+"s" ||
		subscriptionPath == "/"+string(registry.ResourceTypeSender)+"s" ||
		subscriptionPath == "/"+string(registry.ResourceTypeReceiver)+"s" {
		return true
	}

	return false
}

type WebSocketMessage struct {
	SubscriptionID string                `json:"subscription_id,omitempty"`
	ResourceType   registry.ResourceType `json:"resource_type"`
	Action         string                `json:"action"`
	Data           interface{}           `json:"data"`
}
