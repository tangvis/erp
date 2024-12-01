package engine

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/tangvis/erp/conf/config"
	ctxUtil "github.com/tangvis/erp/pkg/context"
	logutil "github.com/tangvis/erp/pkg/log"
)

func startRequest(ctx *gin.Context) {
	ctx.Set(startTimeKey, time.Now())
	// 先写入response Header
	globalID := GetTraceID(ctx)
	ctx.Writer.Header().Add("x-trace-id", globalID) // 这个key是给前端用的
}

func PanicWrapper(c *gin.Context) {
	startRequest(c)
	// 业务处理
	defer func(start time.Time) {
		if ev := recover(); ev != nil {
			stack := make([]byte, 16*1024)
			runtime.Stack(stack, false)
			logutil.CtxErrorF(ctxUtil.AutoWrapContext(context.Background(), GetTraceID(c)), "[PANIC]%+v, %s", ev, stack)
			String(c, http.StatusInternalServerError, "panic")
		}
	}(time.Now())
	//time.Sleep(300 * time.Millisecond)
	c.Next()
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogWrapper(c *gin.Context) {
	if !config.Config.GetEnableLogRequest() {
		c.Next()
		return
	}
	logCtx := ctxUtil.AutoWrapContext(context.Background(), GetTraceID(c))
	if c.Request.Body != nil {
		body, _ := io.ReadAll(c.Request.Body)
		buff := bytes.NewBuffer(body)
		c.Request.Body = io.NopCloser(buff)
		logutil.CtxInfoF(logCtx, "------> path - %s, request body - <%s>", c.Request.URL.Path, bytes.NewBuffer(body).String())
	}
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	logutil.CtxInfoF(logCtx, "<------ path - %s, response body - <%s>", c.Request.URL.Path, blw.body.String())
}
