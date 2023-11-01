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

	viper.AddConfigPath("/env/token.yaml")
	client := openai.NewClient(token)
	messages := make([]openai.ChatCompletionMessage, 0)
	input, _ := reader.ReadString('\n')
	firstWord := strings.ToLower(strings.TrimSpace(input))

	if firstWord == "контекст" {
		for {
			fmt.Print("Введите свое обращение: ")
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
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
				fmt.Printf("ChatCompletion error: %v\n", err)
				continue
			}

			content := resp.Choices[0].Message.Content
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: content,
			})
			fmt.Println(content)
		}

	}
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
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

func loadToken() (string, error) {
	viper.SetConfigFile("env/token.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	token := viper.GetString("token")
	return token, nil
}
