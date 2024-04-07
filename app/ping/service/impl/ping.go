package impl

import "github.com/tangvis/erp/app/ping/service"

type Ping struct {
}

func NewPing() service.APP {
	return &Ping{}
}

func (p *Ping) Ping() {
	panic("ping")
}
