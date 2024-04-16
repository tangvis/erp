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
		engine.NewRouter(http.MethodPost, "/user/signup", c.engine.JSON(c.Signup)),
		engine.NewRouter(http.MethodPost, "/user/login", c.engine.JSON(c.Login)),
	}
}

func (c *Controller) Signup(ctx engine.Context) (any, error) {
	var req define.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	createdUser, err := c.app.CreateUser(ctx, define.UserEntity{
		Username: req.Username,
		Passwd:   req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (c *Controller) Login(ctx engine.Context) (any, error) {
	var req define.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	createdUser, err := c.app.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
