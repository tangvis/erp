package repository

import "time"

type RateSetting struct {
	ID         uint64
	UserID     string
	Path       string
	QPSLimit   int
	TotalLimit int
	RateUsed   int
	ExpireTime int64
	Ctime      int64
	Mtime      int64
}

func (r *RateSetting) Valid() bool {
	return time.UnixMilli(r.ExpireTime).After(time.Now()) && r.RateUsed < r.TotalLimit
}
