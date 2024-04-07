package common

import (
	"github.com/tangvis/erp/libs/ecode"
)

var (
	// ErrConfPing ping
	ErrConfPing       = ecode.NewBusinessErrorCode(ecode.Ping, 1)
	ErrPingFailed     = ecode.NewErrorConf(ErrConfPing)
	ErrPingFailedTest = ErrPingFailed.New("ping failed test")
)
