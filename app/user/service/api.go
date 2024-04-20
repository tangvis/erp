package service

import (
	"context"
	"github.com/tangvis/erp/app/user/define"

	"github.com/tangvis/erp/app/user/repository"
)

type APP interface {
	GetUserByID(ctx context.Context, id uint64) (repository.UserTab, error)
	CreateUser(ctx context.Context, user define.UserEntity) (define.UserEntity, error)
	Login(ctx context.Context, req define.LoginRequest) (define.UserEntity, error)
	OnlineUsers(ctx context.Context) ([]define.UserEntity, error)
}
