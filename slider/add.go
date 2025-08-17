package slider

import (
	"channyein/db"
	"database/sql"
	"encoding/json"
	"net/http"
)

// UpdateSlider updates a slider row by id.
func UpdateSlider(conn *sql.DB, s Slider) error {
	_, err := conn.Exec(`UPDATE slider SET forwardlink = ?, link = ? WHERE id = ?`, s.ForwardLink, s.Link, s.ID)
	return err
}

// PostSliderHandler handles create or update of slider.
// If JSON contains id (>0) it updates that row; otherwise it inserts a new row.
func PostSliderHandler(w http.ResponseWriter, r *http.Request) {
	var s Slider
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")

	if s.ID > 0 {
		// update path
		if err := UpdateSlider(conn, s); err != nil {
			http.Error(w, "update failed", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "updated",
			"data":   s,
		})
		return
	}

	// insert path
	newID, err := InsertSlider(conn, s)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	s.ID = int(newID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "inserted",
		"data":   s,
	})
}
