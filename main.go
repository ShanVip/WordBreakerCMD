package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func main() {
	token, err := loadToken()
	if err != nil {
		fmt.Println("Ошибка при загрузке токена:", err)
		return
	}
	fmt.Println("Для диалога с ботом введите контекст в начале предложения. А если вам нужен только один ответ от него, просто напишите ваше обращение.")
	reader := bufio.NewReader(os.Stdin)

	client := openai.NewClient(token)
	messages := make([]openai.ChatCompletionMessage, 0)
	input, _ := reader.ReadString('\n')
	firstWord := strings.ToLower(strings.TrimSpace(input))

	if firstWord == "контекст" {
		for {
			fmt.Print("Введите свое обращение: ")
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)

			response, newMessages, err := CreateChatResponse(client, messages, text)
			if err != nil {
				fmt.Printf("ChatCompletion error: %v\n", err)
				continue
			}

			messages = newMessages
			fmt.Println(response)
		}
	}

	single, err := CreateSingleChatCompletion(input, token)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(single)
}

func loadToken() (string, error) {
	viper.SetConfigFile("env/token.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	token := viper.GetString("token")
	return token, nil
}

func CreateSingleChatCompletion(input string, token string) (string, error) {

	client := openai.NewClient(token)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: input,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
func CreateChatResponse(client *openai.Client, messages []openai.ChatCompletionMessage, text string) (string, []openai.ChatCompletionMessage, error) {
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return "", messages, err
	}

	content := resp.Choices[0].Message.Content
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	})

	return content, messages, nil
}
