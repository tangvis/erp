package ping

import (
	"github.com/google/wire"

	"github.com/tangvis/erp/biz/ping/access"
	"github.com/tangvis/erp/biz/ping/service/impl"
)

var APISet = wire.NewSet(
	ServiceSet,
	access.NewController,
)

var ServiceSet = wire.NewSet(
	impl.NewPing,
)
