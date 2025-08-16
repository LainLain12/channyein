package settings

import (
	"channyein/db"
	"encoding/json"
	"net/http"
)

// GetSettingsHandler returns the settings as JSON
func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	dbConn := db.InitDB()
	defer dbConn.Close()

	var s SettingsJson
	err := dbConn.QueryRow(`SELECT iamshow, version, updbody, updtitle, iambody, iamtile, needads FROM settings LIMIT 1`).
		Scan(&s.IamShow, &s.Version, &s.UpdBody, &s.UpdTitle, &s.IamBody, &s.IamTitle, &s.NeedAds)
	if err != nil {
		http.Error(w, "No settings found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}
