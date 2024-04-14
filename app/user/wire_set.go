package user

import (
	"github.com/google/wire"
	"github.com/tangvis/erp/app/user/repository"
	"github.com/tangvis/erp/app/user/service/impl"
)

var ServiceSet = wire.NewSet(
	impl.NewUserAPP,
	RepoSet,
)

var RepoSet = wire.NewSet(
	repository.NewUserRepo,
)
