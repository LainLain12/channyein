package slider

import (
	"channyein/db"
	"encoding/json"
	"net/http"
)

// GetAllHandler returns all slider rows as JSON.
func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	conn := db.InitDB()
	if conn == nil {
		http.Error(w, "db init error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	_ = EnsureSliderTable(conn)

	sliders, err := GetAllSliders(conn)
	if err != nil {
		http.Error(w, "failed to fetch sliders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sliders)
}
