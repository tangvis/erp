package common

import "github.com/gin-gonic/gin"

type Router struct {
	Method   string
	URL      string
	Handlers gin.HandlersChain
}

var AllRouters []Router

func RegisterAllRouters(all []Router) {
	AllRouters = all
}
