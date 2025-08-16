package threed

import (
	"channyein/db"
	"encoding/json"
	"net/http"
)

// GetAllThreeDHandler returns all threed data ordered by date desc
func GetAllThreeDHandler(w http.ResponseWriter, r *http.Request) {
	dbConn := db.InitDB()
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT date, result FROM threed ORDER BY date DESC")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := make([]ThreeDJson, 0)
	for rows.Next() {
		var t ThreeDJson
		if err := rows.Scan(&t.Date, &t.Result); err == nil {
			results = append(results, t)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
