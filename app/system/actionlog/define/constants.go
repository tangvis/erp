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
