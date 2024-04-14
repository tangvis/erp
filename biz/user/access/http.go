package access

import (
	"github.com/tangvis/erp/biz/user/service"
	"github.com/tangvis/erp/biz/user/service/define"
	"net/http"

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
		engine.NewRouter(http.MethodPost, "/user/signup", c.engine.JSON(c.Create)),
	}
}

func (c *Controller) Create(ctx engine.Context) (any, error) {
	var req define.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.biz.Create(ctx, req)
}
