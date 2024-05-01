//go:build wireinject
// +build wireinject

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/tangvis/erp/access"
	pingHTTP "github.com/tangvis/erp/access/ping"
	"github.com/tangvis/erp/access/product"
	systemHTTP "github.com/tangvis/erp/access/system"
	userHTTP "github.com/tangvis/erp/access/user"
	"github.com/tangvis/erp/agent/email"
	"github.com/tangvis/erp/agent/templates"
	"github.com/tangvis/erp/app"
	"github.com/tangvis/erp/app/apirate/service"
	"github.com/tangvis/erp/app/system"
	getter "github.com/tangvis/erp/conf/config"
	"github.com/tangvis/erp/middleware"
	"github.com/tangvis/erp/middleware/engine"
	"net/http"
)

type application struct {
	pingController    *pingHTTP.Controller
	userController    *userHTTP.Controller
	productController *product.Controller
	systemController  *systemHTTP.Controller

	rateLimiterAPP service.APP
	sessionStore   engine.Store
}

func (app *application) GetRouterGroups() []engine.Controller {
	return []engine.Controller{
		app.pingController,
		app.userController,
		app.productController,
		app.systemController,
	}
}

func initializeApplication(
	dep *dependence,
) (*application, error) {
	wire.Build(
		middleware.Set,
		access.HTTPSet,
		app.ServiceSet,
		templates.NewDefaultTemplate,
		email.NewDefaultClient,
		system.Set,
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
