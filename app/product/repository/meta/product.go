package meta

import (
	"context"

	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/agent/redis"
)

type RepoImpl struct {
	db    *mysql.DB
	cache redis.Cache
}

func Save[T any](ctx context.Context, r *RepoImpl, entity T) (T, error) {
	if err := r.db.WithContext(ctx).Save(&entity).Error; err != nil {
		var zero T
		return zero, err // Return zero value of T if error
	}
	return entity, nil
}

func (r RepoImpl) CreateSpu(ctx context.Context, spu SpuTab) (SpuTab, error) {
	return Save(ctx, &r, spu)
}

// CreateBrand inserts a new brand into the database
func (r RepoImpl) CreateBrand(ctx context.Context, brand BrandTab) (BrandTab, error) {
	return Save(ctx, &r, brand)
}

// CreateSku inserts a new Sku into the database
func (r RepoImpl) CreateSku(ctx context.Context, sku SkuTab) (SkuTab, error) {
	return Save(ctx, &r, sku)
}

// CreateUnit inserts a new Unit into the database
func (r RepoImpl) CreateUnit(ctx context.Context, unit UnitTab) (UnitTab, error) {
	return Save(ctx, &r, unit)
}

// CreateSkuAttr inserts a new SkuAttr into the database
func (r RepoImpl) CreateSkuAttr(ctx context.Context, skuAttr SkuAttrTab) (SkuAttrTab, error) {
	return Save(ctx, &r, skuAttr)
}

// CreateAttributeKey inserts a new AttributeKey into the database
func (r RepoImpl) CreateAttributeKey(ctx context.Context, attributeKey AttributeKeyTab) (AttributeKeyTab, error) {
	return Save(ctx, &r, attributeKey)
}

// CreateAttributeValue inserts a new AttributeValue into the database
func (r RepoImpl) CreateAttributeValue(ctx context.Context, attributeValue AttributeValueTab) (AttributeValueTab, error) {
	return Save(ctx, &r, attributeValue)
}

// CreateURL inserts a new URL into the database
func (r RepoImpl) CreateURL(ctx context.Context, url URLTab) (URLTab, error) {
	return Save(ctx, &r, url)
}

func NewRepoImpl(
	db *mysql.DB,
	cache redis.Cache,
) Repo {
	return &RepoImpl{
		db:    db,
		cache: cache,
	}
}
