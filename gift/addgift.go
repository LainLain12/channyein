package gift

import (
	"channyein/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// PostGiftHandler handles POST requests to add or update a gift with image and JSON data
func PostGiftHandler(w http.ResponseWriter, r *http.Request) {
	const imgDir = "gift/images"

	// Detect scheme and host for image URL
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	hostPrefix := scheme + "://" + r.Host + "/gift/images/"

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	// Get JSON string
	jsonStr := r.FormValue("jsonstring")
	if jsonStr == "" {
		http.Error(w, "Missing jsonstring", http.StatusBadRequest)
		return
	}

	var gift GiftJson
	if err := json.Unmarshal([]byte(jsonStr), &gift); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get image file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate image filename using current timestamp
	now := time.Now().UnixNano()
	ext := filepath.Ext(handler.Filename)
	imgName := fmt.Sprintf("%d%s", now, ext)
	imgPath := filepath.Join(imgDir, imgName)
	imgLink := imgName             // Save only filename in DB
	imgURL := hostPrefix + imgName // Full image URL for client

	// Ensure directory exists
	os.MkdirAll(imgDir, 0755)

	// Open DB
	dbConn := db.InitDB()
	defer dbConn.Close()

	var oldImgLink sql.NullString
	err = dbConn.QueryRow("SELECT img_link FROM gift WHERE id = ?", gift.ID).Scan(&oldImgLink)
	isUpdate := err == nil

	// Save new image
	out, err := os.Create(imgPath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(out, file)
	out.Close()
	if err != nil {
		http.Error(w, "Failed to write image", http.StatusInternalServerError)
		return
	}

	// Save to DB
	if isUpdate && oldImgLink.Valid && oldImgLink.String != "" {
		// Delete old image file
		oldName := filepath.Base(oldImgLink.String)
		oldPath := filepath.Join(imgDir, oldName)
		os.Remove(oldPath)
		_, err = dbConn.Exec(`UPDATE gift SET name=?, category=?, img_link=? WHERE id=?`,
			gift.Name, gift.Category, imgLink, gift.ID)
	} else if !isUpdate {
		_, err = dbConn.Exec(`INSERT INTO gift (name, category, img_link) VALUES (?, ?, ?)`,
			gift.Name, gift.Category, imgLink)
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","img_link":"` + imgLink + `","img_url":"` + imgURL + `"}`))

}
