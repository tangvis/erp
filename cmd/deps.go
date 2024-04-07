package main

import (
	"github.com/tangvis/erp/agent/mysql"
	getter "github.com/tangvis/erp/conf/config"
	logutil "github.com/tangvis/erp/libs/log"
)

type dependence struct {
	DB *mysql.DB
}

func newDependence() (*dependence, error) {
	ret := &dependence{}
	ret.initDB()
	return ret, nil
}

func (d *dependence) initDB() {
	if getter.Config == nil {
		panic("config not init yet")
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
