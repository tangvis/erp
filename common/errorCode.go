package common

import (
	"github.com/tangvis/erp/libs/ecode"
)

var (
	// ErrConfInvalidArguments 参数错误
	ErrConfInvalidArguments = ecode.NewErrorConf(-10)

	// ErrConfPing ping
	ErrConfPing       = ecode.NewBusinessErrorCode(ecode.Ping, 1)
	ErrPingFailed     = ecode.NewErrorConf(ErrConfPing)
	ErrPingFailedTest = ErrPingFailed.New("ping failed test")

	ErrUser     = ecode.NewErrorConf(ecode.NewBusinessErrorCode(ecode.User, 1))
	ErrUserInfo = ErrUser.New("username or password error")
	ErrAuth     = ErrUser.New("auth failed")

	ErrDB               = ecode.NewSystemErrorCode(ecode.SystemDB, 10)
	ErrDBConf           = ecode.NewErrorConf(ErrDB)
	ErrDBRecordNotFound = ErrDBConf.New("record not found")
)
