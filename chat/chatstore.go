package chat

import "sync"

var (
	chatStore   []ChatJson
	chatStoreMu sync.Mutex
)

// AddChat adds a ChatJson to the store, keeping max 50 items
func AddChat(chat ChatJson) {
	chatStoreMu.Lock()
	defer chatStoreMu.Unlock()
	chatStore = append(chatStore, chat)
	if len(chatStore) > 50 {
		chatStore = chatStore[1:]
	}
}

// GetAllChats returns a copy of all chats
func GetAllChats() []ChatJson {
	chatStoreMu.Lock()
	defer chatStoreMu.Unlock()
	copied := make([]ChatJson, len(chatStore))
	copy(copied, chatStore)
	return copied
}