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
}
