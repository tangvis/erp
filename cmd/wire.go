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
	"github.com/tangvis/erp/biz/ping"
	"github.com/tangvis/erp/biz/ping/access"
	"github.com/tangvis/erp/middleware/engine"
)

type application struct {
	pingController *access.Controller

	rateLimiterAPP service.APP
}

func (app *application) GetRouterGroups() []engine.Controller {
	return []engine.Controller{
		app.pingController,
	}
}

func (app *application) Use(g *gin.Engine) {
	app.rateLimiterAPP.InitPublic(map[string]int{})
	g.Use(app.rateLimiterAPP.RateLimitWrapper)
}

func initializeApplication(
	dep *dependence,
) (*application, error) {
	wire.Build(
		ping.APISet,
		engine.Set,
		apirate.ServiceSet,
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
	app.Use(ginEngine)
	return nil
}
