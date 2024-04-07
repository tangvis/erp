package main

import (
	"github.com/spf13/viper"
	"os"

	"github.com/tangvis/erp/agent/mysql"
	getter "github.com/tangvis/erp/conf/config"
	logutil "github.com/tangvis/erp/libs/log"
)

var config getter.Getter

func initConfig() {
	vp := initViper()
	config = getter.NewConfigGetter(vp)
}

func initViper() *viper.Viper {
	env := os.Getenv(getter.EnvKey)
	vp := viper.New()
	vp.SetConfigName("app_" + env)
	vp.AddConfigPath("../conf/")
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

func initDB() *mysql.DB {
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

	return db
}

func initLogger() {
	logConfig := logutil.NewConfig()
	// todo 日志配置
	logutil.InitLogger(logConfig)
}
