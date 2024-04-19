package impl

import (
	"context"

	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/agent/redis"
	"github.com/tangvis/erp/app/ping/service"
	"github.com/tangvis/erp/common"
	logutil "github.com/tangvis/erp/pkg/log"
)

type Ping struct {
	db    *mysql.DB
	cache redis.Cache
}

func NewPing(
	db *mysql.DB,
	cache redis.Cache,
) service.APP {
	return &Ping{
		db:    db,
		cache: cache,
	}
}

func (p *Ping) Ping() string {
	return "pong"
}

func (p *Ping) PingFail(ctx context.Context) (string, error) {
	logutil.CtxErrorF(ctx, "manual fail ping")
	return "", common.ErrPingFailedTest
}
