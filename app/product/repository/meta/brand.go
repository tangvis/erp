package meta

import (
	"context"
	"gorm.io/gorm"
)

type BrandQuery struct {
	Name   string
	Offset int
	Limit  int
}

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
	_, err = r.getAndCacheBrand(ctx, brand.CreateBy, BrandQuery{})

	return data, err
}

func (r RepoImpl) DeleteBrandsByIDs(ctx context.Context, userEmail string, id ...uint64) error {
	if err := r.db.WithContext(ctx).Where("create_by = ? and id in (?)", userEmail, id).Delete(&BrandTab{}).Error; err != nil {
		return err
	}
	_, err := r.getAndCacheBrand(ctx, userEmail, BrandQuery{})
	return err
}

func (r RepoImpl) GetBrandsByUser(ctx context.Context, userEmail string, query BrandQuery) ([]BrandTab, error) {
	//qJson, _ := json.Marshal(query)
	//var (
	//	data     = make([]BrandTab, 0)
	//	cacheKey = common.BrandKey(userEmail, string(qJson))
	//)
	//
	//if err := r.cache.GetExUnmarshal(ctx, cacheKey.Key, &data, cacheKey.Expiry); err != nil {
	//	return nil, err
	//}
	//if len(data) > 0 {
	//	return data, nil
	//}
	return r.getAndCacheBrand(ctx, userEmail, query)
}

func (r RepoImpl) brandListGetQuery(ctx context.Context, userEmail string, query BrandQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&BrandTab{}).Where("create_by = ?", userEmail)
	if len(query.Name) > 0 {
		db = db.Where("name like ?", "%"+query.Name+"%")
	}

	return db
}

func (r RepoImpl) getAndCacheBrand(ctx context.Context, userEmail string, query BrandQuery) ([]BrandTab, error) {
	//qJson, _ := json.Marshal(query)
	var (
		data = make([]BrandTab, 0)
		//cacheKey = common.BrandKey(userEmail, string(qJson))
	)
	db := r.brandListGetQuery(ctx, userEmail, query).Offset(query.Offset)
	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	//return data, r.cache.SetExMarshal(ctx, cacheKey.Key, &data, cacheKey.Expiry)
	return data, nil
}

func (r RepoImpl) CountBrand(ctx context.Context, userEmail string, query BrandQuery) (int64, error) {
	var count int64
	err := r.brandListGetQuery(ctx, userEmail, query).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
