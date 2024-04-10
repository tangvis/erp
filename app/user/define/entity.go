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
