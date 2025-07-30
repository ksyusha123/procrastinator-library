package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func init() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		panic("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(fmt.Sprintf("Error creating bot: %v", err))
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
}

type APIGatewayRequest struct {
	OperationID string `json:"operationId"`
	Resource    string `json:"resource"`

	HTTPMethod string `json:"httpMethod"`

	Path           string            `json:"path"`
	PathParameters map[string]string `json:"pathParameters"`

	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`

	QueryStringParameters           map[string]string   `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"`

	Parameters           map[string]string   `json:"parameters"`
	MultiValueParameters map[string][]string `json:"multiValueParameters"`

	Body            string `json:"body"`
	IsBase64Encoded bool   `json:"isBase64Encoded,omitempty"`

	RequestContext interface{} `json:"requestContext"`
}

type APIGatewayResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
}

func Greet(ctx context.Context, event *APIGatewayRequest) (*APIGatewayResponse, error) {
	req := &tgbotapi.Update{}

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing body: %v", err)
	}

	fmt.Println(event.HTTPMethod, event.Path)

	msg := tgbotapi.NewMessage(req.Message.Chat.ID, "You said: "+"meeeow")
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("an error has occurred when sending message: ", err)
		return nil, err
	}

	return &APIGatewayResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Hello, %s", req.Message.Chat.ID),
	}, nil
}
