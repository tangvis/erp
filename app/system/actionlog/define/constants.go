package define

type Action int8
type Module int8

const (
	ADD    Action = 1
	DELETE Action = 2
	UPDATE Action = 3

	Category Module = 1
	Brand    Module = 2
)

func (a Action) String() string {
	switch a {
	case ADD:
		return "创建"
	case DELETE:
		return "删除"
	case UPDATE:
		return "修改"
	default:
		return ""
	}
}
