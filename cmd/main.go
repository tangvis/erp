package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	// Custom CORS configuration that allows all origins
	ginEngine.Use(disableAllCORSPolicy())
	initGlobalResources()
	dep, err := newDependence()
	if err != nil {
		panic(err)
	}
	app, err := initializeApplication(dep)
	if err != nil {
		panic(err)
	}
	if err = app.registerHTTP(ginEngine, dep); err != nil {
		panic(err)
	}
	_ = ginEngine.Run("0.0.0.0:8081")
}
