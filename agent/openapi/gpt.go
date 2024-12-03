package openapi

import (
	"context"
	"github.com/sashabaranov/go-openai"
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
