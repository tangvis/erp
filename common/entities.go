package common

import json2 "encoding/json"

type UserInfo struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	LoginTime   int64  `json:"login_time"`
	IP          string `json:"ip"`
}

func (u *UserInfo) String() string {
	b, _ := json2.Marshal(u)

	return string(b)
}
