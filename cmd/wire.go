//go:build wireinject
// +build wireinject

package main

import (
	"fmt"
	"github.com/tangvis/erp/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/tangvis/erp/access"
	pingHTTP "github.com/tangvis/erp/access/ping"
	userHTTP "github.com/tangvis/erp/access/user"
	"github.com/tangvis/erp/app/apirate"
	"github.com/tangvis/erp/app/apirate/service"
	"github.com/tangvis/erp/app/ping"
	userAPP "github.com/tangvis/erp/app/user"
	getter "github.com/tangvis/erp/conf/config"
	"github.com/tangvis/erp/middleware/engine"
)

type application struct {
	pingController *pingHTTP.Controller
	userController *userHTTP.Controller

	rateLimiterAPP service.APP
	sessionStore   engine.Store
}

func (app *application) GetRouterGroups() []engine.Controller {
	return []engine.Controller{
		app.pingController,
		app.userController,
	}
}

func initializeApplication(
	dep *dependence,
) (*application, error) {
	wire.Build(
		ping.ServiceSet,
		middleware.Set,
		apirate.ServiceSet,
		userAPP.ServiceSet,
		access.HTTPSet,
		wire.FieldsOf(
			new(*dependence),
			"DB",
			"Cache",
		),
		wire.Struct(new(application), "*"),
	)
	return &application{}, nil
}

func (app *application) registerHTTP(ginEngine *gin.Engine, dep *dependence) error {
	app.userMiddlewares(ginEngine)
	controllers := app.GetRouterGroups()
	for _, v := range controllers {
		for _, router := range v.URLPatterns() {
			switch router.Method {
			case http.MethodGet:
				ginEngine.GET(router.Path, router.Handlers...)
			case http.MethodPost:
				ginEngine.POST(router.Path, router.Handlers...)
			default:
				return fmt.Errorf("unsupported http method %s", router.Method)
			}
		}
	}
	app.InitCommonRateLimiter(ginEngine)
	return nil
}

func (app *application) userMiddlewares(ginEngine *gin.Engine) {
	ginEngine.Use(
		app.sessionStore.SessionHandler(),
		app.rateLimiterAPP.RateLimitWrapper,
		engine.PanicWrapper,
		engine.LogWrapper,
	)
}

func (app *application) InitCommonRateLimiter(g *gin.Engine) {
	cfg, err := getter.Config.GetMiddleWareConfig()
	if err != nil {
		panic(err)
	}
	m := make(map[string]int)
	for _, route := range g.Routes() {
		m[route.Path] = cfg.PublicQpsLimit
	}
	app.rateLimiterAPP.InitPublic(m)
}
