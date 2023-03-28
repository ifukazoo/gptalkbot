package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {

	path := os.Getenv("HOME")

	err := godotenv.Load(filepath.Join(path, ".env"))
	if err != nil {
		log.Fatal(err)
	}
	key := os.Getenv("APIKEY")
	if len(key) == 0 {
		log.Fatal(fmt.Errorf("not found apikey"))
	}
	reader := bufio.NewScanner(os.Stdin)

	dialog := []openai.ChatCompletionMessage{}
	for reader.Scan() {
		input := reader.Text()
		dialog = append(dialog, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})
		resp, totalToken, err := response(key, dialog)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
		dialog = append(dialog, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: resp,
		})
		if totalToken > 3000 {
			dialog = dialog[2:]
		}
	}
}

func response(key string, dialog []openai.ChatCompletionMessage) (string, int, error) {
	client := openai.NewClient(key)
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT432K,
			Messages: dialog,
		},
	)
	if err != nil {
		return "", 0, err
	}
	return resp.Choices[0].Message.Content, resp.Usage.TotalTokens, nil
}
