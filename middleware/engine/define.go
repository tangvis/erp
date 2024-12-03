package engine

import (
	"github.com/gin-gonic/gin"

	"github.com/tangvis/erp/common"
)

const (
	startTimeKey = "__start_time"
)

type HTTPAPIJSONHandler func(ctx Context) (any, error)
type HTTPAPIJSONUserHandler func(ctx Context, userInfo *common.UserInfo) (any, error)
type RawHandler func(ctx *gin.Context) error

type JSONResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}
