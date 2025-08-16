package chat

import (
	"channyein/db"
	"database/sql"
	"encoding/json"
	"net/http"
)

// BanEntry represents a banned id (user/client)
type BanEntry struct {
	ID string `json:"id"`
}

// EnsureBanTable creates the ban table if not exists.
func EnsureBanTable(conn *sql.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS ban (
    id TEXT PRIMARY KEY
);`
	_, err := conn.Exec(ddl)
	return err
}

// InsertBan inserts a banned id (ignores if exists).
func InsertBan(conn *sql.DB, id string) error {
	if id == "" {
		return nil
	}
	_, err := conn.Exec(`INSERT OR IGNORE INTO ban (id) VALUES (?)`, id)
	return err
}

// DeleteBan removes a banned id.
func DeleteBan(conn *sql.DB, id string) error {
	_, err := conn.Exec(`DELETE FROM ban WHERE id = ?`, id)
	return err
}

// ListBans returns all banned ids.
func ListBans(conn *sql.DB) ([]BanEntry, error) {
	rows, err := conn.Query(`SELECT id FROM ban ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []BanEntry
	for rows.Next() {
		var b BanEntry
		if err := rows.Scan(&b.ID); err == nil {
			list = append(list, b)
		}
	}
	return list, rows.Err()
}

// IsBanned checks if id is banned.
func IsBanned(conn *sql.DB, id string) (bool, error) {
	var c int
	err := conn.QueryRow(`SELECT COUNT(1) FROM ban WHERE id = ?`, id).Scan(&c)
	return c > 0, err
}

// GetBanHandler handles GET /chat/ban -> list all bans.
func GetBanHandler(w http.ResponseWriter, r *http.Request) {
	conn := db.InitDB()
	defer conn.Close()
	_ = EnsureBanTable(conn)
	list, err := ListBans(conn)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

// PostBanHandler handles POST /chat/ban with JSON { "id": "..." }.
func PostBanHandler(w http.ResponseWriter, r *http.Request) {
	conn := db.InitDB()
	defer conn.Close()
	_ = EnsureBanTable(conn)

	var b BanEntry
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil || b.ID == "" {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := InsertBan(conn, b.ID); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(b)
}

// DeleteBanHandler handles DELETE /chat/ban?id=...
func DeleteBanHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	conn := db.InitDB()
	defer conn.Close()
	_ = EnsureBanTable(conn)
	if err := DeleteBan(conn, id); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
