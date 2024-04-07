package impl

import (
	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/biz/ping/service"
)

type Ping struct {
	db *mysql.DB
}

func NewPing(
	db *mysql.DB,
) service.APP {
	return &Ping{
		db: db,
	}
}

func (p *Ping) Ping() string {
	return "pong"
}
