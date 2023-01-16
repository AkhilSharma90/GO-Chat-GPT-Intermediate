package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		panic("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	const inputFile = "./input_with_code.txt"
	fileBytes, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	
	msgPrefix := "give me a short list of libraries that are used in the code\n```python\n"
	msgSuffix := "\n```"
	msg := msgPrefix + string(fileBytes) + msgSuffix

	outputBuilder := strings.Builder{}
	err = client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			msg,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		outputBuilder.WriteString(resp.Choices[0].Text)
	})
	if err != nil {
		log.Fatalln(err)
	}
	output := strings.TrimSpace(outputBuilder.String())

	const outputFile = "./output.txt"
	err = os.WriteFile(outputFile, []byte(output), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
}
