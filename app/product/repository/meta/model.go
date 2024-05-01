package meta

import (
	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/app/product/define"
	"github.com/tangvis/erp/common"
)

type BrandTab struct {
	Name        string        `al:"品牌名"`
	Desc        string        `al:"品牌描述"`
	URL         string        // URLTab `al:"图片链接"`
	BrandStatus define.Status `al:"状态"`
	CreateBy    string        // email of creator

	mysql.BaseModel
}

type CategoryTab struct {
	PID        uint64 `gorm:"column:pid"`
	Name       string
	Desc       string
	URL        string // URLTab
	Path       string // {pid} -> {child_id} -> ...
	CateStatus define.Status
	CreateBy   string // email of creator

	mysql.BaseModel
}

func (tab *CategoryTab) TableName() string {
	return "category_tab"
}

type SpuTab struct {
	Code        string
	Name        string
	EnglishName string
	UnitID      uint64 // UnitTab
	Desc        string
	URLs        string // URLTab split by ","
	CategoryID  uint64 // CategoryTab
	BrandID     uint64 // BrandTab
	SpuStatus   define.Status
	CreateBy    string

	mysql.BaseModel
}

type SkuTab struct {
	SpuID       uint64 // SpuTab
	URLs        string // URLTab split by ","
	Price       int    // real price * 100
	MarketPrice int    // real market price * 100
	Stock       int    // 库存
	CreateBy    string // email of creator
	SkuStatus   define.Status
	DefaultSku  common.Boolean // if is the default sku of spu

	mysql.BaseModel
}

type UnitTab struct {
	Name     string
	CreateBy string

	mysql.BaseModel
}

// SkuAttrTab sku和属性映射表，(sku_id,attr_key_id)应该是唯一键，
// 也就是说一个sku某一个属性只能对应一个值，不同的值组合成不同的sku
type SkuAttrTab struct {
	SkuID       uint64
	AttrKeyID   uint64
	AttrValueID uint64

	mysql.BaseModel
}

type AttributeKeyTab struct {
	Name     string
	CreateBy string

	mysql.BaseModel
}

type AttributeValueTab struct {
	KeyID    uint64
	Value    string
	CreateBy string

	mysql.BaseModel
}

type URLTab struct {
	URL     string
	URLType define.URLType

	mysql.BaseModel
}
