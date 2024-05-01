package repository

import (
	"context"
	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/app/system/actionlog/define"
)

type ListQuery struct {
	ModuleID define.Module
	BizID    int64
	Offset   int
	Limit    int
}

type Repo interface {
	Save(ctx context.Context, tab ActionLogTab) error
	List(ctx context.Context, query ListQuery) ([]ActionLogTab, error)
}

type RepoImpl struct {
	db *mysql.DB
}

func (r RepoImpl) List(ctx context.Context, query ListQuery) ([]ActionLogTab, error) {
	data := make([]ActionLogTab, 0)
	err := r.db.WithContext(ctx).Model(&ActionLogTab{}).
		Where("module_id = ? and biz_id = ?", query.ModuleID, query.BizID).Order("id").
		Offset(query.Offset).
		Limit(query.Limit).
		Find(&data).
		Error

	return data, err
}

func (r RepoImpl) Save(ctx context.Context, tab ActionLogTab) error {
	return r.db.WithContext(ctx).Model(&ActionLogTab{}).Save(&tab).Error
}

func NewRepoImpl(db *mysql.DB) Repo {
	return &RepoImpl{db: db}
}
