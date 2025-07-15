package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ksyusha123/procrastinator-library/bot"
	// "github.com/ksyusha123/procrastinator-library/config"
	"github.com/ksyusha123/procrastinator-library/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN not set in .env file")
	}

	db, err := storage.NewSQLiteDB("articles.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// articleBot := bot.New(botAPI, db, cfg.OpenAIKey)

	articleBot := bot.New(botAPI, db)

	// go bot.StartNotificationScheduler(articleBot)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	articleBot.Start()

	<-sigChan
	log.Println("Shutting down bot...")
}
