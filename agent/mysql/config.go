package mysql

import "time"

type Config struct {
	DSN         string
	MaxIdle     int
	MaxOpen     int
	MaxIdleTime time.Duration
	MaxLifeTime time.Duration
}
