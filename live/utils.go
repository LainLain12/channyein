package live

type JsonData struct {
	MSet       string `json:"mset"`
	MValue     string `json:"mvalue"`
	MResult    string `json:"mresult"`
	ESet       string `json:"eset"`
	EValue     string `json:"evalue"`
	EResult    string `json:"eresult"`
	NModern    string `json:"nmodern"`
	NInternet  string `json:"ninternet"`
	TModern    string `json:"tmodern"`
	TInternet  string `json:"tinternet"`
	UpdateTime string `json:"updatetime"`
	Date       string `json:"date"`
	Status     string `json:"status"`
	Live       string `json:"live"`
}
