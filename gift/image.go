package gift

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ShowGiftImageHandler serves images from gift/images/ with no-cache headers
func ShowGiftImageHandler(w http.ResponseWriter, r *http.Request) {
	const imgDir = "gift/images"

	imgName := strings.TrimPrefix(r.URL.Path, "/gift/images/")
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

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	http.ServeContent(w, r, imgName, fStat(f), f)
}

func fStat(f *os.File) (modTime time.Time) {
	info, err := f.Stat()
	if err == nil {
		return info.ModTime()
	}
	return
}
