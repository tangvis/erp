package openai

import (
	"context"
	"io"
)

type API interface {
	SimpleChat(ctx context.Context, text string, model string) (string, error)
	SimpleChat4oMini(ctx context.Context, text string) (string, error)
	Caption(ctx context.Context, reader io.Reader) (string, error)
}
