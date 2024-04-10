package repository

import "github.com/tangvis/erp/app/user/define"

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
