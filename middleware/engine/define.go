package engine

import (
	jsonLib "encoding/json"
	"github.com/gin-gonic/gin"
)

const (
	startTimeKey = "__start_time"
	UserInfoKey  = "user_info"
)

type Config struct {
	ResponseTraceID bool `toml:"response_trace_id"`
	RequestLog      bool `toml:"log_request"`
	PublicQpsLimit  int  `toml:"public_qps_limit"`
}

type UserInfo struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	LoginTime   int64  `json:"login_time"`
	IP          string `json:"ip"`
}

func (u *UserInfo) String() string {
	b, _ := jsonLib.Marshal(u)

	return string(b)
}

type HTTPAPIJSONHandler func(ctx Context) (any, error)
type HTTPAPIJSONUserHandler func(ctx Context, userInfo *UserInfo) (any, error)
type RawHandler func(ctx *gin.Context) error

type JSONResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Data     any    `json:"data"`
	TranceID string `json:"trance_id"`
}
