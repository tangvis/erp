package common

import "encoding/json"

type UserInfo struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	LoginTime   int64  `json:"login_time"`
	IP          string `json:"ip"`
	SessionID   string `json:"session_id,omitempty"`
}

func (u *UserInfo) String() string {
	b, _ := json.Marshal(u)

	return string(b)
}
