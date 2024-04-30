package common

import (
	"github.com/tangvis/erp/pkg/ecode"
)

var (
	// ErrConfInvalidArguments 参数错误
	ErrConfInvalidArguments = ecode.NewErrorConf(-10)

	// ErrConfPing ping
	ErrConfPing       = ecode.NewBusinessErrorCode(ecode.Ping, 1)
	ErrPingFailed     = ecode.NewErrorConf(ErrConfPing)
	ErrPingFailedTest = ErrPingFailed.New("ping failed test")

	ErrUser             = ecode.NewErrorConf(ecode.NewBusinessErrorCode(ecode.User, 1))
	ErrUserInfo         = ErrUser.New("username or password error")
	ErrUserTooManyLogin = ErrUser.New("too many login, please logout other sessions then try again")
	ErrAuth             = ErrUser.New("auth failed")

	ErrDB               = ecode.NewSystemErrorCode(ecode.SystemDB, 10)
	ErrDBConf           = ecode.NewErrorConf(ErrDB)
	ErrDBRecordNotFound = ErrDBConf.New("record not found")

	ErrCategory                = ecode.NewErrorConf(ecode.NewBusinessErrorCode(ecode.BusinessCategory, 1))
	ErrCategoryParentNotExists = ErrCategory.New("parent category not exists")
	ErrCategoryNotExists       = ErrCategory.New("category not exists")
	ErrCategoryNameConflict    = ErrCategory.New("name conflict")
	ErrCategoryHasChildren     = ErrCategory.New("current category has children, can't remove it")
	ErrBrand                   = ecode.NewErrorConf(ecode.NewBusinessErrorCode(ecode.BusinessBrand, 1))
	ErrBrandNameConflict       = ErrCategory.New("name conflict")
	ErrBrandNotExists          = ErrCategory.New("brand not exists")
)
