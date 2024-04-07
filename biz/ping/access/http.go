package access

import (
	"net/http"

	"github.com/tangvis/erp/biz/ping/service"
	"github.com/tangvis/erp/middleware/engine"
)

type Controller struct {
	engine engine.HTTPEngine
	app    service.APP
}

func NewController(
	engine engine.HTTPEngine,
	app service.APP,
) *Controller {
	return &Controller{
		engine: engine,
		app:    app,
	}
}

func (c *Controller) URLPatterns() []engine.Router {
	return []engine.Router{
		engine.NewRouter(http.MethodGet, "/ping", c.engine.JSON(c.Ping)),
	}
}

func (c *Controller) Ping(ctx engine.Context) (interface{}, error) {
	panic("")
}
