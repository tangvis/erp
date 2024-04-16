package engine

import (
	"bytes"
	"context"
	jsonLib "encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
	"time"

	"github.com/tangvis/erp/common"
	"github.com/tangvis/erp/conf/config"
	ctxUtil "github.com/tangvis/erp/libs/context"
	"github.com/tangvis/erp/libs/ecode"
)

type Context interface {
	context.Context
	// ContentType 获取底层的
	ContentType() string
	// ShouldBind 读取参数
	ShouldBind(dest any) error
	ShouldBindJSON(dest any) error
	Data(code int, contentType string, data []byte)
	Header(key, value string)
	GetCtx() context.Context
}

type HTTPEngine interface {
	JSON(handler HTTPAPIJSONHandler) gin.HandlersChain
}

type HttpContext struct {
	ginCtx *gin.Context

	Ctx context.Context
}

func (c *HttpContext) Deadline() (deadline time.Time, ok bool) {
	return c.ginCtx.Deadline()
}

func (c *HttpContext) Done() <-chan struct{} {
	return c.ginCtx.Done()
}

func (c *HttpContext) Err() error {
	return c.ginCtx.Err()
}

func (c *HttpContext) Value(key any) any {
	return c.ginCtx.Value(key)
}

func (c *HttpContext) ContentType() string {
	//TODO implement me
	panic("implement me")
}

func (c *HttpContext) ShouldBind(dest any) error {
	if err := c.ginCtx.ShouldBind(dest); err != nil {
		return c.convertParamError(err)
	}
	return nil
}

func (c *HttpContext) Data(code int, contentType string, data []byte) {
	c.ginCtx.Data(code, contentType, data)
}

func (c *HttpContext) Header(key, value string) {
	c.ginCtx.Header(key, value)
}

func (c *HttpContext) ShouldBindJSON(dest any) error {
	if err := c.ginCtx.ShouldBindJSON(dest); err != nil {
		return c.convertParamError(err)
	}
	if customValidator, ok := dest.(Validator); ok {
		if err := customValidator.Validate(); err != nil {
			return common.ErrConfInvalidArguments.New(err.Error())
		}
	}
	return nil
}

func (c *HttpContext) convertParamError(err error) error {
	var (
		typeError        *jsonLib.UnmarshalTypeError
		syntaxError      *jsonLib.SyntaxError
		validationErrors validator.ValidationErrors
	)

	switch {
	case errors.As(err, &typeError):
		return common.ErrConfInvalidArguments.NewF("mismatch proto: %s", typeError.Error())
	case errors.As(err, &syntaxError):
		return common.ErrConfInvalidArguments.NewF("invalid json: %s", err)
	case errors.As(err, &validationErrors):
		buf := bytes.NewBuffer(nil)
		buf.WriteString("request params failed with validate: ")
		for _, v := range validationErrors {
			buf.WriteString(fmt.Sprintf("field %s failed with %s,", v.Field(), v.Tag()))
		}
		return common.ErrConfInvalidArguments.New(buf.String())
	default:
		return common.ErrConfInvalidArguments.New(err.Error())
	}
}

func (c *HttpContext) GetCtx() context.Context {
	return c.Ctx
}

func NewHttpContext(ginCtx *gin.Context) Context {
	return &HttpContext{
		ginCtx: ginCtx,
		Ctx:    ctxUtil.AutoWrapContext(context.Background(), GetTraceID(ginCtx)),
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
		ctx.Writer.Header().Add("x-response-et", fmt.Sprintf("%.2f", float64(duration.Microseconds())/1000))
	}
}

func toResponse(ctx *gin.Context, data any, err error) JSONResponse {
	resp := JSONResponse{
		Data:    data,
		Message: "Success",
	}
	if config.Config.GetEnableResponseTraceID() {
		resp.TranceID = GetTraceID(ctx)
	}
	// 如果是空数据不返回nil，而是返回一个空的map给前端
	if IsNilValue(data) {
		resp.Data = map[string]any{}
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

func json(ctx *gin.Context, data any, err error) {
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

func NewRouter(method, path string, handlers gin.HandlersChain) Router {
	if len(handlers) == 0 {
		panic("handlers is empty")
	}
	return Router{
		Method:   method,
		Path:     path,
		Handlers: handlers,
	}
}

type Router struct {
	Method   string
	Path     string
	Handlers gin.HandlersChain
}

func IsNilValue(object any) bool {
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

type Validator interface {
	Validate() error
}
