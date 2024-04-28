package common

type Boolean uint8

func (b Boolean) True() bool {
	return b == T
}

const (
	UserInfoKey = "user_info"

	F Boolean = 0
	T Boolean = 1
)
