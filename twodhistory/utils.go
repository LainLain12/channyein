package twodhistory

import (
	"database/sql"
)

// CreateHistoryTable creates the history table if it does not exist
func CreateHistoryTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS history (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        mset TEXT,
        mvalue TEXT,
        mresult TEXT,
        eset TEXT,
        evalue TEXT,
        eresult TEXT,
        ninternet TEXT,
        nmodern TEXT,
        tmodern TEXT,
        tinternet TEXT,
        date TEXT
    )`
	_, err := db.Exec(query)
	return err
}

type HistoryJson struct {
	MSet      string `json:"mset"`
	MValue    string `json:"mvalue"`
	MResult   string `json:"mresult"`
	ESet      string `json:"eset"`
	EValue    string `json:"evalue"`
	EResult   string `json:"eresult"`
	NInternet string `json:"ninternet"`
	NModern   string `json:"nmodern"`
	TModern   string `json:"tmodern"`
	TInternet string `json:"tinternet"`
	Date      string `json:"date"`
}
