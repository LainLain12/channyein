package chat

import (
	"encoding/json"
	"net/http"

	"channyein/db"
)

// PostChatHandler handles POST requests to receive and store ChatJson
func PostChatHandler(w http.ResponseWriter, r *http.Request) {
	var chat ChatJson
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Ban check
	if chat.ID != "" {
		conn := db.InitDB()
		if conn != nil {
			defer conn.Close()
			_ = EnsureBanTable(conn)
			if banned, _ := IsBanned(conn, chat.ID); banned {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"status":"error","message":"you are ban"}`))
				return
			}
		}
	}

	AddChat(chat)
	chatSSE.BroadcastNewChat(chat)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
