package config

import "os"

type Env string

const (
	EnvKey = "env"

	DEV  Env = "dev"
	TEST Env = "test"
	LIVE Env = "live"
)

func IsLive() bool {
	return Env(os.Getenv(EnvKey)) == LIVE
}
