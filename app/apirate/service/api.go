package service

import "github.com/gin-gonic/gin"

type APP interface {
	InitPublic(publicLimitSetting map[string]int)
	RateLimitWrapper(c *gin.Context)
}
