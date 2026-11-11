# AI Harness

A basic Go CLI that sends a prompt to Gemini and prints its response.
Requires Go 1.27.2 or newer.

Create an API key in [Google AI Studio](https://aistudio.google.com/apikey), then run:

```sh
export GEMINI_API_KEY="your-api-key"
go run cmd/main.go -p "Explain Go interfaces"
```

Alternatively, create a `.env` file in the directory where you run the command:

```dotenv
GEMINI_API_KEY=your-api-key
```

The CLI loads `.env` automatically. Exported environment variables take precedence.
The `.env` file is ignored by Git.

The default model is `gemini-flash-latest`. Override it with `-model`:

```sh
go run cmd/main.go -p "Hello" -model gemini-flash-latest
```

Uses the [Gemini generateContent REST API](https://ai.google.dev/api/generate-content).
Responses go to stdout; errors go to stderr with a nonzero exit status.
Requests time out after two minutes.
