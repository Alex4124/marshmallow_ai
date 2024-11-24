package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	// Loading environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading the .env file: %v\n", err)
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")

	if telegramToken == "" || openaiAPIKey == "" {
		log.Fatalf("Telegram or OpenAI tokens are not found. Check the .env file")
	}

	// Initializing a Telegram bot
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatalf("Initalization error: %v\n", err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Initializing OpenAI
	openaiClient := openai.NewClient(option.WithAPIKey(openaiAPIKey))

	// Setting up an update channel to receive messages
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	// Processing updates
	for update := range updates {
		if update.Message != nil {

			shouldRespond := false
			mentionBot := false

			if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {

				// Remove the name mentioned
				for _, ent := range update.Message.Entities {
					if ent.Type == "mention" {
						mention := update.Message.Text[ent.Offset : ent.Offset+ent.Length]
						if strings.EqualFold(mention, "@"+bot.Self.UserName) {
							mentionBot = true
							break
						}
					}
				}
				// The message is a reply to a bot message
				isReplyToBot := false
				if update.Message.ReplyToMessage != nil {
					if update.Message.ReplyToMessage.From.UserName == bot.Self.UserName {
						isReplyToBot = true
					}
				}

				// If none of the conditions are met, skip the message
				if !mentionBot && !isReplyToBot {
					continue
				}

				shouldRespond = true

			} else {
				shouldRespond = true
			}

			if !shouldRespond {
				continue
			}

			// Receive text from the user
			chatID := update.Message.Chat.ID
			userMessage := update.Message.Text

			if mentionBot {
				// Remove bot mention if it's present in the message
				userMessage = strings.Replace(userMessage, fmt.Sprintf("@%s", bot.Self.UserName), "", 1)
				userMessage = strings.TrimSpace(userMessage)
			}

			if userMessage == "" {
				log.Printf("Empty message from user: [%s]\n", update.Message.From.UserName)
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, userMessage)

			// Receive response from OpenAI
			chatCompletion, err := openaiClient.Chat.Completions.New(
				context.TODO(),
				openai.ChatCompletionNewParams{
					Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
						openai.UserMessage(userMessage),
					}),
					Model: openai.F(openai.ChatModelGPT3_5Turbo),
				},
			)

			if err != nil {
				log.Printf("Error sending message to OpenAI: %v", err)
			}

			// Create a message to send back the same text
			responseMessage := chatCompletion.Choices[0].Message.Content
			msg := tgbotapi.NewMessage(chatID, responseMessage)

			// Highlights the message to which you are replying
			msg.ReplyToMessageID = update.Message.MessageID

			_, err = bot.Send(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
