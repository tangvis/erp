package ecode

type Business int
type SubCode int
type System int

const (
	MainBusinessCode = 89
	SubBusinessCode  = 01
)

const (
	Ping Business = 10
	User Business = 11
)

const (
	SystemUnknown   System = 1
	SystemDB        System = 2
	SystemCache     System = 3
	SystemReadWrite System = 4
	SystemNetwork   System = 5
)

const (
	BusinessSku      Business = 1
	BusinessOrder    Business = 2
	BusinessCategory Business = 3
	BusinessBrand    Business = 4
)

func NewBusinessErrorCode(business Business, subCode SubCode) int {
	return MainBusinessCode*10e6 + SubBusinessCode*10e4 + int(business)*10e2 + int(subCode)
}

func NewSystemErrorCode(system System, subCode SubCode) int {
	return -(MainBusinessCode*10e4 + int(system)*10e2 + int(subCode))
}

const (
	CodeUnknown = -1 // PANIC之类
)

// 未知错误相关
var (
	CodeDeadlineExceed  = NewSystemErrorCode(SystemUnknown, 30)
	CodeCanceled        = NewSystemErrorCode(SystemUnknown, 31)
	CodeInvalidArgument = NewSystemErrorCode(SystemUnknown, 32)
)
