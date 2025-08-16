package gift

import (
	"channyein/db"
	"database/sql"
	"encoding/json"
	"net/http"
)

// GiftJson represents the structure of the gift data

// GetAllGiftHandler returns all gift data as JSON
func GetAllGiftHandler(w http.ResponseWriter, r *http.Request) {
	dbConn := db.InitDB()
	defer dbConn.Close()

	// Detect scheme (http/https)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	hostPrefix := scheme + "://" + r.Host + "/gift/images/"

	category := r.URL.Query().Get("category")
	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = dbConn.Query(`SELECT id, name, category, img_link FROM gift WHERE category = ? ORDER BY id DESC`, category)
	} else {
		rows, err = dbConn.Query(`SELECT id, name, category, img_link FROM gift ORDER BY id DESC`)
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := make([]GiftJson, 0)
	for rows.Next() {
		var g GiftJson
		if err := rows.Scan(&g.ID, &g.Name, &g.Category, &g.ImgLink); err == nil {
			// Prepend hostPrefix to img_link if not empty
			if g.ImgLink != "" {
				g.ImgLink = hostPrefix + g.ImgLink
			}
			results = append(results, g)
		}
	}
	json.NewEncoder(w).Encode(results)
}
