package live

// SaveLiveData saves the received JsonData to storage (to be implemented)

import (
	"encoding/json"
	"os"
)

func SaveLiveData(data JsonData) error {
	liveDataMu.Lock()
	changed := data.Live != liveData.Live
	liveData = data
	liveDataMu.Unlock()

	// Write to Live/live.json
	f, err := os.Create("Live/live.json")
	if err == nil {
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		enc.Encode(data)
		f.Close()
	}

	// --- History logic ---
	// now := time.Now().In(time.Local)
	// weekday := now.Weekday()
	// if weekday != time.Saturday && weekday != time.Sunday {
	// 	afterTime, _ := time.Parse("15:04:05", "16:31:00")
	// 	nowTime, _ := time.Parse("15:04:05", now.Format("15:04:05"))
	// 	if nowTime.After(afterTime) || nowTime.Equal(afterTime) {
	// 		dbConn := db.InitDB()
	// 		defer dbConn.Close()
	// 		// Ensure table exists
	// 		History.CreateHistoryTable(dbConn)
	// 		// Check if today's date exists
	// 		var count int
	// 		err := dbConn.QueryRow("SELECT COUNT(*) FROM history WHERE date = ?", data.Date).Scan(&count)
	// 		if err == nil && count == 0 {
	// 			if data.EResult != "--" {
	// 				// Insert into history
	// 				_, _ = dbConn.Exec(`INSERT INTO history (mset, mvalue, mresult, eset, evalue, eresult, ninternet, nmodern, tmodern, tinternet, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	// 					data.MSet, data.MValue, data.MResult, data.ESet, data.EValue, data.EResult, data.NInternet, data.NModern, data.TModern, data.TInternet, data.Date)
	// 			}
	// 		}
	// 	}
	// }

	if changed {
		liveSSE.Broadcast(data)
	}
	return nil
}
