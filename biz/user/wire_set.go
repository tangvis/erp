package user

import (
	"github.com/google/wire"
	"github.com/tangvis/erp/biz/user/access"
	"github.com/tangvis/erp/biz/user/service/impl"
)

var APISet = wire.NewSet(
	BizSet,
	access.NewController,
)

var BizSet = wire.NewSet(
	impl.NewUserBiz,
)
