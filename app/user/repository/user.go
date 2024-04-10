package repository

import (
	"context"

	"github.com/tangvis/erp/agent/mysql"
)

type User interface {
	GetUserByName(ctx context.Context, username ...string) ([]UserTab, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber ...string) ([]UserTab, error)
}

type UserRepo struct {
	db *mysql.DB
}

func NewUserRepo(
	db *mysql.DB,
) User {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) GetUserByName(ctx context.Context, usernames ...string) ([]UserTab, error) {
	users := make([]UserTab, 0)
	// Use Gorm's `Where` with `IN` clause for matching usernames
	if err := u.db.WithContext(ctx).Where("username IN ?", usernames).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserRepo) GetUserByPhoneNumber(ctx context.Context, phoneNumbers ...string) ([]UserTab, error) {
	users := make([]UserTab, 0)
	// Use Gorm's `Where` with `IN` clause for matching phone numbers
	if err := u.db.WithContext(ctx).Where("phone_number IN ?", phoneNumbers).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
