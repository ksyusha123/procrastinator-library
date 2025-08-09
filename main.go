package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const telegramAPIBaseURL = "https://api.telegram.org/bot"

type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
	Result      bool   `json:"result"`
}

func main() {
	// Get bot token from environment variable or prompt
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		fmt.Print("Enter your Telegram bot token: ")
		fmt.Scanln(&botToken)
		if botToken == "" {
			fmt.Println("Bot token is required")
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	initBot(ctx, nil)

	go articleBot.Start(ctx)

	<-stopChan
	log.Println("Shutdown signal received")

	cancel()

	time.Sleep(1 * time.Second)
	log.Println("Bot stopped gracefully")

	//initBot()

	// Get webhook URL from environment variable or prompt
	//webhookURL := os.Getenv("TELEGRAM_WEBHOOK_URL")
	//if webhookURL == "" {
	//	fmt.Print("Enter your webhook URL (HTTPS required): ")
	//	fmt.Scanln(&webhookURL)
	//	if webhookURL == "" {
	//		fmt.Println("Webhook URL is required")
	//		os.Exit(1)
	//	}
	//}
	//
	//// Create the API URL
	//apiURL := fmt.Sprintf("%s%s/setWebhook", telegramAPIBaseURL, botToken)
	//
	//// Prepare the request payload
	//payload := map[string]interface{}{
	//	"url":                  webhookURL,
	//	"drop_pending_updates": true,
	//	// You can add more options here like:
	//	// "max_connections":     40,
	//	// "allowed_updates":     []string{"message", "callback_query"},
	//	// "secret_token":       "your_secret_token",
	//}
	//
	//jsonPayload, err := json.Marshal(payload)
	//if err != nil {
	//	fmt.Printf("Error creating payload: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//// Make the request to set the webhook
	//resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	//if err != nil {
	//	fmt.Printf("Error making request: %v\n", err)
	//	os.Exit(1)
	//}
	//defer resp.Body.Close()
	//
	//// Read and parse the response
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Printf("Error reading response: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//var telegramResp TelegramResponse
	//err = json.Unmarshal(body, &telegramResp)
	//if err != nil {
	//	fmt.Printf("Error parsing response: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//// Output the result
	//if telegramResp.OK {
	//	fmt.Println("Webhook set successfully!")
	//	fmt.Println("Response:", telegramResp.Description)
	//} else {
	//	fmt.Println("Failed to set webhook:", telegramResp.Description)
	//}
}
