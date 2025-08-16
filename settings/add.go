package settings

import (
	"bytes"
	"channyein/db"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

// PostSettingsHandler handles POST (create / partial update) of the single settings row.
// If the row exists: only columns present in JSON payload are updated.
// If it does not exist: an empty row is inserted first, then partial update applied.
func PostSettingsHandler(w http.ResponseWriter, r *http.Request) {
	bodyBuf := new(bytes.Buffer)
	if _, err := bodyBuf.ReadFrom(r.Body); err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}
	if bodyBuf.Len() == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// Parse into generic map to know which keys are present
	var incoming map[string]interface{}
	if err := json.Unmarshal(bodyBuf.Bytes(), &incoming); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Mapping of JSON keys -> DB column names
	fieldMap := map[string]string{
		"iamshow":  "iamshow",
		"version":  "version",
		"updbody":  "updbody",
		"updtitle": "updtitle",
		"iambody":  "iambody",
		"iamtile":  "iamtile", // keep typo for compatibility
		"iamlink":  "iamlink",
		"needads":  "needads",
	}

	// Collect valid updates
	setClauses := []string{}
	args := []interface{}{}
	updatedKeys := []string{}

	for k, v := range incoming {
		col, ok := fieldMap[strings.ToLower(k)]
		if !ok {
			continue // ignore unknown keys
		}
		setClauses = append(setClauses, col+"=?")
		switch val := v.(type) {
		case string:
			args = append(args, val)
		default:
			// Coerce non-strings via JSON marshal (simple)
			b, _ := json.Marshal(val)
			args = append(args, string(b))
		}
		updatedKeys = append(updatedKeys, col)
	}

	if len(setClauses) == 0 {
		http.Error(w, "no valid fields provided", http.StatusBadRequest)
		return
	}

	conn := db.InitDB()
	if conn == nil {
		http.Error(w, "db init error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	if err := ensureSingleRow(conn); err != nil {
		http.Error(w, "db ensure error", http.StatusInternalServerError)
		return
	}

	query := "UPDATE settings SET " + strings.Join(setClauses, ", ")
	_, err := conn.Exec(query, args...)
	if err != nil {
		http.Error(w, "db update error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"updated_keys": updatedKeys,
	})
}

// ensureSingleRow inserts an empty row if table has zero rows.
func ensureSingleRow(dbConn *sql.DB) error {
	var count int
	if err := dbConn.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		_, err := dbConn.Exec(`INSERT INTO settings (iamshow, version, updbody, updtitle, iambody, iamtile, needads) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			0, "", "", "", "", "", 0)
		return err
	}
	return nil
}
