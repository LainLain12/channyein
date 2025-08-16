package chat

import (
	"channyein/db"
	"database/sql"
	"encoding/json"
	"net/http"
)

// Report represents a user reporting another user.
type Report struct {
	ID       string `json:"id"`        // reporter id
	ReportID string `json:"report_id"` // reported id
}

// EnsureReportTable creates the report table (unique pair).
func EnsureReportTable(conn *sql.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS report (
    id TEXT NOT NULL,
    report_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(id, report_id)
);`
	_, err := conn.Exec(ddl)
	return err
}

// InsertReportIfNew inserts only if (id, report_id) pair not present.
func InsertReportIfNew(conn *sql.DB, r Report) (inserted bool, err error) {
	if r.ID == "" || r.ReportID == "" {
		return false, nil
	}
	res, err := conn.Exec(`INSERT OR IGNORE INTO report (id, report_id) VALUES (?, ?)`, r.ID, r.ReportID)
	if err != nil {
		return false, err
	}
	aff, _ := res.RowsAffected()
	return aff > 0, nil
}

func countReportsAgainst(conn *sql.DB, reportedID string) (int, error) {
	var c int
	err := conn.QueryRow(`SELECT COUNT(1) FROM report WHERE report_id = ?`, reportedID).Scan(&c)
	return c, err
}

// PostReportHandler handles POST /chat/report with JSON { "id": "...", "report_id": "..." }.
func PostReportHandler(w http.ResponseWriter, r *http.Request) {
	var rp Report
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if rp.ID == "" || rp.ReportID == "" {
		http.Error(w, "id and report_id required", http.StatusBadRequest)
		return
	}
	// Optional: prevent self-report
	if rp.ID == rp.ReportID {
		http.Error(w, "cannot report self", http.StatusBadRequest)
		return
	}

	conn := db.InitDB()
	if conn == nil {
		http.Error(w, "db init error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	if err := EnsureReportTable(conn); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	_ = EnsureBanTable(conn)

	inserted, err := InsertReportIfNew(conn, rp)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !inserted {
		// Pair already exists
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "already_reported",
			"message": "already reported",
			"data":    rp,
		})
		return
	}

	// Newly inserted: check total reports against this ReportID
	total, err := countReportsAgainst(conn, rp.ReportID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "inserted",
			"data":   rp,
			"error":  "count_failed",
		})
		return
	}

	banned := false
	if total >= 10 {
		// Ban the reported user (id = ReportID)
		_ = InsertBan(conn, rp.ReportID)
		banned = true
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "inserted",
		"data":          rp,
		"total_reports": total,
		"banned":        banned,
	})
}
