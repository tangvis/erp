package repository

import (
	"context"
	"time"

	"github.com/tangvis/erp/agent/mysql"
)

type RateSettingTab struct {
	mysql.BaseModel
	UserID     uint64
	Path       string
	QPSLimit   int
	TotalLimit int
	RateUsed   int
	ExpireTime int64
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
