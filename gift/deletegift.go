package gift

import (
	"channyein/db"
	"net/http"
	"strconv"
)

// DeleteGiftHandler handles DELETE /gift?id=123 requests
func DeleteGiftHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	_, err = dbConn.Exec(`DELETE FROM gift WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted","id":` + idStr + `}`))
}
