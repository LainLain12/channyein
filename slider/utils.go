package slider

import (
	"database/sql"
)

// Slider represents a slider entry for JSON and DB.
type Slider struct {
	ID          int    `json:"id"`
	ForwardLink string `json:"forwardlink"`
	Link        string `json:"link"`
}

// EnsureSliderTable creates the slider table if it does not exist.
func EnsureSliderTable(db *sql.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS slider (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    forwardlink TEXT,
    link TEXT
);`
	_, err := db.Exec(ddl)
	return err
}

// InsertSlider inserts a slider row into the slider table.
func InsertSlider(db *sql.DB, s Slider) (int64, error) {
	res, err := db.Exec(`INSERT INTO slider (forwardlink, link) VALUES (?, ?)`, s.ForwardLink, s.Link)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetAllSliders returns all slider rows.
func GetAllSliders(db *sql.DB) ([]Slider, error) {
	rows, err := db.Query(`SELECT id, forwardlink, link FROM slider ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Slider
	for rows.Next() {
		var s Slider
		if err := rows.Scan(&s.ID, &s.ForwardLink, &s.Link); err == nil {
			out = append(out, s)
		}
	}
	return out, rows.Err()
}

// DeleteSlider deletes a slider by id.
func DeleteSlider(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM slider WHERE id = ?`, id)
	return err
}
