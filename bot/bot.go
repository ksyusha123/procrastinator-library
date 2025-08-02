package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ksyusha123/procrastinator-library/storage"
)

type Article = storage.Article

type Bot struct {
	botAPI         *tgbotapi.BotAPI
	articleStorage storage.ArticleStorage
	userStorage    storage.UserStorage
	commands       map[string]string
}

func New(botAPI *tgbotapi.BotAPI, db storage.SQLiteDb) *Bot {
	return &Bot{
		botAPI:         botAPI,
		articleStorage: &db,
		userStorage:    &db,
		commands: map[string]string{
			"save":   "Save an article (reply to message or provide URL)",
			"list":   "List your saved articles",
			"read":   "Mark article as read (provide article ID)",
			"delete": "Delete article (provide article ID)",
			"help":   "Show available commands",
		},
	}
}
