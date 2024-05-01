package system

import (
	"github.com/google/wire"

	actionLogRepo "github.com/tangvis/erp/app/system/actionlog/repository"
	actionLogImpl "github.com/tangvis/erp/app/system/actionlog/service/impl"
	emailRepo "github.com/tangvis/erp/app/system/email/repository"
	emailImpl "github.com/tangvis/erp/app/system/email/service/impl"
)

var Set = wire.NewSet(
	EmailSet,
	ActionLogSet,
)

var EmailSet = wire.NewSet(
	emailImpl.NewEmailAPP,
	emailRepo.NewRepo,
)

var ActionLogSet = wire.NewSet(
	actionLogRepo.NewRepoImpl,
	actionLogImpl.NewActionLogAPP,
)
