package service

import (
	"context"

	"github.com/tangvis/erp/app/product/define"
	"github.com/tangvis/erp/common"
)

type Category interface {
	Add(ctx context.Context, user *common.UserInfo, req *define.AddCateRequest) (*define.Category, error)
	List(ctx context.Context, user *common.UserInfo) ([]*define.Category, error)
	Update(ctx context.Context, user *common.UserInfo, req *define.UpdateCateRequest) (*define.Category, error)
	Remove(ctx context.Context, user *common.UserInfo, id ...uint64) error
}

type Brand interface {
	Add(ctx context.Context, user *common.UserInfo, req *define.AddBrandRequest) (*define.Brand, error)
	List(ctx context.Context, req *define.ListBrandRequest, user *common.UserInfo) (*define.ListBrandResponse, error)
	Update(ctx context.Context, user *common.UserInfo, req *define.UpdateBrandRequest) (*define.Brand, error)
	Remove(ctx context.Context, user *common.UserInfo, id ...uint64) error
}
