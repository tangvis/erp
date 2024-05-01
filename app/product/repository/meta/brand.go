package meta

import (
	"context"
	"github.com/tangvis/erp/common"
)

func (r RepoImpl) GetBrandByName(ctx context.Context, userEmail string, brandName ...string) ([]BrandTab, error) {
	ret := make([]BrandTab, 0)
	err := r.db.WithContext(ctx).Model(&BrandTab{}).Where("create_by = ? and name in (?)", userEmail, brandName).Find(&ret).Error
	return ret, err
}

func (r RepoImpl) GetBrandByID(ctx context.Context, userEmail string, id ...uint64) ([]BrandTab, error) {
	ret := make([]BrandTab, 0)
	err := r.db.WithContext(ctx).Model(&BrandTab{}).Where("create_by = ? and id in (?)", userEmail, id).Find(&ret).Error
	return ret, err
}

func (r RepoImpl) SaveBrand(ctx context.Context, brand BrandTab) (BrandTab, error) {
	data, err := Save(ctx, &r, brand)
	if err != nil {
		return data, err
	}
	// refresh cache
	_, err = r.getAndCacheBrand(ctx, brand.CreateBy)

	return data, err
}

func (r RepoImpl) DeleteBrandsByIDs(ctx context.Context, userEmail string, id ...uint64) error {
	if err := r.db.WithContext(ctx).Where("create_by = ? and id in (?)", userEmail, id).Delete(&BrandTab{}).Error; err != nil {
		return err
	}
	_, err := r.getAndCacheBrand(ctx, userEmail)
	return err
}

func (r RepoImpl) GetBrandsByUser(ctx context.Context, userEmail string) ([]BrandTab, error) {
	var (
		data     = make([]BrandTab, 0)
		cacheKey = common.BrandKey(userEmail)
	)

	if err := r.cache.GetExUnmarshal(ctx, cacheKey.Key, &data, cacheKey.Expiry); err != nil {
		return nil, err
	}
	if len(data) > 0 {
		return data, nil
	}
	return r.getAndCacheBrand(ctx, userEmail)
}

func (r RepoImpl) getAndCacheBrand(ctx context.Context, userEmail string) ([]BrandTab, error) {
	var (
		data     = make([]BrandTab, 0)
		cacheKey = common.BrandKey(userEmail)
	)
	err := r.db.WithContext(ctx).Model(&BrandTab{}).Where("create_by = ?", userEmail).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, r.cache.SetExMarshal(ctx, cacheKey.Key, &data, cacheKey.Expiry)
}
