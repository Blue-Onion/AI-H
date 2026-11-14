package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"google.golang.org/genai"
)

var output string = ` 
Bound to fall in love
Bound to fall in love (uh-huh, honey)
All them other niggas lame, and you know it now
When a real nigga hold you down, (alright) you s'pposed to drown
Bound (bound) to fall in love
Bound (bound) to fall in love (uh-huh, honey)
What you doing in the club on a Thursday?
She say she only here for her girl birthday
They ordered champagne but still look thirsty
Rock Forever 21 but just turned thirty
I know I got a bad reputation
Walking 'round, always mad reputation
Leave a pretty girl sad reputation
Start a Fight Club, Brad reputation
I turnt the nightclub out of the basement
I'll turn the plane 'round, your ass keep complaining
How you gon' be mad on vacation?
Dutty whining 'round all these Jamaicans
Uh, this that prom shit
This that what we do, don't tell your mom shit
This that red cup, all on the lawn shit
Got a fresh cut, straight out the salon, bitch
I know you're tired of loving, of loving
With nobody to love, nobody, no- (uh-huh, honey)
Close your eyes and let the word paint a thousand pictures
One good girl (alright) is worth a thousand bitches
Bound (bound) to fall in love
Bound (bound) to fall in love (uh-huh, honey)
I wanna fuck you hard on the sink
After that, give you something to drink
Step back, can't get spunk on the mink
I mean damn, what would Jeromey Romey Romey Rome think?
Hey, you remember where we first met?
Okay, I don't remember where we first met
But hey, admitting is the first step
And, ay, you know ain't nobody perfect
And I know, with the hoes I got the worst rep
But ay, the backstroke I'm tryna perfect
And ay, ayo, we made to Thanksgiving
So ay, maybe we can make it to Christmas
She asked me what I wished for on my wishlist
Have you ever asked your bitch for other bitches?
Maybe we could still make it to the church steps
But first, you gon' remember how to forget
After all these long-ass verses
I'm tired, you tired, (uh-huh, honey) Jesus wept
I know you're tired (tired) of loving, of loving
With nobody to love, nobody, nobody
So just grab somebody, no leaving this party
With nobody to love, nobody, nobody (to love)
(Uh-huh, honey)
Jerome's in the house, watch your mouth
Jerome's in the house, watch your mouth
Bound (bound) to fall in love
Bound (bound) to fall in love (uh-huh, honey)
`

// dummyTransport mimics Gemini's server-sent events without network access.
type dummyTransport struct {
	chunks []string
}

func (d dummyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var events strings.Builder
	for i, chunk := range d.chunks {
		event := map[string]any{
			"candidates": []any{map[string]any{
				"content": map[string]any{
					"role":  "model",
					"parts": []any{map[string]string{"text": chunk}},
				},
			}},
		}
		// These are dummy usage counts, independent of our simulated chunk sizes.
		if i == len(d.chunks)-1 {
			event["usageMetadata"] = map[string]int{
				"promptTokenCount": 5, "candidatesTokenCount": 4,
				"thoughtsTokenCount": 2, "totalTokenCount": 11,
			}
		}
		data, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(&events, "data: %s\n\n", data)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"text/event-stream"}},
		Body:       io.NopCloser(strings.NewReader(events.String())),
		Request:    req,
	}, nil
}

type chunkWriter struct {
	chunks []string
}

func (w *chunkWriter) Write(p []byte) (int, error) {
	w.chunks = append(w.chunks, string(p))
	return len(p), nil
}

func Tokenize(text string, maxChars int) []string {
	if maxChars <= 0 {
		panic("maxChars must be greater than zero")
	}

	var tokens []string

	for _, word := range strings.Fields(text) {
		chars := []rune(word)

		for len(chars) > maxChars {
			tokens = append(tokens, string(chars[:maxChars]))
			chars = chars[maxChars:]
		}

		if len(chars) > 0 {
			tokens = append(tokens, string(chars))
		}
	}

	return tokens
}
func TestStreamDummyLLM(t *testing.T) {
	ctx := context.Background()
	want := Tokenize(output, 6)
	sdk, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:     "dummy-key",
		Backend:    genai.BackendGeminiAPI,
		HTTPClient: &http.Client{Transport: dummyTransport{chunks: want}},
	})
	if err != nil {
		t.Fatal(err)
	}
	client := &Client{client: sdk}
	out := &chunkWriter{}
	usage, err := client.Stream(ctx, "Say hello", "", out)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out.chunks, want) {
		t.Fatalf("streamed chunks = %q, want %q", out.chunks, want)
	}
	if usage == nil {
		t.Fatal("expected token usage")
	}
	if usage.PromptTokenCount != 5 || usage.CandidatesTokenCount != 4 || usage.ThoughtsTokenCount != 2 || usage.TotalTokenCount != 11 {
		t.Fatalf("unexpected token usage: %+v", usage)
	}
	t.Logf("Streamed chunks: %q", out.chunks)
	t.Logf("Response: %s", strings.Join(out.chunks, ""))
	t.Logf("input_tokens=%d output_tokens=%d thought_tokens=%d total_tokens=%d",
		usage.PromptTokenCount, usage.CandidatesTokenCount, usage.ThoughtsTokenCount, usage.TotalTokenCount)
}
