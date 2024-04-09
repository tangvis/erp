package main

import "github.com/gin-gonic/gin"

func main() {
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	initGlobalResources()
	dep, err := newDependence()
	if err != nil {
		panic(err)
	}
	app, err := initializeApplication(dep)
	if err != nil {
		panic(err)
	}
	ginEngine.Use()
	if err = app.registerHTTP(ginEngine); err != nil {
		panic(err)
	}
	// nolint:errcheck
	ginEngine.Run()
}
