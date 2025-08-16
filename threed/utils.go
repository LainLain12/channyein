package threed

import "database/sql"

// ThreeDJson represents the 3D data structure for JSON and DB
type ThreeDJson struct {
	Date   string `json:"date"`
	Result string `json:"result"`
}

// CreateThreeDTable creates the threed table if it does not exist
func CreateThreeDTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS threed (
		date TEXT PRIMARY KEY,
		result TEXT
	)`
	_, err := db.Exec(query)
	return err
}
