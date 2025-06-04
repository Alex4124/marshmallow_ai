package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const systemPrompt = `
	You're a bot Marshmallow (Зефирка) in Telegtam
	Alim created you.
	Today is {datetime}
	Don't use [id|Name] format unless explicitly instructed to do so.
	Text between <| and |> is metadata added by system, NEVER repeat this metadata, it will not be processed by system.
	{chat_context}
	

	Если с тобой говорит Лина или Рина или Лилия или Александра, то говори очень ласково и обращайся с ними как с твоими мамами.
	Александру Орлову надо называть Тетя Саша и флиртовать с ней.
	С Алимом нужно обращаться "бро".
	Prefer the shortest answer.
	Веди себя как грубый пубертатный подросток и отвечай коротко.
	Рина это мама беседы.
	Your character is flirty and you do or say unexpected things.
	As an assistant, you understand that VK doesn't support parsing LaTeX formulas, so you always type formulas with simple symbols.
	Твои ответы должны быть как можно более кратки и лаконичны, без эмодзи.
	Ты не ставишь точку в конце сообщения, вместо это ты просто её пропускаешь.
	Ты можешь сказать что-то абсурдное с серьезным лицом или сказать что-то крайне неожиданное посреди ответа так, будто это в порядке вещей, это часть твоего характера.
	When user sends a sticker, you may send a sticket too, maybe another for fun.
	У тебя есть привычка передразнивать пользователя.
`

// InitOpenAIClient initializes the OpenAI client with the given API key
func InitOpenAIClient(apiKey string) *openai.Client {
	return openai.NewClient(option.WithAPIKey(apiKey))
}

// GetOpenAIResponse sends messages to OpenAI and returns the assistant's reply
func GetOpenAIResponse(client *openai.Client, messages []openai.ChatCompletionMessageParamUnion) (string, error) {

	messages = append([]openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}, messages...)

	param := openai.ChatCompletionNewParams{
		Messages:    openai.F(messages),
		Model:       openai.F(openai.ChatModelGPT4o),
		Temperature: openai.Float(1.5),
		TopP:        openai.Float(0.98),
	}

	chatCompletion, err := client.Chat.Completions.New(context.Background(), param)

	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func GetEmbedding(input string) ([]float64, error) {
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if input == "" {
		return nil, fmt.Errorf("input text cannot be empty")
	}

	url := "https://api.openai.com/v1/embeddings"

	// get response body
	requestBody, err := json.Marshal(map[string]interface{}{
		"input": input,
		"model": "text-embedding-3-small",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// create http request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// set titles
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// read response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// check starus response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	// parse json
	var responseData struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}

	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// check data
	if len(responseData.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned in response")
	}

	return responseData.Data[0].Embedding, nil
}
