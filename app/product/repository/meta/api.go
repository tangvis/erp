package meta

import "context"

type Repo interface {
	CreateSpu(ctx context.Context, spu SpuTab) (SpuTab, error)
	SaveCategory(ctx context.Context, category CategoryTab) (CategoryTab, error)
	CreateSku(ctx context.Context, sku SkuTab) (SkuTab, error)
	CreateUnit(ctx context.Context, unit UnitTab) (UnitTab, error)
	CreateURL(ctx context.Context, url URLTab) (URLTab, error)

	GetCategoryByID(ctx context.Context, userEmail string, id ...uint64) ([]CategoryTab, error)
	GetCategoryByPID(ctx context.Context, userEmail string, pid ...uint64) ([]CategoryTab, error)
	GetCategoryByName(ctx context.Context, userEmail string, name ...string) ([]CategoryTab, error)
	GetCategoriesByUser(ctx context.Context, userEmail string) ([]CategoryTab, error)
	DeleteCategoryByIDs(ctx context.Context, userEmail string, id ...uint64) error

	GetBrandByID(ctx context.Context, userEmail string, id ...uint64) ([]BrandTab, error)
	GetBrandByName(ctx context.Context, userEmail string, brandName ...string) ([]BrandTab, error)
	GetBrandsByUser(ctx context.Context, userEmail string, query BrandQuery) ([]BrandTab, error)
	SaveBrand(ctx context.Context, brand BrandTab) (BrandTab, error)
	DeleteBrandsByIDs(ctx context.Context, userEmail string, id ...uint64) error
	CountBrand(ctx context.Context, userEmail string, query BrandQuery) (int64, error)
}
