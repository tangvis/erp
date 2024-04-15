package ping

import (
	"github.com/google/wire"

	"github.com/tangvis/erp/app/ping/service/impl"
)

var ServiceSet = wire.NewSet(
	impl.NewPing,
)
