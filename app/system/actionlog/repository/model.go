package repository

import (
	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/app/system/actionlog/define"
)

type ActionLogTab struct {
	ID         uint64
	ModuleID   define.Module
	BizID      uint64
	ActionType define.Action
	Operator   string
	// json type
	Content string
	mysql.BaseModel
}

func (tab *ActionLogTab) TableName() string {
	return "action_log_tab"
}
