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
)
