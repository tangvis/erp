package category

import (
	"net/http"

	"github.com/tangvis/erp/app/product/define"
	"github.com/tangvis/erp/app/product/service"
	"github.com/tangvis/erp/common"
	"github.com/tangvis/erp/middleware/engine"
)

type Controller struct {
	engine engine.HTTPEngine
	app    service.Category
}

func NewController(
	engine engine.HTTPEngine,
	app service.Category,
) *Controller {
	return &Controller{
		engine: engine,
		app:    app,
	}
}

func (c *Controller) URLPatterns() []engine.Router {
	return []engine.Router{
		engine.NewRouter(http.MethodPost, "/category/add", c.engine.JSONAuth(c.Add)),
		engine.NewRouter(http.MethodPost, "/category/update", c.engine.JSONAuth(c.Update)),
		engine.NewRouter(http.MethodGet, "/category/list", c.engine.JSONAuth(c.List)),
		engine.NewRouter(http.MethodPost, "/category/delete", c.engine.JSONAuth(c.Remove)),
	}
}

func (c *Controller) Add(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.AddCateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.app.Add(ctx.GetCtx(), userInfo, &req)
}

func (c *Controller) List(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	return c.app.List(ctx.GetCtx(), userInfo)
}

func (c *Controller) Update(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.UpdateCateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.app.Update(ctx.GetCtx(), userInfo, &req)
}

func (c *Controller) Remove(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.RemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return nil, c.app.Remove(ctx.GetCtx(), userInfo, req.IDs...)
}
