package apirate

import (
	"github.com/google/wire"

	"github.com/tangvis/erp/app/apirate/repository"
	"github.com/tangvis/erp/app/apirate/service/impl"
)

var ServiceSet = wire.NewSet(
	impl.NewLimiters,
	RepoSet,
)

var RepoSet = wire.NewSet(
	repository.NewRepoImpl,
)
