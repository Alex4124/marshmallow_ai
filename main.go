package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading the .env file: %v\n", err)
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")

	if telegramToken == "" || openaiAPIKey == "" {
		log.Fatalf("Telegram or OpenAI tokens are not found. Check the .env file")
	}

	// Initialize the Telegram bot and OpenAI client
	bot, err := InitTelegramBot(telegramToken)
	if err != nil {
		log.Fatalf("Error initializing bot: %v\n", err)
	}

	openaiClient := InitOpenAIClient(openaiAPIKey)

	// Start handling updates
	HandleUpdates(bot, openaiClient)
}
