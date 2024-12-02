package system

import (
	"github.com/tangvis/erp/app/system/actionlog/define"
	actionLogAPP "github.com/tangvis/erp/app/system/actionlog/service"
	"github.com/tangvis/erp/common"
	"github.com/tangvis/erp/middleware/engine"

	"net/http"
)

type Controller struct {
	engine    engine.HTTPEngine
	actionLog actionLogAPP.APP
}

func NewController(
	engine engine.HTTPEngine,
	actionLog actionLogAPP.APP,
) *Controller {
	return &Controller{
		engine:    engine,
		actionLog: actionLog,
	}
}

func (c *Controller) URLPatterns() []engine.Router {
	return []engine.Router{
		// action log
		engine.NewRouter(http.MethodPost, "/action_log/list", c.engine.JSONAuth(c.ActionLogList)),
	}
}

func (c *Controller) ActionLogList(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.ListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.actionLog.List(ctx.GetCtx(), &req)
}
