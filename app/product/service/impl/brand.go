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
	if err := b.CheckBrandName(ctx, user.Email, req.Name, 0); err != nil {
		return nil, err
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
	brands, err := b.repo.GetBrandsByUser(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	return converter.BrandsConvert(brands), nil
}

func (b BrandImpl) Update(ctx context.Context, user *common.UserInfo, req *define.UpdateBrandRequest) (*define.Brand, error) {
	if err := b.CheckBrandName(ctx, user.Email, req.Name, req.ID); err != nil {
		return nil, err
	}
	brands, err := b.repo.GetBrandByID(ctx, user.Email, req.ID)
	if err != nil {
		return nil, err
	}
	if len(brands) == 0 {
		return nil, common.ErrBrandNotExists
	}
	brand := brands[0]
	brand.Name = req.Name
	brand.Desc = req.Desc
	brand.URL = req.URL
	ret, err := b.repo.SaveBrand(ctx, brand)
	if err != nil {
		return nil, err
	}
	return converter.BrandConvert(ret), nil
}

func (b BrandImpl) CheckBrandName(ctx context.Context, userEmail, name string, id uint64) error {
	brands, err := b.repo.GetBrandByName(ctx, userEmail, name)
	if err != nil {
		return err
	}
	if len(brands) != 0 && (id == 0 || brands[0].ID != id) {
		return common.ErrBrandNameConflict
	}
	return nil
}

func (b BrandImpl) Remove(ctx context.Context, user *common.UserInfo, id ...uint64) error {
	return b.repo.DeleteBrandsByIDs(ctx, user.Email, id...)
}

func NewBrandImpl(
	repo meta.Repo,
) service.Brand {
	return &BrandImpl{
		repo: repo,
	}
}
