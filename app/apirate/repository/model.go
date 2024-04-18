package repository

import (
	"context"
	"time"
)

type RateSettingTab struct {
	ID         uint64
	UserID     uint64
	Path       string
	QPSLimit   int
	TotalLimit int
	RateUsed   int
	ExpireTime int64
	Ctime      int64
	Mtime      int64
}

func (r *RateSettingTab) Valid() bool {
	return time.UnixMilli(r.ExpireTime).After(time.Now()) && r.RateUsed < r.TotalLimit
}

type Repo interface {
	GetRateLimitSettings(ctx context.Context) ([]RateSettingTab, error)
}

type RepoImpl struct {
}

func (r RepoImpl) GetRateLimitSettings(ctx context.Context) ([]RateSettingTab, error) {
	return nil, nil
}

func NewRepoImpl() Repo {
	return &RepoImpl{}
}
