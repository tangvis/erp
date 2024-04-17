package ping

import (
	"net/http"

	"github.com/tangvis/erp/app/ping/service"
	"github.com/tangvis/erp/middleware/engine"
)

type Controller struct {
	engine engine.HTTPEngine
	biz    service.APP
}

func NewController(
	engine engine.HTTPEngine,
	app service.APP,
) *Controller {
	return &Controller{
		engine: engine,
		biz:    app,
	}
}

func (c *Controller) URLPatterns() []engine.Router {
	return []engine.Router{
		engine.NewRouter(http.MethodGet, "/ping", c.engine.JSONAuth(c.Ping)),
		engine.NewRouter(http.MethodPost, "/ping_failed", c.engine.JSON(c.Error)),
	}
}

func (c *Controller) Ping(ctx engine.Context, userInfo engine.UserInfo) (any, error) {
	return c.biz.Ping(), nil
}

func (c *Controller) Error(ctx engine.Context) (any, error) {
	var req FailPingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.biz.PingFail(ctx.GetCtx())
}

type FailPingRequest struct {
	ID   uint64 `json:"id" binding:"required"`
	Name string `json:"name"`
}
