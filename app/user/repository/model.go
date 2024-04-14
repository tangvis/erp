package repository

import (
	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/app/user/define"
)

func (q *UserTab) TableName() string {
	return "user_tab"
}

type UserTab struct {
	mysql.BaseModel

	Username    string
	Passwd      string
	PhoneNumber string
	Email       string
	UserStatus  define.UserStatus
}
