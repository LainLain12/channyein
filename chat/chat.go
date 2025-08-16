package chat

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type chatSSEClient struct {
	id   int
	ch   chan []byte
	done chan struct{}
}

type ChatSSE struct {
	mu      sync.RWMutex
	clients map[int]*chatSSEClient
	lastID  int
}

var chatSSE = &ChatSSE{
	clients: make(map[int]*chatSSEClient),
}

// AddClient adds a new SSE client
func (s *ChatSSE) AddClient() *chatSSEClient {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastID++
	client := &chatSSEClient{
		id:   s.lastID,
		ch:   make(chan []byte, 8),
		done: make(chan struct{}),
	}
	s.clients[client.id] = client
	return client
}

// RemoveClient removes an SSE client
func (s *ChatSSE) RemoveClient(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.clients[id]; ok {
		close(c.done)
		delete(s.clients, id)
	}
}

// Broadcast sends data to all clients
func (s *ChatSSE) Broadcast(data []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients {
		select {
		case c.ch <- data:
		default:
		}
	}
}

// BroadcastAllChats sends all chat data to all clients
func (s *ChatSSE) BroadcastAllChats() {
	allChats := GetAllChats()
	payload, _ := json.Marshal(struct {
		Type string     `json:"type"`
		Data []ChatJson `json:"data"`
	}{
		Type: "all",
		Data: allChats,
	})
	s.Broadcast(payload)
}

// BroadcastNewChat sends a new chat message to all clients
func (s *ChatSSE) BroadcastNewChat(chat ChatJson) {
	payload, _ := json.Marshal(struct {
		Type string   `json:"type"`
		Data ChatJson `json:"data"`
	}{
		Type: "new",
		Data: chat,
	})
	s.Broadcast(payload)
}

// SSEHandler handles SSE connections for chat
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := chatSSE.AddClient()
	defer chatSSE.RemoveClient(client.id)

	// Send all chat data immediately on connect
	allChats := GetAllChats()
	payload, _ := json.Marshal(struct {
		Type string     `json:"type"`
		Data []ChatJson `json:"data"`
	}{
		Type: "all",
		Data: allChats,
	})
	w.Write([]byte("data: " + string(payload) + "\n\n"))
	flusher.Flush()

	pingTicker := time.NewTicker(15 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case msg := <-client.ch:
			w.Write([]byte("event: message\ndata: " + string(msg) + "\n\n"))
			flusher.Flush()
		case <-pingTicker.C:
			w.Write([]byte("ping: {}\n\n"))
			flusher.Flush()
		case <-r.Context().Done():
			return
		case <-client.done:
			return
		}
	}
}

// Call this in your PostChatHandler after AddChat(chat):
// chatSSE.BroadcastNewChat(chat)
