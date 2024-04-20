package main

import (
	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/agent/redis"
	getter "github.com/tangvis/erp/conf/config"
	logutil "github.com/tangvis/erp/pkg/log"
)

type dependence struct {
	DB    *mysql.DB
	Cache redis.Cache
}

func newDependence() (*dependence, error) {
	ret := &dependence{}
	ret.initDB()
	ret.initCache()
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
