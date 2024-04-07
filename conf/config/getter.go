package config

import (
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/tangvis/erp/agent/mysql"
)

var Config Getter

func InitConfig() {
	vp := initViper()
	Config = NewConfigGetter(vp)
}

func initViper() *viper.Viper {
	env := os.Getenv(EnvKey)
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

type Getter interface {
	GetMySQLConfig() (mysql.Config, error)
	GetEnableResponseTraceID() bool
}

type configGetter struct {
	viper *viper.Viper
}

func NewConfigGetter(vp *viper.Viper) Getter {
	return &configGetter{
		viper: vp,
	}
}

func (c configGetter) GetMySQLConfig() (mysql.Config, error) {
	var tempCfg struct {
		DSN         string `toml:"dsn"`
		MaxIdle     int    `toml:"max_idle"`
		MaxOpen     int    `toml:"max_open"`
		MaxIdleTime string `toml:"max_idle_time"`
		MaxLifeTime string `toml:"max_life_time"`
	}
	if err := c.viper.UnmarshalKey("mysql", &tempCfg); err != nil {
		return mysql.Config{}, err
	}
	expectedMaxIdleTime, err := time.ParseDuration(tempCfg.MaxIdleTime)
	if err != nil {
		return mysql.Config{}, err
	}
	expectedMaxLifeTime, err := time.ParseDuration(tempCfg.MaxLifeTime)
	if err != nil {
		return mysql.Config{}, err
	}

	return mysql.Config{
		DSN:         tempCfg.DSN,
		MaxIdle:     tempCfg.MaxIdle,
		MaxOpen:     tempCfg.MaxOpen,
		MaxIdleTime: expectedMaxIdleTime,
		MaxLifeTime: expectedMaxLifeTime,
	}, err
}

func (c configGetter) GetEnableResponseTraceID() bool {
	return c.viper.GetBool("middleware.response_trace_id")
}
