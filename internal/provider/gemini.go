package provider

import (
	"context"

	"github.com/Blue-Onion/AI-h/internal/gemini"
	"github.com/Blue-Onion/AI-h/internal/types"
)

func CreateGemini() (*types.Provider, error) {
	client, err := gemini.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	provider := types.Provider{
		Name:         "gemini",
		DefaultModel: "gemini-3.5-flash",
		Stream: func(opts types.StreamOpts) <-chan types.StreamEvent {

		},
	}
	return &provider, nil
}
