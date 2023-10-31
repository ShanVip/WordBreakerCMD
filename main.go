package main

import (
	"bufio"
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"os"
)

func main() {
	token, err := loadToken()
	if err != nil {
		fmt.Println("Ошибка при загрузке токена:", err)
		return
	}
	fmt.Println("Введите строку:")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	viper.AddConfigPath("/env/token.yaml")

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
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

func loadToken() (string, error) {
	viper.SetConfigFile("env/token.yaml") // Укажите путь к файлу token.yaml
	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	token := viper.GetString("token")
	return token, nil
}
