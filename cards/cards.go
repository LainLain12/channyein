package cards

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetAllCardImagesHandler returns all image file names grouped by daily and weekly folders
func GetAllCardImagesHandler(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	hostPrefix := scheme + "://" + r.Host + "/cards/images/"

	allowedExt := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
		".webp": true,
		".bmp":  true,
	}

	result := map[string][]string{
		"daily":  {},
		"weekly": {},
	}

	// Helper to collect image links from a subfolder
	collect := func(folder string) []string {
		var links []string
		dir := filepath.Join("cards/images", folder)
		files, err := os.ReadDir(dir)
		if err != nil {
			return links
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if allowedExt[ext] {
				links = append(links, hostPrefix+folder+"/"+file.Name())
			}
		}
		return links
	}

	result["daily"] = collect("daily")
	result["weekly"] = collect("weekly")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
