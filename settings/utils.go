package settings

import "database/sql"

// SettingsJson represents the settings data structure for JSON and DB
type SettingsJson struct {
	IamShow  string `json:"iamshow"`
	Version  string `json:"version"`
	UpdBody  string `json:"updbody"`
	UpdTitle string `json:"updtitle"`
	IamBody  string `json:"iambody"`
	IamTitle string `json:"iamtile"` // kept existing json tag (typo preserved for compatibility)
	IamLink  string `json:"iamlink"` // new field
	NeedAds  string `json:"needads"`
}

// CreateSettingsTable creates the settings table if it does not exist
func CreateSettingsTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS settings (
        iamshow TEXT,
        version TEXT,
        updbody TEXT,
        updtitle TEXT,
        iambody TEXT,
        iamtile TEXT,
        iamlink TEXT,
        needads TEXT
    )`
	_, err := db.Exec(query)
	return err
}
