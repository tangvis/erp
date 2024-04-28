package define

import "github.com/shopspring/decimal"

type Goods struct {
	SpuInfo Spu `json:"spu_info"`
	SkuInfo Sku `json:"sku_info"`
}

type BrandBrief struct {
	ID          uint64 `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Desc        string `json:"desc,omitempty"`
	URL         string `json:"url,omitempty"`
	BrandStatus Status `json:"brand_status,omitempty"`
}

type Brand struct {
	BrandBrief
	CreateBy string `json:"create_by,omitempty"`
	Ctime    int64  `json:"ctime,omitempty"`
	Mtime    int64  `json:"mtime,omitempty"`
}

type UnitBrief struct {
	ID   uint64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Unit struct {
	UnitBrief
	CreateBy string `json:"create_by,omitempty"`
	Ctime    int64  `json:"ctime,omitempty"`
	Mtime    int64  `json:"mtime,omitempty"`
}

type CateBrief struct {
	ID   uint64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Desc string `json:"desc,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Category struct {
	CateBrief
	PID        uint64 `json:"pid,omitempty"`
	Path       string `json:"path,omitempty"`
	CateStatus Status `json:"cate_status,omitempty"`
	Ctime      int64  `json:"ctime,omitempty"`
	Mtime      int64  `json:"mtime,omitempty"`
}

type Spu struct {
	ID           uint64     `json:"id,omitempty"`
	Code         string     `json:"code,omitempty"`
	Name         string     `json:"name,omitempty"`
	EnglishName  string     `json:"english_name,omitempty"`
	UnitInfo     UnitBrief  `json:"unit_info,omitempty"`
	Desc         string     `json:"desc,omitempty"`
	URLs         string     `json:"ur_ls,omitempty"`
	CategoryInfo CateBrief  `json:"category_info,omitempty"`
	BrandInfo    BrandBrief `json:"brand_info,omitempty"`
	SpuStatus    Status     `json:"spu_status,omitempty"`
	CreateBy     string     `json:"create_by,omitempty"`
	Ctime        int64      `json:"ctime,omitempty"`
	Mtime        int64      `json:"mtime,omitempty"`
}

type Sku struct {
	ID          uint64          `json:"id,omitempty"`
	SpuID       uint64          `json:"spu_id,omitempty"`
	URLs        []string        `json:"ur_ls,omitempty"` // the first element is the main pic url
	Price       decimal.Decimal `json:"price,omitempty"`
	MarketPrice decimal.Decimal `json:"market_price,omitempty"`
	Stock       int             `json:"stock,omitempty"`
	CreateBy    string          `json:"create_by,omitempty"`
	SkuStatus   Status          `json:"sku_status,omitempty"`
	DefaultSku  bool            `json:"default_sku,omitempty"`
	Ctime       int64           `json:"ctime,omitempty"`
	Mtime       int64           `json:"mtime,omitempty"`

	Attrs []Attributes `json:"attrs"`
}

type Attributes struct {
	KeyID   int64  `json:"key_id,omitempty"`
	Key     string `json:"key,omitempty"`
	ValueID int64  `json:"value_id,omitempty"`
	Value   string `json:"value,omitempty"`
}
