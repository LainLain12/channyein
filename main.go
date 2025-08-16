package main

import (
	cards "channyein/cards"
	chat "channyein/chat"
	"channyein/db"
	gift "channyein/gift"
	"channyein/holiday"
	live "channyein/live"
	"channyein/settings"
	"channyein/threed"
	history "channyein/twodhistory"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	//LiveRunner
	// Initialize the live runner
	// Start periodic timer: calls checktimerTick() every 500ms

	// Initialize the database
	db := db.InitDB()
	history.CreateHistoryTable(db)
	gift.CreateGiftTable(db)
	threed.CreateThreeDTable(db)
	settings.CreateSettingsTable(db)
	holiday.CreateHolidayTable(db)

	os.Setenv("TZ", "Asia/Yangon")
	time.Local, _ = time.LoadLocation("Asia/Yangon")
	///Live
	live.StartLiveBroadcaster()
	http.HandleFunc("/postlive", live.PostLiveHandler)
	http.HandleFunc("/live", live.SSEHandler)
	///History
	http.HandleFunc("/history", history.GetAllHistoryHandler)
	http.HandleFunc("/history/add", history.PostHistoryHandler)
	///Chat
	http.HandleFunc("/chat/send", chat.PostChatHandler)
	http.HandleFunc("/chat/sse", chat.SSEHandler)
	//Chat_Ban
	http.HandleFunc("/chat/addban", chat.PostBanHandler)
	http.HandleFunc("/chat/ban", chat.GetBanHandler)
	//Report
	http.HandleFunc("/chat/report", chat.PostReportHandler)
	///Cards
	http.HandleFunc("/cards/images/daily/", cards.ShowCardImageHandler)
	http.HandleFunc("/cards/images/weekly/", cards.ShowCardImageHandler)
	http.HandleFunc("/cards/", cards.GetAllCardImagesHandler)

	///Gift
	http.HandleFunc("/gift/images/", gift.ShowGiftImageHandler)
	http.HandleFunc("/gift", gift.GetAllGiftHandler)
	http.HandleFunc("/gift/add", gift.PostGiftHandler)
	http.HandleFunc("/gift/delete", gift.DeleteGiftHandler) // Assuming deleteGiftHandler is defined in gift package
	///3D
	http.HandleFunc("/threed/add", threed.PostThreeDHandler)
	http.HandleFunc("/threed", threed.GetAllThreeDHandler)
	///settings
	http.HandleFunc("/settings", settings.GetSettingsHandler)
	http.HandleFunc("/settings/add", settings.PostSettingsHandler)
	///Holidays
	http.HandleFunc("/holiday", holiday.GetAllHandler(db))
	http.HandleFunc("/holiday/add", holiday.PostHolidayHandler(db))
	//Holiday

	// Register the handlers
	log.Println("Server started at :8080 (Asia/Yangon time zone)")
	if err := http.ListenAndServe(":2222", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
