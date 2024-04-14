package impl

import (
	"context"
	appDefine "github.com/tangvis/erp/app/user/define"
	userAPP "github.com/tangvis/erp/app/user/service"
	"github.com/tangvis/erp/biz/user/service"
	"github.com/tangvis/erp/biz/user/service/define"

	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/agent/redis"
)

type User struct {
	db    *mysql.DB
	cache redis.Cache
	app   userAPP.APP
}

func NewUserBiz(
	db *mysql.DB,
	cache redis.Cache,
	app userAPP.APP,
) service.APP {
	return &User{
		db:    db,
		cache: cache,
		app:   app,
	}
}

func (p *User) Create(ctx context.Context, req define.SignupRequest) (uint64, error) {
	createdUser, err := p.app.CreateUser(ctx, appDefine.UserEntity{
		Username: req.Username,
		Passwd:   req.Password,
		Email:    req.Email,
	})

	return createdUser.ID, err
}
