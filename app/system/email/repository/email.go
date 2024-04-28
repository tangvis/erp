package repository

import (
	"context"

	"github.com/tangvis/erp/agent/mysql"
)

type Repo interface {
	Save(ctx context.Context, record EmailRecordTab) error
}

type EmailRepo struct {
	db *mysql.DB
}

func (r EmailRepo) Save(ctx context.Context, record EmailRecordTab) error {
	return r.db.WithContext(ctx).Save(&record).Error
}
