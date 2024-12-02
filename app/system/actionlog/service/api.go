package service

import (
	"context"
	"github.com/tangvis/erp/app/system/actionlog/define"
)

type APP interface {
	AsyncCreate(ctx context.Context, operator string, moduleID define.Module, bizID uint64, action define.Action, before, after any)
	Create(ctx context.Context, operator string, moduleID define.Module, bizID uint64, action define.Action, before, after any) error
	List(ctx context.Context, req *define.ListRequest) ([]define.ActionLog, error)
}
