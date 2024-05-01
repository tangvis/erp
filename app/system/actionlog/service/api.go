package service

import (
	"context"
	"github.com/tangvis/erp/app/system/actionlog/define"
)

type APPActionLog interface {
	Create(ctx context.Context, operator string, moduleID, bizID uint64, action define.Action, before, after any) error
}
