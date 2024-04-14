//go:build wireinject
// +build wireinject

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/tangvis/erp/app/apirate"
	"github.com/tangvis/erp/app/apirate/service"
	userAPP "github.com/tangvis/erp/app/user"
	"github.com/tangvis/erp/biz/ping"
	pingHTTP "github.com/tangvis/erp/biz/ping/access"
	"github.com/tangvis/erp/biz/user"
	userHTTP "github.com/tangvis/erp/biz/user/access"
	"github.com/tangvis/erp/middleware/engine"
)

type application struct {
	pingController *pingHTTP.Controller
	userController *userHTTP.Controller

	rateLimiterAPP service.APP
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
		ping.APISet,
		engine.Set,
		apirate.ServiceSet,
		user.APISet,
		userAPP.ServiceSet,
		wire.FieldsOf(
			new(*dependence),
			"DB",
			"Cache",
		),
		wire.Struct(new(application), "*"),
	)
	return &application{}, nil
}

func (app *application) registerHTTP(ginEngine *gin.Engine) error {
	ginEngine.Use(app.rateLimiterAPP.RateLimitWrapper)
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

func (app *application) InitCommonRateLimiter(g *gin.Engine) {
	m := make(map[string]int)
	for _, route := range g.Routes() {
		m[route.Path] = 1
	}
	app.rateLimiterAPP.InitPublic(m)
}
