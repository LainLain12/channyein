package slider

import (
	"channyein/db"
	"encoding/json"
	"net/http"
	"strconv"
)

// DeleteSliderHandler handles DELETE /slider?id=123 or GET /slider/delete?id=123
func DeleteSliderHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	conn := db.InitDB()
	if conn == nil {
		http.Error(w, "db init error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	if err := EnsureSliderTable(conn); err != nil {
		http.Error(w, "db ensure error", http.StatusInternalServerError)
		return
	}

	if err := DeleteSlider(conn, id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "deleted",
		"id":     id,
	})
}
