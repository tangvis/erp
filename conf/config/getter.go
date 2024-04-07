package config

import (
	"github.com/spf13/viper"
	"github.com/tangvis/erp/agent/mysql"
	"time"
)

type Getter interface {
	GetMySQLConfig() (mysql.Config, error)
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
