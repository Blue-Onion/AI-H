// Package gemini provides streaming access to Gemini.
package gemini

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
}

func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("create Gemini client: %w", err)
	}
	return &Client{client: client}, nil
}

// Stream writes response text to out and returns the last reported token usage.
func (c *Client) Stream(ctx context.Context, prompt, model string, out io.Writer) (*genai.GenerateContentResponseUsageMetadata, error) {
	var usage *genai.GenerateContentResponseUsageMetadata
	modelName := "models/gemini-3.5-flash"
	if len(model) > 0 {
		modelName = fmt.Sprintf("models/%s", model)
	}
	for chunk, err := range c.client.Models.GenerateContentStream(
		ctx,
		modelName,
		genai.Text(prompt),
		nil,
	) {
		if err != nil {
			return nil, err
		}
		if chunk == nil {
			continue
		}

		if _, err := io.WriteString(out, chunk.Text()); err != nil {
			return nil, err
		}
		if chunk.UsageMetadata != nil {
			usage = chunk.UsageMetadata
		}
	}
	return usage, nil
}
