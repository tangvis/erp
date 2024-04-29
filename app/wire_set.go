package app

import (
	"github.com/google/wire"

	"github.com/tangvis/erp/app/apirate"
	"github.com/tangvis/erp/app/ping"
	"github.com/tangvis/erp/app/product"
	"github.com/tangvis/erp/app/user"
)

var ServiceSet = wire.NewSet(
	apirate.ServiceSet,
	ping.ServiceSet,
	user.ServiceSet,
	product.ServiceSet,
)
