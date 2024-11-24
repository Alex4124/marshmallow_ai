# MarshMallowAI - Telegram Bot Powered by OpenAI API

MarshMallowAI is a Telegram bot that integrates OpenAI's GPT-3.5 Turbo model for intelligent and conversational responses. Whether in group chats or private messages, MarshMallowAI processes user queries and provides thoughtful replies, creating a seamless experience for engaging AI-driven conversations.

---

## Features

- **Smart Replies**: Leverages OpenAI's GPT-3.5 Turbo model to respond intelligently to user queries.
- **Group Chat Support**: Responds only when mentioned or replied to directly in group chats.
- **Private Messages**: Fully interactive and responsive in private conversations.
- **Customizable**: Built with environment variables to securely store API keys.

---

## Installation and Setup

### Prerequisites

- [Go](https://golang.org/dl/) installed (version 1.19 or higher recommended).
- A Telegram bot token from [BotFather](https://core.telegram.org/bots#botfather).
- An OpenAI API key from the [OpenAI website](https://openai.com/).
- [Git](https://git-scm.com/).

### Clone the Repository

```bash
git clone https://github.com/Alex4124/marshmallow_ai
cd marshmallow_ai
```

### Install Dependencies

Ensure that the required dependencies are installed:

```bash
go mod tidy
```

### Set Up Environment Variables

Create a `.env` file in the project root directory and add your API keys:

```plaintext
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
OPENAI_API_KEY=your_openai_api_key
```

### Run the Bot

Start the bot with the following command:

```bash
go run main.go
```

---

## How to Use

1. **Add MarshMallowAI to Your Telegram**: Start a chat with the bot or add it to a group.
2. **Interact**:
   - In **private messages**, simply type your query.
   - In **group chats**, mention the bot using `@MarshMallowAIBot` or reply to a message from the bot.
3. **Wait for MarshMallowAI's Response**: The bot will respond using OpenAI's GPT-3.5 Turbo model.

---

## Development

### Project Structure

- `main.go`: Core application logic.
- `.env`: Stores environment variables securely (not included in the repository).
- `go.mod`, `go.sum`: Dependency management files.

### Key Libraries Used

- [tgbotapi](https://github.com/go-telegram-bot-api/telegram-bot-api): For interacting with the Telegram Bot API.
- [godotenv](https://github.com/joho/godotenv): For managing environment variables.
- [openai-go](https://github.com/openai/openai-go): For communicating with the OpenAI API.

---

## Troubleshooting

### Common Issues

- **Error: Tokens not found**: Ensure `TELEGRAM_BOT_TOKEN` and `OPENAI_API_KEY` are correctly set in the `.env` file.
- **Error: Bot not responding in groups**: Verify that the bot has proper permissions and is mentioned correctly.
- **OpenAI API errors**: Check API key validity and rate limits on your OpenAI account.

### Debugging

Enable debug mode for the Telegram bot by setting `bot.Debug = true` in the code. Logs will display additional information.

---

## Future Enhancements

- Implement command-based interactions (e.g., `/start`, `/help`).
- Add support for multimedia messages.
- Improve error handling and logging.
- Expand configuration options.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Contributing

Contributions are welcome! Please fork the repository, create a new branch for your changes, and submit a pull request.

---

## Contact

For inquiries, reach out to the project maintainer at **leskinen01@mail.ru**.

Happy chatting with MarshMallowAI! 🎉
