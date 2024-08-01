package xlerror

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	codes = map[int]struct{}{}
)

var (
	NullError       = add(0, "")
	Success         = add(200, "SUCCESS")
	ErrRequest      = add(400, "bad request")
	ErrNotFind      = add(404, "not found")
	ErrForbidden    = add(403, "forbidden")
	ErrNoPermission = add(405, "no permission")
	ErrServer       = add(500, "server error")
	ErrTimeout      = add(504, "request timeout")
	ErrRateLimit    = add(600, "wait for retry")
	ErrForWait      = add(900, "wait for retry")
	ErrToken        = add(1000, "Token Error")
)

// New 创建一个错误
func New(code int, msg string, link ...string) Error {
	if code < 1000 && code > 0 {
		panic("error code must be greater than 1000")
	}
	return add(code, msg, link...)
}

// add only inner error
func add(code int, msg string, link ...string) Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", code))
	}
	codes[code] = struct{}{}
	var _link string
	if len(link) > 0 {
		_link = link[0]
	}
	return Error{
		code: code, message: msg, link: _link,
	}
}

type Errors interface {
	// Error sometimes Error return Code in string form
	Error() string
	// Code get error code.
	Code() int
	// Link get error code.
	Link() string
	// Message get code message.
	Message() string
	// Details get error detail,it may be nil.
	Details() []interface{}
	// Equal for compatible.
	Equal(error) bool
	// Reload Message
	Reload(string) Error
}

type Error struct {
	code    int
	message string
	link    string
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Reload(message string) Error {
	e.message = message
	return e
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Link() string {
	return e.link
}

func (e Error) Details() []interface{} { return nil }

func (e Error) Equal(err error) bool { return Equal(err, e) }

func String(e string) Error {
	if e == "" {
		return NullError
	}
	return Error{
		code:    500,
		message: "",
	}
}

// Cause 解析错误码
func Cause(err error) Errors {
	if err == nil {
		return NullError
	}
	if ec, ok := errors.Cause(err).(Errors); ok {
		return ec
	}
	return ErrServer
}

// Equal 两个错误错误码是否一致
func Equal(err error, e Error) bool {
	return Cause(err).Code() == e.Code()
}

// Wrap 对错误描述进行包装
func Wrap(err error, message string) Errors {
	ec := Cause(err)
	return Error{
		code:    ec.Code(),
		message: fmt.Sprintf("%s: %s", ec.Message(), message),
	}
}
