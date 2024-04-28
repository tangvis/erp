package define

type Status int8

const (
	Init   Status = 0
	Send   Status = 1
	Failed Status = 2
)
