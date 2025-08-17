package gift

import (
	"channyein/db"
	"encoding/json"
	"net/http"
)

// GiftItem matches your gift table JSON form
type GiftItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	ImgLink  string `json:"img_link"`
}

// SliderItem matches slider table JSON form
type SliderItem struct {
	ID          int    `json:"id"`
	ForwardLink string `json:"forwardlink"`
	Link        string `json:"link"`
}

// GiftSliderResponse contains both lists
type GiftSliderResponse struct {
	Gifts   []GiftItem   `json:"gifts"`
	Sliders []SliderItem `json:"sliders"`
}

// GetGiftSliderHandler returns gifts and sliders in one JSON response.
func GetGiftSliderHandler(w http.ResponseWriter, r *http.Request) {
	conn := db.InitDB()
	if conn == nil {
		http.Error(w, "db init error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// load gifts
	gifts := []GiftItem{}
	grows, err := conn.Query(`SELECT id, name, category, img_link FROM gift ORDER BY id DESC`)
	if err == nil {
		defer grows.Close()
		for grows.Next() {
			var g GiftItem
			if err := grows.Scan(&g.ID, &g.Name, &g.Category, &g.ImgLink); err == nil {
				gifts = append(gifts, g)
			}
		}
	}

	// load sliders
	sliders := []SliderItem{}
	srows, err := conn.Query(`SELECT id, forwardlink, link FROM slider ORDER BY id`)
	if err == nil {
		defer srows.Close()
		for srows.Next() {
			var s SliderItem
			if err := srows.Scan(&s.ID, &s.ForwardLink, &s.Link); err == nil {
				sliders = append(sliders, s)
			}
		}
	}

	resp := GiftSliderResponse{
		Gifts:   gifts,
		Sliders: sliders,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
