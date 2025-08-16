package live

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// JsonData represents the structure of the live data JSON

type sseClient struct {
	id   int
	ch   chan []byte
	done chan struct{}
}

type LiveSSE struct {
	mu       sync.RWMutex
	clients  map[int]*sseClient
	lastID   int
	lastData JsonData
	lastRaw  []byte
}

var liveSSE = &LiveSSE{
	clients: make(map[int]*sseClient),
}

// In-memory storage for the latest live data
var liveData JsonData
var liveDataMu sync.RWMutex

func (s *LiveSSE) AddClient() *sseClient {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastID++
	client := &sseClient{
		id:   s.lastID,
		ch:   make(chan []byte, 8),
		done: make(chan struct{}),
	}
	s.clients[client.id] = client
	return client
}

func (s *LiveSSE) RemoveClient(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.clients[id]; ok {
		close(c.done)
		delete(s.clients, id)
	}
}

type LivePayload struct {
	JsonData
	ClientCount int `json:"client_count"`
}

func (s *LiveSSE) Broadcast(data JsonData) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	payload := LivePayload{
		JsonData:    data,
		ClientCount: len(s.clients),
	}
	raw, _ := json.Marshal(payload)
	for _, c := range s.clients {
		select {
		case c.ch <- raw:
		default:
		}
	}
	s.lastData = data
	s.lastRaw = raw
}

func (s *LiveSSE) BroadcastLast() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Rebuild payload with current client count
	payload := LivePayload{
		JsonData:    s.lastData,
		ClientCount: len(s.clients),
	}
	raw, _ := json.Marshal(payload)
	for _, c := range s.clients {
		select {
		case c.ch <- raw:
		default:
		}
	}
	s.lastRaw = raw
}

// SSEHandler handles Server-Sent Events for live data
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

	client := liveSSE.AddClient()
	defer liveSSE.RemoveClient(client.id)

	pingTicker := time.NewTicker(10 * time.Second)
	defer pingTicker.Stop()
	broadTicker := time.NewTicker(3 * time.Second)
	defer broadTicker.Stop()

	// Send initial data with client count
	liveDataMu.RLock()
	liveSSE.mu.RLock()
	payload := LivePayload{
		JsonData:    liveData,
		ClientCount: len(liveSSE.clients),
	}
	initRaw, _ := json.Marshal(payload)
	liveSSE.mu.RUnlock()
	liveDataMu.RUnlock()
	w.Write([]byte("data: "))
	w.Write(initRaw)
	w.Write([]byte("\n\n"))
	flusher.Flush()

	notify := w.(http.CloseNotifier).CloseNotify()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-notify:
			return
		case msg := <-client.ch:
			w.Write([]byte("data: "))
			w.Write(msg)
			w.Write([]byte("\n\n"))
			flusher.Flush()
		case <-pingTicker.C:
			w.Write([]byte(": ping\n\n"))
			flusher.Flush()
		case <-broadTicker.C:
			liveSSE.BroadcastLast()
		}
	}
}

// StartLiveBroadcaster starts the smart broadcast loop.
func StartLiveBroadcaster() {
	go func() {
		var lastLive string
		var lastBroadcast time.Time
		for {
			time.Sleep(1 * time.Second)
			liveDataMu.RLock()
			currentLive := liveData.Live
			liveDataMu.RUnlock()
			now := time.Now()
			if currentLive != lastLive {
				liveSSE.Broadcast(liveData)
				lastLive = currentLive
				lastBroadcast = now
			} else if now.Sub(lastBroadcast) >= 5*time.Second {
				liveSSE.Broadcast(liveData)
				lastBroadcast = now
			}
		}
	}()
}
