package repository

import (
	"github.com/tangvis/erp/app/user/define"
)

func (q *UserTab) TableName() string {
	return "user_tab"
}

type UserTab struct {
	ID          uint64
	Username    string
	Passwd      string
	PhoneNumber string
	Email       string
	Status      define.UserStatus

	Ctime int64
	Mtime int64
}
