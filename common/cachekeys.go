package common

import "time"

const (
	Category = "product:"
	Brand    = "brand:"
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

func BrandKey(email string) CacheKey {
	return CacheKey{
		Key:    Brand + email,
		Expiry: 30 * time.Minute,
	}
}
