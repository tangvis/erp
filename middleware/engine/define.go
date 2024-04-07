package engine

import "github.com/gin-gonic/gin"

const (
	startTimeKey = "__start_time"
	TraceIDKey   = "__trace_id"
)

type HTTPAPIJSONHandler func(ctx Context) (interface{}, error)
type GinHandler func(ctx *gin.Context)
type RawHandler func(ctx *HttpContext) error

type JSONResponse struct {
	TranceID string      `json:"trance_id"`
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
}

type Router struct {
	Method  string
	URL     string
	Handler GinHandler
}
