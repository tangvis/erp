package meta

import (
	"context"
	"github.com/tangvis/erp/common"
)

func (r RepoImpl) SaveCategory(ctx context.Context, category CategoryTab) (CategoryTab, error) {
	data, err := Save(ctx, &r, category)
	if err != nil {
		return data, err
	}
	// refresh cache
	_, err = r.getAndCacheCategory(ctx, category.CreateBy)

	return data, err
}

func (r RepoImpl) GetCategoryByID(ctx context.Context, userEmail string, id ...uint64) ([]CategoryTab, error) {
	ret := make([]CategoryTab, 0)
	err := r.db.WithContext(ctx).Model(&CategoryTab{}).Where("create_by = ? and id in (?)", userEmail, id).Find(&ret).Error
	return ret, err
}

func (r RepoImpl) GetCategoryByPID(ctx context.Context, userEmail string, pid ...uint64) ([]CategoryTab, error) {
	ret := make([]CategoryTab, 0)
	err := r.db.WithContext(ctx).Model(&CategoryTab{}).Where("create_by = ? and pid in (?)", userEmail, pid).Find(&ret).Error
	return ret, err
}

func (r RepoImpl) DeleteCategoryByIDs(ctx context.Context, userEmail string, id ...uint64) error {
	if err := r.db.WithContext(ctx).Where("create_by = ? and id in (?)", userEmail, id).Delete(&CategoryTab{}).Error; err != nil {
		return err
	}
	_, err := r.getAndCacheCategory(ctx, userEmail)
	return err
}

func (r RepoImpl) GetCategoryByName(ctx context.Context, userEmail string, name ...string) ([]CategoryTab, error) {
	ret := make([]CategoryTab, 0)
	err := r.db.WithContext(ctx).Model(&CategoryTab{}).Where("create_by = ? and name in (?)", userEmail, name).Find(&ret).Error
	return ret, err
}

func (r RepoImpl) GetCategoriesByUser(ctx context.Context, userEmail string) ([]CategoryTab, error) {
	var (
		data     = make([]CategoryTab, 0)
		cacheKey = common.CategoryKey(userEmail)
	)

	if err := r.cache.GetExUnmarshal(ctx, cacheKey.Key, &data, cacheKey.Expiry); err != nil {
		return nil, err
	}
	if len(data) > 0 {
		return data, nil
	}
	return r.getAndCacheCategory(ctx, userEmail)
}

func (r RepoImpl) getAndCacheCategory(ctx context.Context, userEmail string) ([]CategoryTab, error) {
	var (
		data     = make([]CategoryTab, 0)
		cacheKey = common.CategoryKey(userEmail)
	)
	err := r.db.WithContext(ctx).Model(&CategoryTab{}).Where("create_by = ?", userEmail).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, r.cache.SetExMarshal(ctx, cacheKey.Key, &data, cacheKey.Expiry)
}
