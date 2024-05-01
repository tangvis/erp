package product

import (
	"github.com/google/wire"

	"github.com/tangvis/erp/app/product/repository/meta"
	"github.com/tangvis/erp/app/product/service/impl"
)

var RepoSet = wire.NewSet(
	meta.NewRepoImpl,
)

var ServiceSet = wire.NewSet(
	RepoSet,
	impl.NewCategoryImpl,
	impl.NewBrandImpl,
)
