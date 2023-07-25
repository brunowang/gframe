package gferr

import (
	"fmt"
	"strconv"
	"strings"
)

type IError interface {
	error
	fmt.Stringer
	Code() int
	Msg() string
	Wrap(err error) IError
	Equal(err error) bool
}

type ecode struct {
	code int
	msg  string
	errs []error
}

func New(code int, msg string, errs ...error) IError {
	err := &ecode{
		code: code,
		msg:  msg,
		errs: errs,
	}
	return err
}

func (e ecode) Error() string {
	return e.String()
}

func (e ecode) String() string {
	var sb strings.Builder
	sb.WriteString("code: ")
	sb.WriteString(strconv.Itoa(e.code))
	sb.WriteString(", msg: ")
	sb.WriteString(e.msg)
	for _, err := range e.errs {
		sb.WriteString("; ")
		if err == nil {
			sb.WriteString("<nil>")
			continue
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

func (e ecode) Code() int {
	return e.code
}

func (e ecode) Msg() string {
	return e.msg
}

func (e ecode) Wrap(err error) IError {
	cp := e
	cp.errs = make([]error, len(e.errs))
	copy(cp.errs, e.errs)
	cp.errs = append(cp.errs, err)
	return &cp
}

func (e ecode) Equal(other error) bool {
	if other == nil {
		return false
	}
	err, ok := other.(IError)
	if !ok {
		return false
	}
	return e.Code() == err.Code()
}

type Timeout struct{ IError }

func (Timeout) Timeout() bool   { return true }
func (Timeout) Temporary() bool { return true }

func NewTimeout(code int, msg string, errs ...error) IError {
	return &Timeout{New(code, msg, errs...)}
}
