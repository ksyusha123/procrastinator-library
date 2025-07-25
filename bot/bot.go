package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

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

func (b *Bot) Start(ctx context.Context) {
	log.Println("Starting article bot...")

	_, err := readLastUpdateId()
	if err != nil {
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botAPI.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			b.handleUpdate(&update)
			writeLastUpdateId(update.UpdateID)
		case <-ctx.Done():
			log.Println("Stopping bot updates")
			return
		}
	}
}

func readLastUpdateId() (int, error) {
	content, err := os.ReadFile("offset")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return 0, err
	}
	num, err := strconv.Atoi(string(content))
	if err != nil {
		fmt.Printf("Error converting line: %v\n", err)
		return 0, err
	}
	return num, nil
}

func writeLastUpdateId(lastUpdateId int) error {
	return os.WriteFile("offset", []byte(strconv.Itoa(lastUpdateId+1)), 0644)
}
