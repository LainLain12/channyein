package twodhistory

import (
	"channyein/db"
	"database/sql"
	"encoding/json"
	"net/http"
)

// PostHistoryHandler handles POST requests to add or update a HistoryJson record
func PostHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var data HistoryJson
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	var exists int
	err := dbConn.QueryRow("SELECT COUNT(*) FROM history WHERE date = ?", data.Date).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists > 0 {
		_, err = dbConn.Exec(`UPDATE history SET mset=?, mvalue=?, mresult=?, eset=?, evalue=?, eresult=?, ninternet=?, nmodern=?, tmodern=?, tinternet=? WHERE date=?`,
			data.MSet, data.MValue, data.MResult, data.ESet, data.EValue, data.EResult, data.NInternet, data.NModern, data.TModern, data.TInternet, data.Date)
	} else {
		_, err = dbConn.Exec(`INSERT INTO history (mset, mvalue, mresult, eset, evalue, eresult, ninternet, nmodern, tmodern, tinternet, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			data.MSet, data.MValue, data.MResult, data.ESet, data.EValue, data.EResult, data.NInternet, data.NModern, data.TModern, data.TInternet, data.Date)
	}
	if err != nil {
		http.Error(w, "Database write error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
