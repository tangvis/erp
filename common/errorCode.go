package common

import (
	"github.com/tangvis/erp/pkg/ecode"
)

var (
	// ErrConfInvalidArguments 参数错误
	ErrConfInvalidArguments = ecode.NewErrorConf(-10)

	// ErrConfPing ping
	ErrConfPing       = ecode.NewSystemErrorCode(ecode.Ping, 1)
	ErrPingFailed     = ecode.NewErrorConf(ErrConfPing)
	ErrPingFailedTest = ErrPingFailed.New("ping failed test")

	ErrUser             = ecode.NewErrorConf(ecode.NewSystemErrorCode(ecode.User, 2))
	ErrUserInfo         = ErrUser.New("username or password error")
	ErrUserTooManyLogin = ErrUser.New("too many login, please logout other sessions then try again")
	ErrAuth             = ecode.NewErrorConf(ecode.NewSystemErrorCode(ecode.User, 1)).New("auth failed")

	ErrDB               = ecode.NewSystemErrorCode(ecode.SystemDB, 10)
	ErrDBConf           = ecode.NewErrorConf(ErrDB)
	ErrDBRecordNotFound = ErrDBConf.New("record not found")

	ErrCategory                = ecode.NewErrorConf(ecode.NewBusinessErrorCode(ecode.BusinessCategory, 1))
	ErrCategoryParentNotExists = ErrCategory.New("parent product not exists")
	ErrCategoryNotExists       = ErrCategory.New("product not exists")
	ErrCategoryNameConflict    = ErrCategory.New("name conflict")
	ErrCategoryHasChildren     = ErrCategory.New("current product has children, can't remove it")
	ErrBrand                   = ecode.NewErrorConf(ecode.NewBusinessErrorCode(ecode.BusinessBrand, 1))
	ErrBrandNameConflict       = ErrCategory.New("name conflict")
	ErrBrandNotExists          = ErrCategory.New("brand not exists")
)
