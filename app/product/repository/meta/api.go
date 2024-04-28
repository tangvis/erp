package meta

import "context"

type Repo interface {
	CreateSpu(ctx context.Context, spu SpuTab) (SpuTab, error)
	CreateBrand(ctx context.Context, brand BrandTab) (BrandTab, error)
	CreateCategory(ctx context.Context, category CategoryTab) (CategoryTab, error)
	CreateSku(ctx context.Context, sku SkuTab) (SkuTab, error)
	CreateUnit(ctx context.Context, unit UnitTab) (UnitTab, error)
	CreateURL(ctx context.Context, url URLTab) (URLTab, error)
}
