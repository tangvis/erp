package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/agent/openai"
	"github.com/tangvis/erp/agent/redis"
	getter "github.com/tangvis/erp/conf/config"
	logutil "github.com/tangvis/erp/pkg/log"
)

type dependence struct {
	DB     *mysql.DB
	Cache  redis.Cache
	OpenAI openai.API
}

func newDependence() (*dependence, error) {
	ret := &dependence{}
	ret.initDB()
	ret.initCache()
	ret.initOpenAI()
	return ret, nil
}

func (d *dependence) initDB() {
	if getter.Config == nil {
		panic("initDB config not init yet")
	}
	dbConfig, err := getter.Config.GetMySQLConfig()
	if err != nil {
		panic(err)
	}
	db, err := mysql.NewMySQL(dbConfig)
	if err != nil {
		panic(err)
	}
	d.DB = db
}

func (d *dependence) initCache() {
	if getter.Config == nil {
		panic("initCache config not init yet")
	}
	cacheConfig, err := getter.Config.GetCacheConfig()
	if err != nil {
		panic(err)
	}
	d.Cache = redis.NewCache(cacheConfig)
}

func (d *dependence) initOpenAI() {
	if getter.Config == nil {
		panic("initOpenAI config not init yet")
	}
	d.OpenAI = openai.NewGPTClient(getter.Config.GetOpenAIToken())
}

func initLogger() {
	logConfig := logutil.NewConfig()
	logConfig.DisableJSONFormat()
	logConfig.SetFileOut("./logs", "test_log", 1, 2)
	logutil.InitLogger(logConfig)
}

func initGlobalResources() {
	getter.InitConfig()
	initLogger()
}

func disableAllCORSPolicy() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: false, // Must be false when using AllowCredentials
		AllowOriginFunc: func(origin string) bool {
			return true // Allow all origins
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
