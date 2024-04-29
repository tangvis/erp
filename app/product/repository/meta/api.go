package meta

import "context"

type Repo interface {
	CreateSpu(ctx context.Context, spu SpuTab) (SpuTab, error)
	CreateBrand(ctx context.Context, brand BrandTab) (BrandTab, error)
	SaveCategory(ctx context.Context, category CategoryTab) (CategoryTab, error)
	CreateSku(ctx context.Context, sku SkuTab) (SkuTab, error)
	CreateUnit(ctx context.Context, unit UnitTab) (UnitTab, error)
	CreateURL(ctx context.Context, url URLTab) (URLTab, error)

	GetCategoryByID(ctx context.Context, id ...uint64) ([]CategoryTab, error)
	GetCategoryByName(ctx context.Context, name ...string) ([]CategoryTab, error)
	GetCategoriesByUser(ctx context.Context, userEmail string) ([]CategoryTab, error)
}
