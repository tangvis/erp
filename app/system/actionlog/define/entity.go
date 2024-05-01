package define

import (
	"github.com/tangvis/erp/common"
)

type ListRequest struct {
	ModuleID Module `json:"module_id" binding:"required"`
	BizID    int64  `json:"biz_id" binding:"required"`

	common.PageInfo
}

type ActionLog struct {
	ID       uint64 `json:"id"`
	ModuleID Module
	BizID    uint64
	Action   string
	Operator string
	Content  string
	Ctime    int64
}
