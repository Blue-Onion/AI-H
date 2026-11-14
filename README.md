# AI Harness

A basic Go CLI that sends a prompt to Gemini and prints its response.
Requires Go 1.27.2 or newer.

Create an API key in [Google AI Studio](https://aistudio.google.com/apikey), then run:

```sh
export GOOGLE_API_KEY="your-api-key"
go run cmd/main.go -p "Explain Go interfaces"
```

Alternatively, create a `.env` file in the directory where you run the command:

```dotenv
GOOGLE_API_KEY=your-api-key
```

The CLI loads `.env` automatically. Exported environment variables take precedence.
The `.env` file is ignored by Git.

The default model is `gemini-3.5-flash`. Override it with `-m`:

```sh
go run cmd/main.go -p "Hello" -m gemini-3.5-flash-lite
```

Uses the [Gemini generateContent REST API](https://ai.google.dev/api/generate-content).
Responses go to stdout; errors go to stderr with a nonzero exit status.
Press Ctrl+C to cancel a request.

## Project structure

- `cmd/main.go`: process startup, cancellation, and exit status.
- `internal/cli`: flag parsing, application wiring, and terminal formatting.
- `internal/gemini`: Gemini client creation, and response streaming.
- `config`: environment and `.env` configuration loading.

The `internal` packages are private to this module. The API client accepts a
context and output writer so callers control cancellation and output destination.
