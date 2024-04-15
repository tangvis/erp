package user

import (
	"net/http"

	"github.com/tangvis/erp/app/user/define"
	userAPP "github.com/tangvis/erp/app/user/service"
	"github.com/tangvis/erp/middleware/engine"
)

type Controller struct {
	engine engine.HTTPEngine
	app    userAPP.APP
}

func NewController(
	engine engine.HTTPEngine,
	app userAPP.APP,
) *Controller {
	return &Controller{
		engine: engine,
		app:    app,
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
	return c.app.CreateUser(ctx, define.UserEntity{
		Username: req.Username,
		Passwd:   req.Password,
		Email:    req.Email,
	})
}
