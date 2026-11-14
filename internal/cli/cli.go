// Package cli handles command-line arguments and application setup.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/Blue-Onion/AI-h/config"
	"github.com/Blue-Onion/AI-h/internal/gemini"
)

// Run executes one invocation without relying on process-wide flags or output.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("ai-h", flag.ContinueOnError)
	flags.SetOutput(stderr)
	prompt := flags.String("p", "", "Prompt to send to Gemini")
	model := flags.String("m", "", "Gemini model (default: gemini-3.5-flash)")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	client, err := gemini.NewClient(ctx, config.GetConfig().ApiKey)
	if err != nil {
		return err
	}
	usage, err := client.Stream(ctx, *prompt, *model, stdout)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return err
	}
	if usage != nil {
		_, err = fmt.Fprintf(stdout, "input_tokens=%d output_tokens=%d thought_tokens=%d total_tokens=%d",
			usage.PromptTokenCount, usage.CandidatesTokenCount, usage.ThoughtsTokenCount, usage.TotalTokenCount)
	}
	return err
}
