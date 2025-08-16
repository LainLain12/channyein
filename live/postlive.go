package live

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// PostLiveHandler handles POST requests with live data JSON
func PostLiveHandler(w http.ResponseWriter, r *http.Request) {
	var data JsonData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := SaveLiveData(data); err != nil {
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		return
	}

	// Update in-memory live data and broadcast to SSE clients
	liveDataMu.Lock()
	liveData = data
	liveDataMu.Unlock()
	if changed := data.Live != liveData.Live; changed {
		liveSSE.Broadcast(data)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
