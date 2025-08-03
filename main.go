package main

//import (
//	"context"
//	"fmt"
//	"github.com/ksyusha123/procrastinator-library/storage/migrations"
//	"github.com/ydb-platform/ydb-go-sdk/v3"
//	"log"
//	"os"
//	"os/signal"
//	"syscall"
//	"time"
//
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/joho/godotenv"
//	"github.com/ksyusha123/procrastinator-library/bot"
//	"github.com/ksyusha123/procrastinator-library/storage"
//	yc "github.com/ydb-platform/ydb-go-yc"
//)
//
//func main() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	token := os.Getenv("TELEGRAM_BOT_TOKEN")
//	if token == "" {
//		log.Fatal("TELEGRAM_BOT_TOKEN not set in .env file")
//	}
//
//	dsn := os.Getenv("DB_ENDPOINT")
//	if dsn == "" {
//		log.Fatal("DB_ENDPOINT not set in .env file")
//	}
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	db, err := ydb.Open(ctx, "grpcs://ydb.serverless.yandexcloud.net:2135/ru-central1/b1ggmhcf617s1nne9i4g/etnuupd9grd4dlmful9c",
//		yc.WithInternalCA(),
//		yc.WithServiceAccountKeyFileCredentials("./path/to/sa/key/file.json"))
//	if err != nil {
//		fmt.Printf("Driver failed: %v", err)
//	}
//	defer db.Close(ctx)
//
//	if err := migrations.Migrate(ctx, db); err != nil {
//		log.Fatalf("Failed to run migrations: %v", err)
//	}
//
//	storageProvider := storage.NewYDBStorageProvider(db)
//
//	botAPI, err := tgbotapi.NewBotAPI(token)
//	if err != nil {
//		log.Fatalf("Failed to create bot: %v", err)
//	}
//
//	articleBot := bot.New(botAPI, storageProvider)
//
//	stopChan := make(chan os.Signal, 1)
//	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
//
//	go articleBot.Start(ctx)
//
//	<-stopChan
//	log.Println("Shutdown signal received")
//
//	cancel()
//
//	time.Sleep(1 * time.Second)
//	log.Println("Bot stopped gracefully")
//}
