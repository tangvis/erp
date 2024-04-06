package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type DB struct {
	*gorm.DB
}

func NewMySQL(config Config) (*DB, error) {
	return newMySQL(config.DSN, config.MaxIdle, config.MaxOpen, config.MaxLifeTime, config.MaxIdleTime)
}

func newMySQL(dsn string, maxIdle, maxOpen int, maxLifetime, maxIdleTime time.Duration) (*DB, error) {
	gormConf := &gorm.Config{
		// todo logger
	}
	db, err := gorm.Open(mysql.Open(dsn), gormConf)
	if err != nil {
		return nil, err
	}
	rawDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	rawDB.SetMaxIdleConns(maxIdle)
	rawDB.SetMaxOpenConns(maxOpen)
	rawDB.SetConnMaxLifetime(maxLifetime)
	rawDB.SetConnMaxIdleTime(maxIdleTime)
	return &DB{db}, nil
}
