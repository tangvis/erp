package logutil

import (
	"context"
	"go.uber.org/zap"
	"testing"

	ctxUtil "github.com/tangvis/erp/libs/context"
)

func TestInfo(t *testing.T) {
	config := NewConfig()
	config.DisableJSONFormat()
	config.SetFileOut("../../logs", "test_log", 1, 2)
	InitLogger(config)
	CtxInfo(ctxUtil.AutoWrapContext(context.Background(), ctxUtil.GenerateTrace()), "abc", zap.String("kk", "vv"))
}
