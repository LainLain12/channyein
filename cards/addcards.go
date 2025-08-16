package cards

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PostCardImageHandler handles POST requests to upload an image to cards/images/ (raw body, not multipart)
func PostCardImageHandler(w http.ResponseWriter, r *http.Request) {
	const imgDir = "cards/images"

	// Check Content-Type for image
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "Content-Type must be image/*", http.StatusBadRequest)
		return
	}

	// Get file extension from Content-Type
	ext := ""
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	case "image/bmp":
		ext = ".bmp"
	default:
		ext = filepath.Ext(contentType)
		if ext == "" {
			http.Error(w, "Unsupported image type", http.StatusBadRequest)
			return
		}
	}

	// Generate unique filename using timestamp
	now := time.Now().UnixNano()
	imgName := fmt.Sprintf("%d%s", now, ext)
	imgPath := filepath.Join(imgDir, imgName)

	// Ensure directory exists
	os.MkdirAll(imgDir, 0755)

	out, err := os.Create(imgPath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, r.Body)
	if err != nil {
		http.Error(w, "Failed to write image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","filename":"` + imgName + `"}`))
}
