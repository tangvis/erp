package repository

import (
	"context"
	"github.com/tangvis/erp/agent/mysql"
)

type Repo interface {
	Save(ctx context.Context, tab ActionLogTab) error
}

type RepoImpl struct {
	db *mysql.DB
}

func (r RepoImpl) Save(ctx context.Context, tab ActionLogTab) error {
	//TODO implement me
	panic("implement me")
}

func NewRepoImpl(db *mysql.DB) Repo {
	return &RepoImpl{db: db}
}
