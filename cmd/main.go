package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Blue-Onion/AI-h/config"
	"google.golang.org/genai"
)

func GetAiResponse(prompt string, model string) error {
	var usage *genai.GenerateContentResponseUsageMetadata
	apiKey := config.GetConfig().ApiKey
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return err
	}
	modelName := "models/gemini-3.5-flash"
	if len(model) > 0 {
		modelName = fmt.Sprintf("models/%s", model)
	}
	for chunk, err := range client.Models.GenerateContentStream(
		ctx,
		modelName,
		genai.Text(prompt),
		nil,
	) {
		if err != nil {
			return err
		}
		if chunk == nil {
			continue
		}

		fmt.Print(chunk.Text())
		if chunk.UsageMetadata != nil {
			usage = chunk.UsageMetadata
		}
	}
	fmt.Println()
	if usage != nil {
		fmt.Printf(
			"input_tokens=%d output_tokens=%d thought_tokens=%d total_tokens=%d",
			usage.PromptTokenCount,
			usage.CandidatesTokenCount,
			usage.ThoughtsTokenCount,
			usage.TotalTokenCount,
		)
	}
	return nil
}
func main() {
	p := flag.String("p", "", "Value for P")
	m := flag.String("m", "", "Value for m")
	flag.Parse()
	err := GetAiResponse(*p, *m)
	if err != nil {
		log.Fatal(err)
	}
}
