package common

import (
	"net/http"

	"github.com/tangvis/erp/libs/ecode"
)

var (
	ErrConfPanic = ecode.NewErrorConf(http.StatusInternalServerError)
	ErrPanicHTTP = ErrConfPanic.New("http")
)
