package config

import (
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/agent/redis"
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
	GetCacheConfig() (redis.Config, error)
	GetMiddleWareConfig() (MiddlewareConfig, error)
	GetEnableResponseTraceID() bool
	GetEnableLogRequest() bool
}

type configGetter struct {
	viper *viper.Viper
}

func (c configGetter) GetCacheConfig() (redis.Config, error) {
	var tempCfg struct {
		Addr   string
		Passwd string
		DB     int
	}
	if err := c.viper.UnmarshalKey("cache", &tempCfg); err != nil {
		return redis.Config{}, err
	}

	return redis.Config{
		Addr:   tempCfg.Addr,
		Passwd: tempCfg.Passwd,
		DB:     tempCfg.DB,
	}, nil
}

func NewConfigGetter(vp *viper.Viper) Getter {
	return &configGetter{
		viper: vp,
	}
}

func (c configGetter) GetMySQLConfig() (mysql.Config, error) {
	var tempCfg struct {
		DSN         string
		MaxIdle     int
		MaxOpen     int
		MaxIdleTime string
		MaxLifeTime string
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

func (c configGetter) GetMiddleWareConfig() (MiddlewareConfig, error) {
	var tempCfg MiddlewareConfig
	if err := c.viper.UnmarshalKey("middleware", &tempCfg); err != nil {
		return MiddlewareConfig{}, err
	}
	return tempCfg, nil
}

func (c configGetter) GetEnableLogRequest() bool {
	return c.viper.GetBool("middleware.log_request")
}
