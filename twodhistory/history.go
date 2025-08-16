package twodhistory

import (
	"channyein/db"
	"encoding/json"
	"net/http"
)

// GetAllHistoryHandler returns all history data ordered by date desc
func GetAllHistoryHandler(w http.ResponseWriter, r *http.Request) {
	dbConn := db.InitDB()
	defer dbConn.Close()

	date := r.URL.Query().Get("date")
	if date != "" {
		var h HistoryJson
		err := dbConn.QueryRow(`SELECT mset, mvalue, mresult, eset, evalue, eresult, ninternet, nmodern, tmodern, tinternet, date FROM history WHERE date = ?`, date).
			Scan(&h.MSet, &h.MValue, &h.MResult, &h.ESet, &h.EValue, &h.EResult, &h.NInternet, &h.NModern, &h.TModern, &h.TInternet, &h.Date)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(h)
			return
		}
		// Not found, return default struct
		def := HistoryJson{
			MSet:      "----.--",
			MValue:    "-----.--",
			MResult:   "--",
			ESet:      "----.--",
			EValue:    "-----.--",
			EResult:   "--",
			NModern:   "---",
			NInternet: "---",
			TModern:   "---",
			TInternet: "---",
			Date:      date,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(def)
		return
	}

	rows, err := dbConn.Query(`SELECT mset, mvalue, mresult, eset, evalue, eresult, ninternet, nmodern, tmodern, tinternet, date FROM history ORDER BY date DESC`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Force results to be an empty slice, not nil
	results := make([]HistoryJson, 0)
	for rows.Next() {
		var h HistoryJson
		if err := rows.Scan(&h.MSet, &h.MValue, &h.MResult, &h.ESet, &h.EValue, &h.EResult, &h.NInternet, &h.NModern, &h.TModern, &h.TInternet, &h.Date); err == nil {
			results = append(results, h)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
