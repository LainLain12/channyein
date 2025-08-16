package holiday

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Holiday struct (if not already defined elsewhere)

// UpsertHoliday updates name if date exists, else inserts new row.
func UpsertHoliday(db *sql.DB, h Holiday) error {
	// Try update first
	res, err := db.Exec(`UPDATE holiday SET name = ? WHERE date = ?`, h.Name, h.Date)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		_, err = db.Exec(`INSERT INTO holiday (name, date) VALUES (?, ?)`, h.Name, h.Date)
	}
	return err
}

// PostHolidayHandler handles POST /holiday requests.
func PostHolidayHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var h Holiday
		if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if h.Name == "" || h.Date == "" {
			http.Error(w, "Missing name or date", http.StatusBadRequest)
			return
		}
		if err := UpsertHoliday(db, h); err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(h)
	}
}
