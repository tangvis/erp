package common

import "time"

const (
	Category = "category:"
)

type CacheKey struct {
	Key    string
	Expiry time.Duration
}

func CategoryKey(email string) CacheKey {
	return CacheKey{
		Key:    Category + email,
		Expiry: 30 * time.Minute,
	}
}
