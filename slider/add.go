package slider

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

// UpdateSlider updates a slider row by id.
func UpdateSlider(conn *sql.DB, s Slider) error {
	_, err := conn.Exec(`UPDATE slider SET forwardlink = ?, link = ? WHERE id = ?`, s.ForwardLink, s.Link, s.ID)
	return err
}

// PostSliderHandler handles multipart POST to create or update a slider.
// It reads the slider object from form field "jsonstring" (JSON).
// If json contains id>0 -> try update (if row exists), otherwise insert.
func PostSliderHandler(w http.ResponseWriter, r *http.Request) {
	const imgDir = "slider/images"
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	// detect scheme/host for full url
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	hostPrefix := scheme + "://" + r.Host + "/slider/images/"

	// read jsonstring (must contain slider data; id optional)
	var s Slider
	if js := r.FormValue("jsonstring"); js != "" {
		if err := json.Unmarshal([]byte(js), &s); err != nil {
			http.Error(w, "invalid jsonstring", http.StatusBadRequest)
			return
		}
	}

	// handle optional image upload, update s.Link if provided
	file, fh, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		ext := filepath.Ext(fh.Filename)
		imgName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		if err := os.MkdirAll(imgDir, 0755); err != nil {
			http.Error(w, "failed to create image dir", http.StatusInternalServerError)
			return
		}
		dstPath := filepath.Join(imgDir, imgName)
		out, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "failed to save image", http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(out, file); err != nil {
			out.Close()
			http.Error(w, "failed to write image", http.StatusInternalServerError)
			return
		}
		out.Close()
		// set full url in Link
		s.Link = hostPrefix + imgName
	}

	conn := db.InitDB()
	if conn == nil {
		http.Error(w, "db init error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	if err := EnsureSliderTable(conn); err != nil {
		http.Error(w, "db ensure error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// If ID provided in JSON, check if row exists -> update, otherwise insert
	if s.ID > 0 {
		var cnt int
		if err := conn.QueryRow("SELECT COUNT(1) FROM slider WHERE id = ?", s.ID).Scan(&cnt); err != nil {
			http.Error(w, "db query error", http.StatusInternalServerError)
			return
		}
		if cnt > 0 {
			if err := UpdateSlider(conn, s); err != nil {
				http.Error(w, "update failed", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "updated",
				"data":   s,
			})
			return
		}
		// If id provided but not found, fall through to insert (will ignore provided id)
	}

	// insert
	newID, err := InsertSlider(conn, s)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	s.ID = int(newID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "inserted",
		"data":   s,
	})
}
