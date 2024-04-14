package service

import (
	"context"
	"github.com/tangvis/erp/biz/user/service/define"
)

type APP interface {
	Create(ctx context.Context, req define.SignupRequest) (uint64, error)
}
