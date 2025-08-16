package holiday

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Holiday struct (if not already defined elsewhere)

// GetAll returns all holiday rows as []Holiday.
func GetAll(db *sql.DB) ([]Holiday, error) {
	rows, err := db.Query(`SELECT name, date FROM holiday ORDER BY date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var holidays []Holiday
	for rows.Next() {
		var h Holiday
		if err := rows.Scan(&h.Name, &h.Date); err == nil {
			holidays = append(holidays, h)
		}
	}
	return holidays, rows.Err()
}

// GetAllHandler returns all holidays as JSON.
func GetAllHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		holidays, err := GetAll(db)
		if err != nil {
			http.Error(w, "Failed to fetch holidays", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(holidays)
	}
}
