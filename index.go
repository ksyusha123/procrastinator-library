package main

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ksyusha123/procrastinator-library/api"
	"github.com/ksyusha123/procrastinator-library/storage"
	"github.com/ksyusha123/procrastinator-library/storage/migrations"
	yc "github.com/ydb-platform/ydb-go-yc-metadata"
	"log"
	"os"

	"github.com/ksyusha123/procrastinator-library/bot"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

var articleBot *bot.Bot

func init() {
	token := getVariable("TELEGRAM_BOT_TOKEN")
	dsn := getVariable("DB_ENDPOINT")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := ydb.Open(ctx, dsn, yc.WithInternalCA(), yc.WithCredentials())
	if err != nil {
		fmt.Printf("Driver failed: %v", err)
	}
	defer db.Close(ctx)
	fmt.Printf("connected to %s, database '%s'", db.Endpoint(), db.Name())

	if err := migrations.Migrate(ctx, db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	storageProvider := storage.NewYDBStorageProvider(db)

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	articleBot = bot.New(botAPI, storageProvider)
}

func getVariable(name string) string {
	variable := os.Getenv(name)
	if variable == "" {
		panic("DB_ENDPOINT environment variable not set")
	}
	return variable
}

func Greet(ctx context.Context, event *api.APIGatewayRequest) (*api.APIGatewayResponse, error) {
	update := &tgbotapi.Update{}

	if err := json.Unmarshal([]byte(event.Body), &update); err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing body: %v", err)
	}

	fmt.Println(event.HTTPMethod, event.Path)

	articleBot.HandleUpdate(ctx, update)

	return &api.APIGatewayResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Hello, %s", update.Message.Chat.ID),
	}, nil
}
