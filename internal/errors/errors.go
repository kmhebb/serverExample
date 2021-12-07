package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

const DefaultMessage = "An unexpected error occurred. Please try again or contact Ken."

type Kind string

const (
	EINVALID  Kind = "invalid"
	EINTERNAL Kind = "internal"
	ETODO     Kind = "todo"
)

func codeForKind(k Kind) int {
	switch k {
	case EINVALID:
		return http.StatusBadRequest
	case ETODO:
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}

type Mapper interface {
	Map() map[string]interface{}
}

func E(k Kind) *Error {
	err := &Error{kind: k}
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return err
	}
	if f := runtime.FuncForPC(pc); f != nil {
		err.op = f.Name()
	}
	return err
}

type Error struct {
	kind  Kind
	op    string
	err   error
	msg   string
	extra Mapper
}

func (e Error) Code() int {
	return codeForKind(e.kind)
}

func (e Error) Kind() Kind {
	return e.kind
}

func (e Error) Message() string {
	if e.msg == "" {
		return DefaultMessage
	}
	return e.msg
}

func (e Error) Error() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s: %s", e.op, e.kind)
	if e.err != nil {
		fmt.Fprintf(&sb, ": %s", e.err)
	}
	return sb.String()
}

func (e Error) WithMessage(msg string, args ...interface{}) *Error {
	e.msg = fmt.Sprintf(msg, args...)
	return &e
}

func (e Error) WithError(err error) *Error {
	e.err = err
	return &e
}

func (e Error) WithErrorf(msg string, args ...interface{}) *Error {
	e.err = fmt.Errorf(msg, args...)
	return &e
}

func (e Error) WithExtra(extra Mapper) *Error {
	e.extra = extra
	return &e
}
