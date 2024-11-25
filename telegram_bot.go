package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/openai/openai-go"
)

var chatHistories = make(map[int64][]openai.ChatCompletionMessageParamUnion)

// InitTelegramBot initializes the Telegram bot with the given token
func InitTelegramBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot, nil
}

// HandleUpdates processes incoming updates from Telegram
func HandleUpdates(bot *tgbotapi.BotAPI, openaiClient *openai.Client) {
	// Set up an update channel to receive messages
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	// Process incoming updates
	for update := range updates {
		if update.Message != nil {
			shouldRespond := false
			mentionBot := false

			if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
				// Check if the bot is mentioned in the group chat
				for _, ent := range update.Message.Entities {
					if ent.Type == "mention" {
						mention := update.Message.Text[ent.Offset : ent.Offset+ent.Length]
						if strings.EqualFold(mention, "@"+bot.Self.UserName) {
							mentionBot = true
							break
						}
					}
				}
				// Check if the message is a reply to the bot
				isReplyToBot := false
				if update.Message.ReplyToMessage != nil {
					if update.Message.ReplyToMessage.From.UserName == bot.Self.UserName {
						isReplyToBot = true
					}
				}

				// If neither condition is met, skip the message
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
				userMessage = strings.Replace(userMessage, "@"+bot.Self.UserName, "", 1)
				userMessage = strings.TrimSpace(userMessage)
			}

			if userMessage == "" {
				log.Printf("Empty message from user: [%s]\n", update.Message.From.UserName)
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, userMessage)

			chatHistories[chatID] = append(chatHistories[chatID], openai.UserMessage(userMessage))

			// Limit the chat history to the last 10 messages
			if len(chatHistories[chatID]) > 10 {
				chatHistories[chatID] = chatHistories[chatID][len(chatHistories[chatID])-10:]
			}

			// Get response from OpenAI
			responseMessage, err := GetOpenAIResponse(openaiClient, chatHistories[chatID])
			if err != nil {
				log.Printf("Error communicating with OpenAI: %v", err)
				continue
			}

			chatHistories[chatID] = append(chatHistories[chatID], openai.AssistantMessage(responseMessage))

			// Send the response back to the user
			msg := tgbotapi.NewMessage(chatID, responseMessage)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err = bot.Send(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
