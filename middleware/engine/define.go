package engine

import "github.com/gin-gonic/gin"

const (
	startTimeKey = "__start_time"
)

type HTTPAPIJSONHandler func(ctx Context) (any, error)
type RawHandler func(ctx *gin.Context) error

type JSONResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Data     any    `json:"data"`
	TranceID string `json:"trance_id"`
}
