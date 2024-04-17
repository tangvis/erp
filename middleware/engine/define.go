package engine

import "github.com/gin-gonic/gin"

const (
	startTimeKey = "__start_time"
	UserInfoKey  = "user_info"
)

type UserInfo struct {
	ID       uint64
	Username string
	Email    string
}

type HTTPAPIJSONHandler func(ctx Context) (any, error)
type HTTPAPIJSONUserHandler func(ctx Context, userInfo UserInfo) (any, error)
type RawHandler func(ctx *gin.Context) error

type JSONResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Data     any    `json:"data"`
	TranceID string `json:"trance_id"`
}
