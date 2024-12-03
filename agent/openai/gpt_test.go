package openai

import (
	"fmt"
	"testing"
)

func foo() error {
	var (
		attempt int
		f       func() error
	)
	f = func() error {
		attempt++
		if attempt > 3 {
			return fmt.Errorf("attempts too many times")
		}
		return f()
	}

	return f()
}

func TestGPTClient_Caption(t *testing.T) {
}
