package ping

import (
	"github.com/google/wire"

	"github.com/tangvis/erp/app/ping/access"
	"github.com/tangvis/erp/app/ping/service/impl"
)

var APISet = wire.NewSet(
	access.NewController,
)

var ServiceSet = wire.NewSet(
	impl.NewPing,
)
