package main

import "github.com/gin-gonic/gin"

func main() {
	ginEngine := gin.New()
	gin.SetMode(gin.ReleaseMode)
	initGlobalResources()
	dep, err := newDependence()
	if err != nil {
		panic(err)
	}
	app, err := initializeApplication(dep)
	if err != nil {
		panic(err)
	}
	if err = app.registerHTTP(ginEngine); err != nil {
		panic(err)
	}
	ginEngine.Run()
}