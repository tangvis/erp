package impl

import (
	"context"
	actionLogDefine "github.com/tangvis/erp/app/system/actionlog/define"
	actionLog "github.com/tangvis/erp/app/system/actionlog/service"

	"github.com/tangvis/erp/app/product/converter"
	"github.com/tangvis/erp/app/product/define"
	"github.com/tangvis/erp/app/product/repository/meta"
	"github.com/tangvis/erp/app/product/service"
	"github.com/tangvis/erp/common"
)

type BrandImpl struct {
	repo      meta.Repo
	actionLog actionLog.APP
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
	b.actionLog.AsyncCreate(ctx, user.Email, actionLogDefine.Brand, brand.ID, actionLogDefine.ADD, nil, nil)
	return converter.BrandConvert(brand), nil
}

func (b BrandImpl) List(ctx context.Context, req *define.ListBrandRequest, user *common.UserInfo) ([]*define.Brand, error) {
	brands, err := b.repo.GetBrandsByUser(ctx, user.Email, meta.BrandQuery{
		Name:   req.Name,
		Offset: req.Offset,
		Limit:  req.Count,
	})
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
	b.actionLog.AsyncCreate(ctx, user.Email, actionLogDefine.Brand, brand.ID, actionLogDefine.UPDATE, brands[0], ret)
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
	err := b.repo.DeleteBrandsByIDs(ctx, user.Email, id...)
	if err != nil {
		return err
	}
	for _, _id := range id {
		b.actionLog.AsyncCreate(ctx, user.Email, actionLogDefine.Brand, _id, actionLogDefine.DELETE, nil, nil)
	}
	return nil
}

func NewBrandImpl(
	repo meta.Repo,
	actionLog actionLog.APP,
) service.Brand {
	return &BrandImpl{
		repo:      repo,
		actionLog: actionLog,
	}
}
