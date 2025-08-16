package threed

import (
	"channyein/db"
	"database/sql"
	"encoding/json"
	"net/http"
)

// PostThreeDHandler handles POST requests to add or update a ThreeDJson record
func PostThreeDHandler(w http.ResponseWriter, r *http.Request) {
	var data ThreeDJson
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	var exists int
	err := dbConn.QueryRow("SELECT COUNT(*) FROM threed WHERE date = ?", data.Date).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists > 0 {
		_, err = dbConn.Exec("UPDATE threed SET result = ? WHERE date = ?", data.Result, data.Date)
	} else {
		_, err = dbConn.Exec("INSERT INTO threed (date, result) VALUES (?, ?)", data.Date, data.Result)
	}
	if err != nil {
		http.Error(w, "Database write error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
