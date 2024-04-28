package system

import (
	"github.com/google/wire"

	"github.com/tangvis/erp/app/system/email/repository"
	"github.com/tangvis/erp/app/system/email/service/impl"
)

var Set = wire.NewSet(
	EmailSet,
)

var EmailSet = wire.NewSet(
	impl.NewEmailAPP,
	repository.NewRepo,
)
