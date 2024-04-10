package impl

import (
	"context"
	"fmt"
	"github.com/tangvis/erp/app/user/define"

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
	return u.doQuery(ctx, define.UserQuery{
		Usernames: []string{username},
	})
}

func (u User) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (repository.UserTab, error) {
	return u.doQuery(ctx, define.UserQuery{
		PhoneNumbers: []string{phoneNumber},
	})
}

func (u User) GetUserByEmail(ctx context.Context, email string) (repository.UserTab, error) {
	return u.doQuery(ctx, define.UserQuery{
		Emails: []string{email},
	})
}

func (u User) doQuery(ctx context.Context, query define.UserQuery) (repository.UserTab, error) {
	users, err := u.repo.QueryUserByName(ctx, query)
	if err != nil {
		return repository.UserTab{}, err
	}
	if len(users) == 0 {
		return repository.UserTab{}, fmt.Errorf("no record found")
	}

	return users[0], nil
}

func (u User) GetUserByID(ctx context.Context, id uint64) (repository.UserTab, error) {
	return u.repo.GetUserByID(ctx, id)
}
