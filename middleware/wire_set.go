package middleware

import (
	"github.com/google/wire"
	"github.com/tangvis/erp/middleware/engine"
)

var Set = wire.NewSet(
	engine.NewEngine,
	engine.NewRedisStore,
)
