package impl

import (
	"context"
	"fmt"

	"github.com/tangvis/erp/app/user/repository"
	"github.com/tangvis/erp/app/user/service"
)

type User struct {
	repo repository.UserRepo
}

func NewUserAPP(
	repo repository.UserRepo,
) service.APP {
	return &User{
		repo: repo,
	}
}

func (u User) GetUserByName(ctx context.Context, username string) (repository.UserTab, error) {
	users, err := u.repo.GetUserByName(ctx, username)
	if err != nil {
		return repository.UserTab{}, err
	}
	if len(users) == 0 {
		return repository.UserTab{}, fmt.Errorf("no record found")
	}

	return users[0], nil
}

func (u User) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (repository.UserTab, error) {
	users, err := u.repo.GetUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return repository.UserTab{}, err
	}
	if len(users) == 0 {
		return repository.UserTab{}, fmt.Errorf("no record found")
	}

	return users[0], nil
}
