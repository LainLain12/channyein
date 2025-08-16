package chat

type ChatJson struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Profile  string `json:"profile"`
	Message  string `json:"message"`
	CreateAt string `json:"createat"`
}
