package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/tangvis/erp/agent/mysql"
	"strings"
	"testing"
	"time"
)

func TestGetMySQLConfig(t *testing.T) {
	// Initialize Viper with test configuration in TOML format
	viper.SetConfigType("toml")
	var tomlExample = strings.NewReader(`
[mysql]
DSN = "user:password@/dbname"
MaxIdle = 10
MaxOpen = 100
MaxIdleTime = "10s"
MaxLifeTime = "1h"
`)
	err := viper.ReadConfig(tomlExample)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Create a new configGetter instance
	cfgGetter := NewConfigGetter(viper.GetViper())

	// Call the GetMySQLConfig method
	mysqlCfg, err := cfgGetter.GetMySQLConfig()

	// Assert no error returned
	assert.NoError(t, err)

	// Assert that the mysql.Config structure is correctly populated
	expectedMaxIdleTime, _ := time.ParseDuration("10s")
	expectedMaxLifeTime, _ := time.ParseDuration("1h")
	expectedCfg := mysql.Config{
		DSN:         "user:password@/dbname",
		MaxIdle:     10,
		MaxOpen:     100,
		MaxIdleTime: expectedMaxIdleTime,
		MaxLifeTime: expectedMaxLifeTime,
	}

	assert.Equal(t, expectedCfg, mysqlCfg)
}
