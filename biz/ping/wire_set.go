package ping

import (
	"github.com/google/wire"

	"github.com/tangvis/erp/biz/ping/access"
	"github.com/tangvis/erp/biz/ping/service/impl"
)

var APISet = wire.NewSet(
	BizSet,
	access.NewController,
)

var BizSet = wire.NewSet(
	impl.NewPing,
)
