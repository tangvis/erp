package openai

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"io"
)

const (
	OpenAIToken = ""
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
