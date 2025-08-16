package gift

import "database/sql"

// GiftJson represents the gift data structure for JSON and DB
type GiftJson struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	ImgLink  string `json:"img_link"`
}

// CreateGiftTable creates the gift table if it does not exist
func CreateGiftTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS gift (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		category TEXT,
		img_link TEXT
	)`
	_, err := db.Exec(query)
	return err
}
