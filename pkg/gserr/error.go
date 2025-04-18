package gserr

import (
	"fmt"
)

type Error struct {
	code uint32
	msg  string
}

func New(code uint32, msg string) *Error {
	return &Error{
		code: code,
		msg:  msg,
	}
}

func NewAttachedMsg(msg string) *Error {
	return &Error{
		code: attachedMsgError,
		msg:  msg,
	}
}

func (e *Error) Code() uint32 {
	return e.code
}

func (e *Error) Message() string {
	return e.msg
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.msg)
}
