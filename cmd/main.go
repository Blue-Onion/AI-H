package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Blue-Onion/AI-h/config"
	"google.golang.org/genai"
)

func GetAiResponse(userQuery string) (string, error) {
	apiKey := config.GetConfig().ApiKey
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return "", err
	}

	resp, err := client.Models.GenerateContent(ctx, "models/gemini-2.5-flash", genai.Text(userQuery), nil)
	if err != nil {
		return "", err

	}
	return resp.Candidates[0].Content.Parts[0].Text, nil
}
func main() {
	p := flag.String("p", "", "Value for P")
	m := flag.String("m", "", "Value for m")
	flag.Parse()
	fmt.Println(*p)
	fmt.Println(*m)
	res, err := GetAiResponse(*p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
