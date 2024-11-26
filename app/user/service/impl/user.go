package impl

import (
	"context"
	"errors"
	"time"

	email "github.com/tangvis/erp/app/system/email/service"
	"github.com/tangvis/erp/app/user/define"
	"github.com/tangvis/erp/app/user/repository"
	"github.com/tangvis/erp/app/user/service"
	"github.com/tangvis/erp/common"
	"github.com/tangvis/erp/middleware/engine"
)

type User struct {
	repo         repository.User
	sessionStore engine.Store
	emailSrv     email.APP
}

func NewUserAPP(
	repo repository.User,
	sessionStore engine.Store,
	emailSrv email.APP,
) service.APP {
	return &User{
		repo:         repo,
		sessionStore: sessionStore,
		emailSrv:     emailSrv,
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
		return repository.UserTab{}, common.ErrDBRecordNotFound
	}

	return users[0], nil
}

func (u User) GetUserByID(ctx context.Context, id uint64) (repository.UserTab, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u User) CreateUser(ctx context.Context, user define.UserEntity) (define.UserEntity, error) {
	if err := u.checkInfoAvailable(ctx, user); err != nil {
		return define.UserEntity{}, err
	}
	createdUser, err := u.repo.SaveUser(ctx, repository.UserTab{
		Username:    user.Username,
		Passwd:      user.Passwd,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		UserStatus:  user.Status,
	})
	if err != nil {
		return define.UserEntity{}, err
	}

	user.ID = createdUser.ID
	user.Ctime = createdUser.Ctime
	user.Mtime = createdUser.Mtime

	return user, nil
}

func (u User) Login(ctx context.Context, req define.LoginRequest) (define.UserEntity, error) {
	if req.Username != "" {
		return u.login(ctx, req.Username, req.Password, u.GetUserByName)
	}
	return u.login(ctx, req.Email, req.Password, u.GetUserByEmail)
}

func (u User) login(ctx context.Context, info, passwd string, getUser func(ctx context.Context, email string) (repository.UserTab, error)) (define.UserEntity, error) {
	user, err := getUser(ctx, info)
	if err != nil {
		if errors.Is(err, common.ErrDBRecordNotFound) {
			return define.UserEntity{}, common.ErrUserInfo
		}
		return define.UserEntity{}, err
	}
	if user.Passwd != passwd {
		return define.UserEntity{}, common.ErrUserInfo
	}
	onlineUsers, err := u.sessionStore.OnlineUsers(ctx, user.ID)
	if err != nil {
		return define.UserEntity{}, err
	}
	loginTime := time.Now()
	resp := define.UserEntity{
		ID:          user.ID,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		LoginTime:   loginTime.Unix(),
	}
	if len(onlineUsers) >= define.MaxOnlineForAUser {
		return resp, common.ErrUserTooManyLogin
	}
	go func() {
		user.LatestLoginTime = loginTime.Unix()
		_, err = u.repo.SaveUser(ctx, user)
	}()
	return resp, err
}

func (u User) OnlineUsers(ctx context.Context, id uint64) ([]define.UserEntity, error) {
	onlineUsers, err := u.sessionStore.OnlineUsers(ctx, id)
	if err != nil {
		return nil, err
	}
	result := make([]define.UserEntity, 0, len(onlineUsers))
	for _, user := range onlineUsers {
		result = append(result, define.UserEntity{
			ID:          user.ID,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			LoginTime:   user.LoginTime,
			SessionID:   user.SessionID,
		})
	}

	return result, nil
}

func (u User) checkInfoAvailable(ctx context.Context, user define.UserEntity) error {
	// Check if the username already exists
	existingUser, err := u.GetUserByName(ctx, user.Username)
	if err != nil && !errors.Is(err, common.ErrDBRecordNotFound) {
		return err
	}
	if existingUser.ID > 0 {
		return common.ErrUser.New("username exists")
	}

	// Check if the email already exists
	existingUser, err = u.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, common.ErrDBRecordNotFound) {
		return err
	}
	if existingUser.ID > 0 {
		return common.ErrUser.New("email exists")
	}

	// If neither username nor email already exists, return nil (no error)
	return nil
}
