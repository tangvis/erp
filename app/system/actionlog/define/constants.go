package define

type Action int8
type Module int8

const (
	Add    Action = 1
	Delete Action = 2
	Update Action = 3

	Category Module = 1
	Brand    Module = 2
)
