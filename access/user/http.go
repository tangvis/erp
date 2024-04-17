package user

import (
	"net/http"
	"time"

	"github.com/tangvis/erp/libs/crypto"

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
		engine.NewRouter(http.MethodPost, "/user/logout", c.engine.JSON(c.SignOut)),
	}
}

func (c *Controller) Signup(ctx engine.Context) (any, error) {
	var req define.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	createdUser, err := c.app.CreateUser(ctx, define.UserEntity{
		Username: req.Username,
		Passwd:   crypto.GetMD5Hash(req.Password),
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
	req.Password = crypto.GetMD5Hash(req.Password)
	userInfo, err := c.app.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	ctx.SetSession(engine.UserInfo{
		ID:        userInfo.ID,
		Username:  userInfo.Username,
		Email:     userInfo.Email,
		LoginTime: time.Now().Unix(),
	})

	return userInfo, nil
}

func (c *Controller) SignOut(ctx engine.Context) (any, error) {
	return nil, ctx.SignOut()
}
