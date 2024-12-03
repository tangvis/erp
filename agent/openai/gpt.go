package openai

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"time"

	logutil "github.com/tangvis/erp/pkg/log"
)

const (
	OpenAIToken        = ""
	maxAttempts        = 3
	retryDelayDuration = 3 * time.Second
)

type GPTClient struct {
	client *openai.Client
	token  string
}

func NewGPTClient(token string) *GPTClient {
	return &GPTClient{
		client: openai.NewClient(token),
		token:  "",
	}
}

func (g *GPTClient) SimpleChat(ctx context.Context, text string, model string) (string, error) {
	return retry(ctx, func() (string, error) {
		return g.simpleChat(ctx, text, model)
	}, maxAttempts, retryDelayDuration)
}

func (g *GPTClient) simpleChat(ctx context.Context, text string, model string) (string, error) {
	resp, err := g.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (g *GPTClient) SimpleChat4oMini(ctx context.Context, text string) (string, error) {
	return g.SimpleChat(ctx, text, openai.GPT4oMini)
}

func (g *GPTClient) Caption(ctx context.Context, reader io.Reader) (string, error) {
	return retry(ctx, func() (string, error) {
		return g.caption(ctx, reader)
	}, maxAttempts, retryDelayDuration)
}

func (g *GPTClient) caption(ctx context.Context, reader io.Reader) (string, error) {
	resp, err := g.client.CreateTranscription(ctx, openai.AudioRequest{
		Model:  openai.Whisper1,
		Reader: reader,
		Format: openai.AudioResponseFormatSRT,
	})
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

// retry function that accepts a retryable operation and tries it until success or max attempts are reached
func retry[T any](ctx context.Context, operation func() (T, error), maxAttempts int, delay time.Duration) (T, error) {
	var result T
	var err error

	// Try the operation up to maxAttempts times
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, err = operation() // Call the provided operation function
		if err == nil {
			return result, nil // Success, return the result
		}
		e := &openai.APIError{}
		if errors.As(err, &e) {
			switch e.HTTPStatusCode {
			case 401:
				// invalid auth or key (do not retry)
				return result, err
			case 429:
				// rate limiting or engine overload (wait and retry)
			case 500:
				// openai server error (retry)
			default:
				return result, err
			}
		}

		// If there is an error and max attempts haven't been reached, retry
		if attempt >= maxAttempts {
			// Max attempts reached, return the error
			return result, fmt.Errorf("failed after %d attempts: %w", maxAttempts, err)

		}
		logutil.CtxInfoF(ctx, "Attempt %d failed, retrying...", attempt)
		time.Sleep(delay) // Optional delay between retries
	}
	return result, err
}
