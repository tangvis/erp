//go:build wireinject
// +build wireinject

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/tangvis/erp/biz/ping"
	"github.com/tangvis/erp/biz/ping/access"
	"github.com/tangvis/erp/middleware/engine"
)

type application struct {
	PingController *access.Controller
}

func (app *application) GetRouterGroups() []engine.Controller {
	return []engine.Controller{
		app.PingController,
	}
}

func initializeApplication(dep *dependence) (*application, error) {
	wire.Build(
		ping.APISet,
		engine.EngineSet,
		wire.FieldsOf(
			new(*dependence),
			"DB",
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
				ginEngine.GET(router.URL, router.Handler)
			case http.MethodPost:
				ginEngine.POST(router.URL, router.Handler)
			default:
				return fmt.Errorf("unsupported http method %s", router.Method)
			}
		}
	}
	return nil
}