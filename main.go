package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/ksyusha123/procrastinator-library/bot"
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

	db, err := storage.NewSQLiteStorage("articles.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// articleBot := bot.New(botAPI, db, cfg.OpenAIKey)

	articleBot := bot.New(botAPI, db)

	go articleBot.StartNotificationScheduler()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go articleBot.Start(ctx)

	<-stopChan
	log.Println("Shutdown signal received")

	cancel()

	time.Sleep(1 * time.Second)
	log.Println("Bot stopped gracefully")
}
