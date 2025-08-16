package cards

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ShowCardImageHandler serves images from cards/images/ with no-cache headers
func ShowCardImageHandler(w http.ResponseWriter, r *http.Request) {
	const imgDir = "cards/images"

	// Get image filename from URL path, e.g. /cards/images/filename.png
	imgName := strings.TrimPrefix(r.URL.Path, "/cards/images/")
	if imgName == "" || strings.Contains(imgName, "..") {
		http.Error(w, "Invalid image name", http.StatusBadRequest)
		return
	}

	imgPath := filepath.Join(imgDir, imgName)
	f, err := os.Open(imgPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	// Set no-cache headers
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	http.ServeContent(w, r, imgName, fStat(f), f)
}

// fStat returns the file's mod time or zero time if error
func fStat(f *os.File) (modTime time.Time) {
	info, err := f.Stat()
	if err == nil {
		return info.ModTime()
	}
	return
}
