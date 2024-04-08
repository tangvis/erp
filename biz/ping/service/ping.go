package service

import "context"

type APP interface {
	Ping() string
	PingFail(ctx context.Context) (string, error)
}
