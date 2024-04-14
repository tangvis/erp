package service

import (
	"context"
	"github.com/tangvis/erp/app/user/define"

	"github.com/tangvis/erp/app/user/repository"
)

type APP interface {
	GetUserByName(ctx context.Context, username string) (repository.UserTab, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (repository.UserTab, error)
	GetUserByEmail(ctx context.Context, email string) (repository.UserTab, error)
	GetUserByID(ctx context.Context, id uint64) (repository.UserTab, error)
	CreateUser(ctx context.Context, user define.UserEntity) (define.UserEntity, error)
}
