package engine

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"

	ctxUtil "github.com/tangvis/erp/libs/context"
	"github.com/tangvis/erp/libs/ecode"
	logutil "github.com/tangvis/erp/libs/log"
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
}

type HTTPEngine interface {
	JSON(handler HTTPAPIJSONHandler) gin.HandlerFunc
}

type HttpContext struct {
	*gin.Context

	Ctx context.Context
}

func NewHttpContext(ginCtx *gin.Context) *HttpContext {
	return &HttpContext{
		Context: ginCtx,
		Ctx:     context.Background(),
	}
}

func (c *HttpContext) GetTraceID() string {
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

func (engine *Engine) startRequest(ctx *HttpContext) {
	ctx.Set(startTimeKey, time.Now())
	ctx.Ctx = ctxUtil.AutoWrapContext(ctx.Ctx, ctx.GetTraceID())
	// 先写入response Header
	globalID := ctx.GetTraceID()
	ctx.Writer.Header().Add("x-trace-id", globalID) // 这个key是给前端用的
}

func (engine *Engine) beforeWriteBody(ctx *HttpContext) {
	if startTime := ctx.GetTime(startTimeKey); !startTime.IsZero() {
		duration := time.Since(startTime)
		ctx.Writer.Header().Add("x-response-et", strconv.FormatInt(duration.Milliseconds(), 10))
	}
}

func (engine *Engine) toResponse(ctx *HttpContext, data interface{}, err error) JSONResponse {
	resp := JSONResponse{
		TranceID: ctx.GetTraceID(),
		Data:     data,
		Message:  "Success",
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

func (engine *Engine) json(ctx *HttpContext, data interface{}, err error) {
	resp := engine.toResponse(ctx, data, err)
	engine.beforeWriteBody(ctx)
	ctx.PureJSON(200, resp)
}

func (engine *Engine) OpenAPIJSON(handler HTTPAPIJSONHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawHandler := func(ctx *HttpContext) error {
			resp, err := handler(ctx)
			engine.json(ctx, resp, err)
			return nil
		}
		engine.handleRaw(NewHttpContext(ctx), rawHandler)
	}
}

func (engine *Engine) String(ctx *HttpContext, code int, msg string) {
	engine.beforeWriteBody(ctx)
	ctx.String(code, msg)
}

func (engine *Engine) handleRaw(ctx *HttpContext, handler RawHandler) {
	engine.startRequest(ctx)
	// 业务处理
	defer func(start time.Time) {
		if ev := recover(); ev != nil {
			stack := make([]byte, 16*1024)
			runtime.Stack(stack, false)
			//log.Printf("%s", stack)
			logutil.CtxErrorF(ctx, "[PANIC]%+v, %s", ev, stack)
			engine.String(ctx, http.StatusInternalServerError, "panic")
		}
	}(time.Now())
	// 逻辑开始
	_ = handler(ctx) // todo error handle
}

func (engine *Engine) handleJSON(ctx *HttpContext, handler HTTPAPIJSONHandler) {
	rawHandler := func(ctx *HttpContext) error {
		resp, err := handler(ctx)
		engine.json(ctx, resp, err)
		return err
	}
	engine.handleRaw(ctx, rawHandler)
}

func (engine *Engine) JSON(handler HTTPAPIJSONHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		engine.handleJSON(NewHttpContext(ctx), handler)
	}
}

type Controller interface {
	URLPatterns() []Router
}

func NewRouter(method, url string, handler gin.HandlerFunc) Router {
	return Router{
		Method:  method,
		URL:     url,
		Handler: handler,
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
