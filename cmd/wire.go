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
	"github.com/tangvis/erp/common"
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

func initializeApplication(
	dep *dependence,
) (*application, error) {
	wire.Build(
		ping.APISet,
		engine.Set,
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
	allRouters := make([]common.Router, 0)
	for _, v := range controllers {
		for _, router := range v.URLPatterns() {
			allRouters = append(allRouters, router)
			switch router.Method {
			case http.MethodGet:
				ginEngine.GET(router.URL, router.Handlers...)
			case http.MethodPost:
				ginEngine.POST(router.URL, router.Handlers...)
			default:
				return fmt.Errorf("unsupported http method %s", router.Method)
			}
		}
	}
	common.RegisterAllRouters(allRouters)
	return nil
}
