package define

type UserTab struct {
	ID          uint64
	Username    string
	Passwd      string
	PhoneNumber string
	Email       string
	Status      UserStatus

	Ctime int64
	Mtime int64
}
