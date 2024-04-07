package ecode

import (
	"errors"
	"fmt"
)

var _ CodeError = (*Error)(nil)

type CodeError interface {
	error
	Code() int
}

type ErrorConf struct {
	code int // 返回的JSON code

}

func NewErrorConf(code int) *ErrorConf {
	return &ErrorConf{
		code: code,
	}
}

func (conf *ErrorConf) New(message string) *Error {
	return &Error{
		conf:    conf,
		message: message,
	}
}

func (conf *ErrorConf) NewF(format string, args ...interface{}) *Error {
	return conf.New(fmt.Sprintf(format, args...))
}

func (conf *ErrorConf) Code() int {
	return conf.code
}

type Error struct {
	conf    *ErrorConf
	message string
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Code() int {
	return e.conf.code
}

func AsError(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var res *Error
	if ok := errors.As(err, &res); ok {
		return res, true
	}
	return nil, false
}

func GetErrCode(err error) int {
	if err == nil {
		return 0
	}
	if e, ok := AsError(err); ok {
		return e.Code()
	}
	return -1
}

func IsErrorCode(err error, code int) bool {
	return GetErrCode(err) == code
}
