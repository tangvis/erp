package meta

import "github.com/tangvis/erp/agent/mysql"

type Category struct {
	mysql.BaseModel
}

type Brand struct {
	mysql.BaseModel
}

type Spu struct {
	mysql.BaseModel
}

type Sku struct {
	mysql.BaseModel
}

type AttributeKey struct {
	mysql.BaseModel
}

type AttributeValue struct {
	mysql.BaseModel
}
