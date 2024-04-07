package main

import (
	"github.com/spf13/viper"
	"os"

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
	if config == nil {
		panic("config not init yet")
	}
	dbConfig, err := config.GetMySQLConfig()
	if err != nil {
		panic(err)
	}
	db, err := mysql.NewMySQL(dbConfig)
	if err != nil {
		panic(err)
	}
	d.DB = db
}

var config getter.Getter

func initConfig() {
	vp := initViper()
	config = getter.NewConfigGetter(vp)
}

func initViper() *viper.Viper {
	env := os.Getenv(getter.EnvKey)
	vp := viper.New()
	vp.SetConfigName("app_" + env)
	vp.AddConfigPath("./conf/")
	vp.SetConfigType("toml")
	if err := vp.ReadInConfig(); err != nil {
		//if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		//	panic(err)
		//} else {
		//	panic(err)
		//}
		panic(err)
	}
	vp.WatchConfig()

	return vp
}

func initLogger() {
	logConfig := logutil.NewConfig()
	logConfig.DisableJSONFormat()
	logConfig.SetFileOut("./logs", "test_log", 1, 2)
	logutil.InitLogger(logConfig)
}

func initGlobalResources() {
	initConfig()
	initLogger()
}
