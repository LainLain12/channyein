package holiday

import (
	"database/sql"
)

// Holiday represents a holiday entry.
type Holiday struct {
	Name string `json:"name"`
	Date string `json:"date"` // format YYYY-MM
}

// CreateHolidayTable creates the holiday table if it does not exist.
func CreateHolidayTable(db *sql.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS holiday (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    date TEXT NOT NULL UNIQUE
);`
	_, err := db.Exec(ddl)
	return err
}
