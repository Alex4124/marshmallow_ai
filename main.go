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

	qdrantClient, err := InitQdrantClient()
	if err != nil {
		log.Fatalf("Error initializing Qdrant client: %v\n", err)
	}

	collectionName := "chat_history"
	vectorSize := 1536

	// Удаляем коллекцию, если она существует
	// err = DeleteCollection(qdrantClient, collectionName)
	// if err != nil {
	// 	log.Printf("Ошибка при удалении коллекции: %v", err)
	// } else {
	// 	log.Printf("Коллекция %s успешно удалена", collectionName)
	// }

	// Создаем коллекцию
	err = CreateCollection(qdrantClient, collectionName, vectorSize)
	if err != nil {
		log.Printf("Ошибка при создании коллекции: %v", err)
	}

	openaiClient := InitOpenAIClient(openaiAPIKey)

	// Start handling updates
	HandleUpdates(bot, openaiClient, qdrantClient, collectionName)
}
