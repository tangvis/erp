package repository

import (
	"context"
	"errors"
	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/app/user/define"
	"gorm.io/gorm"
	"time"
)

type User interface {
	QueryUserByName(ctx context.Context, query define.UserQuery) ([]UserTab, error)
	GetUserByID(ctx context.Context, id uint64) (UserTab, error)
	CreateUser(ctx context.Context, user UserTab) (UserTab, error)
}

type UserRepo struct {
	db *mysql.DB
}

func (u *UserRepo) GetUserByID(ctx context.Context, id uint64) (UserTab, error) {
	var user UserTab
	// Using Gorm's First method to retrieve the first record that matches the query.
	// The method automatically adds a "LIMIT 1" to the query.
	if err := u.db.WithContext(ctx).Model(&UserTab{}).Where("id = ?", id).First(&user).Error; err != nil {
		// Handling the case where the user is not found or there's another error.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// todo handle no found error
			return UserTab{}, err
		}
		// Return the error if it's of a different type.
		return UserTab{}, err
	}
	return user, nil
}

func NewUserRepo(
	db *mysql.DB,
) User {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) QueryUserByName(ctx context.Context, query define.UserQuery) ([]UserTab, error) {
	if err := query.Valid(); err != nil {
		return nil, err
	}
	q := u.db.WithContext(ctx).Model(&UserTab{})
	users := make([]UserTab, 0)
	// Use Gorm's `Where` with `IN` clause for matching usernames
	if len(query.Usernames) > 0 {
		q = q.Where("username IN ?", query.Usernames)
	}
	if len(query.PhoneNumbers) > 0 {
		q = q.Where("phone_number IN ?", query.PhoneNumbers)
	}
	if len(query.Emails) > 0 {
		q = q.Where("email IN ?", query.Emails)
	}
	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserRepo) CreateUser(ctx context.Context, user UserTab) (UserTab, error) {
	now := time.Now()
	user.Ctime = now.UnixMilli()
	user.Mtime = now.UnixMilli()
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		return UserTab{}, err
	}
	return user, nil
}
