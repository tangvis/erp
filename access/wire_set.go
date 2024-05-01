package access

import (
	"github.com/google/wire"
	"github.com/tangvis/erp/access/ping"
	"github.com/tangvis/erp/access/product"
	"github.com/tangvis/erp/access/system"
	"github.com/tangvis/erp/access/user"
)

var HTTPSet = wire.NewSet(
	ping.NewController,
	user.NewController,
	product.NewController,
	system.NewController,
)
