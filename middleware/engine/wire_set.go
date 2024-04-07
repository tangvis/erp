package engine

import (
	"github.com/google/wire"
)

var EngineSet = wire.NewSet(
	NewEngine,
)
