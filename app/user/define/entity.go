package define

import "fmt"

type UserEntity struct {
	ID          uint64
	Username    string
	Passwd      string
	PhoneNumber string
	Email       string
	Status      UserStatus

	Ctime int64
	Mtime int64
}

type UserQuery struct {
	Usernames    []string
	PhoneNumbers []string
	Emails       []string
}

func (u *UserQuery) Valid() error {
	if len(u.Usernames) == 0 && len(u.PhoneNumbers) == 0 && len(u.Emails) == 0 {
		return fmt.Errorf("params invalid")
	}

	return nil
}

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

func (u *LoginRequest) Validate() error {
	if u.Username == "" || u.Email == "" {
		return fmt.Errorf("username or email is needed")
	}
	return nil
}
