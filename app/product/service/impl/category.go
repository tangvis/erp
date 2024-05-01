package impl

import (
	"context"
	"github.com/tangvis/erp/app/product/converter"
	"github.com/tangvis/erp/app/product/define"
	"github.com/tangvis/erp/app/product/repository/meta"
	"github.com/tangvis/erp/app/product/service"
	actionLogDefine "github.com/tangvis/erp/app/system/actionlog/define"
	actionLog "github.com/tangvis/erp/app/system/actionlog/service"
	"github.com/tangvis/erp/common"
)

type CategoryImpl struct {
	repo      meta.Repo
	actionLog actionLog.APP
}

func (c CategoryImpl) List(ctx context.Context, user *common.UserInfo) ([]*define.Category, error) {
	categories, err := c.repo.GetCategoriesByUser(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	return converter.CategoriesConvert(categories), nil
}

func (c CategoryImpl) Add(ctx context.Context, user *common.UserInfo, req *define.AddCateRequest) (*define.Category, error) {
	if err := c.CheckBeforeAdd(ctx, user, req); err != nil {
		return nil, err
	}
	category, err := c.repo.SaveCategory(ctx, meta.CategoryTab{
		PID:      req.PID,
		Name:     req.Name,
		Desc:     req.Desc,
		URL:      req.URL,
		CreateBy: user.Email,
	})
	if err != nil {
		return nil, err
	}
	c.actionLog.AsyncCreate(ctx, user.Email, actionLogDefine.Category, category.ID, actionLogDefine.Add, nil, nil)
	return converter.CategoryConvert(category), nil
}

func (c CategoryImpl) CheckBeforeAdd(ctx context.Context, user *common.UserInfo, req *define.AddCateRequest) error {
	if req.PID != 0 {
		parent, err := c.repo.GetCategoryByID(ctx, user.Email, req.PID)
		if err != nil {
			return err
		}
		if len(parent) == 0 {
			return common.ErrCategoryParentNotExists
		}
	}
	return c.CheckCategoryName(ctx, user.Email, req.Name, 0)
}

func (c CategoryImpl) Remove(ctx context.Context, user *common.UserInfo, id ...uint64) error {
	if err := c.CheckBeforeRemove(ctx, user, id...); err != nil {
		return err
	}
	return c.repo.DeleteCategoryByIDs(ctx, user.Email, id...)
}

func (c CategoryImpl) CheckBeforeRemove(ctx context.Context, user *common.UserInfo, id ...uint64) error {
	categories, err := c.repo.GetCategoryByPID(ctx, user.Email, id...)
	if err != nil {
		return err
	}
	if len(categories) > 0 {
		return common.ErrCategoryHasChildren
	}
	return nil
}

func (c CategoryImpl) Update(ctx context.Context, user *common.UserInfo, req *define.UpdateCateRequest) (*define.Category, error) {
	cate, err := c.CheckBeforeUpdate(ctx, user, req)
	if err != nil {
		return nil, err
	}
	cate.Name = req.Name
	cate.URL = req.URL
	cate.Desc = req.Desc
	category, err := c.repo.SaveCategory(ctx, cate)
	if err != nil {
		return nil, err
	}
	return converter.CategoryConvert(category), nil
}

func (c CategoryImpl) CheckBeforeUpdate(ctx context.Context, user *common.UserInfo, req *define.UpdateCateRequest) (meta.CategoryTab, error) {
	categories, err := c.GetCategoryMap(ctx, user, req.PID, req.ID)
	if err != nil {
		return meta.CategoryTab{}, err
	}
	if req.PID > 0 {
		_, ok := categories[req.PID]
		if !ok {
			return meta.CategoryTab{}, common.ErrCategoryParentNotExists
		}
	}
	cate, ok := categories[req.ID]
	if !ok {
		return meta.CategoryTab{}, common.ErrCategoryNotExists
	}
	if err = c.CheckCategoryName(ctx, user.Email, req.Name, req.ID); err != nil {
		return meta.CategoryTab{}, err
	}

	return cate, nil
}

func (c CategoryImpl) CheckCategoryName(ctx context.Context, userEmail, name string, id uint64) error {
	data, err := c.repo.GetCategoryByName(ctx, userEmail, name)
	if err != nil {
		return err
	}
	if len(data) > 0 && (id == 0 || data[0].ID != id) {
		return common.ErrCategoryNameConflict
	}
	return nil
}

func (c CategoryImpl) GetCategoryMap(ctx context.Context, user *common.UserInfo, id ...uint64) (map[uint64]meta.CategoryTab, error) {
	categories, err := c.repo.GetCategoryByID(ctx, user.Email, id...)
	if err != nil {
		return nil, err
	}
	m := make(map[uint64]meta.CategoryTab)
	for _, category := range categories {
		m[category.ID] = category
	}

	return m, nil
}

func NewCategoryImpl(
	repo meta.Repo,
	actionLog actionLog.APP,
) service.Category {
	return &CategoryImpl{
		repo:      repo,
		actionLog: actionLog,
	}
}
