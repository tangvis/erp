package impl

import (
	"context"
	"github.com/tangvis/erp/app/product/converter"
	"github.com/tangvis/erp/app/product/define"
	"github.com/tangvis/erp/app/product/repository/meta"
	"github.com/tangvis/erp/app/product/service"
	"github.com/tangvis/erp/common"
)

type BrandImpl struct {
	repo meta.Repo
}

func (b BrandImpl) Add(ctx context.Context, user *common.UserInfo, req *define.AddBrandRequest) (*define.Brand, error) {
	// check if name exists
	brands, err := b.repo.GetBrandByName(ctx, user.Email, req.Name)
	if err != nil {
		return nil, err
	}
	if len(brands) > 0 {
		return nil, common.ErrBrandNameConflict
	}
	brand, err := b.repo.SaveBrand(ctx, meta.BrandTab{
		Name:     req.Name,
		Desc:     req.Desc,
		URL:      req.URL,
		CreateBy: user.Email,
	})
	if err != nil {
		return nil, err
	}
	return converter.BrandConvert(brand), nil
}

func (b BrandImpl) List(ctx context.Context, user *common.UserInfo) ([]*define.Brand, error) {
	//TODO implement me
	panic("implement me")
}

func (b BrandImpl) Update(ctx context.Context, user *common.UserInfo, req *define.UpdateBrandRequest) (*define.Brand, error) {
	//TODO implement me
	panic("implement me")
}

func (b BrandImpl) Remove(ctx context.Context, user *common.UserInfo, id ...uint64) error {
	//TODO implement me
	panic("implement me")
}

func NewBrandImpl(
	repo meta.Repo,
) service.Brand {
	return &BrandImpl{
		repo: repo,
	}
}
