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

type PageInfo struct {
	PageNo int `json:"page_no"`
	Count  int `json:"count"`

	Offset int `json:"offset"`
}

func (p *PageInfo) Validate() error {
	if p.PageNo == 0 {
		p.PageNo = 1
	}
	if p.Count == 0 {
		p.Count = 10
	}
	p.Offset = (p.PageNo - 1) * p.Count

	return nil
}
