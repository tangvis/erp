package engine

import (
	"context"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"time"

	"github.com/tangvis/erp/conf/config"
	ctxUtil "github.com/tangvis/erp/libs/context"
	"github.com/tangvis/erp/libs/ecode"
)

type Context interface {
	context.Context
	// ContentType 获取底层的
	ContentType() string
	// ShouldBind 读取参数
	ShouldBind(dest interface{}) error
	ShouldBindJSON(dest interface{}) error
	Data(code int, contentType string, data []byte)
	Header(key, value string)
	GetCtx() context.Context
}

type HTTPEngine interface {
	JSON(handler HTTPAPIJSONHandler) gin.HandlersChain
}

type HttpContext struct {
	*gin.Context

	Ctx context.Context
}

func (c *HttpContext) GetCtx() context.Context {
	return c.Ctx
}

func NewHttpContext(ginCtx *gin.Context) *HttpContext {
	return &HttpContext{
		Context: ginCtx,
		Ctx:     ctxUtil.AutoWrapContext(context.Background(), GetTraceID(ginCtx)),
	}
}

func GetTraceID(c *gin.Context) string {
	// Attempt to get the trace_id from the context
	traceID := c.GetString(ctxUtil.TraceIDKey)
	if len(traceID) > 0 {
		// If it exists, assert the type to string and return it
		return traceID
	}

	// If not found, generate a new trace ID
	newTraceID := ctxUtil.GenerateTrace()
	// Use the Set method to store the trace ID in the context
	c.Set(ctxUtil.TraceIDKey, newTraceID)
	// Return the new trace ID
	return newTraceID
}

type Engine struct {
}

func NewEngine() HTTPEngine {
	return &Engine{}
}

func beforeWriteBody(ctx *gin.Context) {
	if startTime := ctx.GetTime(startTimeKey); !startTime.IsZero() {
		duration := time.Since(startTime)
		ctx.Writer.Header().Add("x-response-et", strconv.FormatInt(duration.Milliseconds(), 10))
	}
}

func toResponse(ctx *gin.Context, data interface{}, err error) JSONResponse {
	resp := JSONResponse{
		Data:    data,
		Message: "Success",
	}
	if config.Config.GetEnableResponseTraceID() {
		resp.TranceID = GetTraceID(ctx)
	}
	// 如果是空数据不返回nil，而是返回一个空的map给前端
	if IsNilValue(data) {
		resp.Data = map[string]interface{}{}
	}
	// no error
	if err == nil {
		return resp
	}
	// error
	resp.Code = ecode.GetErrCode(err)
	resp.Message = err.Error()
	return resp
}

func json(ctx *gin.Context, data interface{}, err error) {
	resp := toResponse(ctx, data, err)
	beforeWriteBody(ctx)
	ctx.PureJSON(200, resp)
}

func String(ctx *gin.Context, code int, msg string) {
	beforeWriteBody(ctx)
	ctx.String(code, msg)
}

func (engine *Engine) JSON(handler HTTPAPIJSONHandler) gin.HandlersChain {
	coreHandler := func(ctx *gin.Context) {
		resp, err := handler(NewHttpContext(ctx))
		json(ctx, resp, err)
		if err != nil {
			_ = ctx.Error(err)
		}
	}
	return append(gin.HandlersChain{PanicWrapper, LogWrapper}, coreHandler)
}

type Controller interface {
	URLPatterns() []Router
}

func NewRouter(method, url string, handlers gin.HandlersChain) Router {
	if len(handlers) == 0 {
		panic("handlers is empty")
	}
	return Router{
		Method:   method,
		URL:      url,
		Handlers: handlers,
	}
}

func IsNilValue(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	switch kind := value.Kind(); kind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
