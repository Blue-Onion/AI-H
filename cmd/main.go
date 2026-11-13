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
	for res, err := range client.Models.GenerateContentStream(
		ctx,
		modelName,
		genai.Text(prompt),
		nil,
	) {
		if err != nil {
			return err
		}
		if res != nil {
			fmt.Print(res.Text())
		}
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
