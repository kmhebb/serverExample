package cloud

import (
	"net/http"
	"runtime"
	"strings"
)

// ErrorKind specifies a class of errors, like internal errors or bad request
// errors.
type ErrorKind string

// Error kinds are recognized in the framework. These should only be augmented
// after a review to ensure consistency. If a new kind is
// added, the Error.Code method will also need to be updated to ensure we're
// returning an appropriate HTTP status code for the new error kind. It may also
// make sense to update the Error.Message method with a suitable default for the
// new kind, if the overall default isn't appropriate.
const (
	ErrKindAuthenticate ErrorKind = "authenticate"
	ErrKindBadRequest   ErrorKind = "bad_request"
	ErrKindConflict     ErrorKind = "conflict"
	ErrKindExternal     ErrorKind = "external"
	ErrKindForbidden    ErrorKind = "forbidden"
	ErrKindInternal     ErrorKind = "internal"
	ErrKindInvalid      ErrorKind = "invalid"
	ErrKindNotFound     ErrorKind = "not_found"
	ErrKindTodo         ErrorKind = "todo"
)

// Error is a cloud-domain error. The NewError constructor should be used to
// create Errors. The zero value is a vague internal error.
type Error struct {
	kind ErrorKind
	op   string
	err  error
	msg  string
}

// ErrOpts are options for adding additional information to an Error. The zero
// value is an acceptable set of options.
type ErrOpts struct {
	// Kind is the class of errors the error falls under. The default is
	// ErrKindInternal. This should never be set to anything other than one of
	// the ErrorKind constants defined in this package. Setting this value
	// affects the output of the Error's Code and Error methods, and also
	// affects the default output of the Error's Message method if a custom
	// Message is not specified.
	Kind ErrorKind

	// Cause is the underlying error that resulted in the creation of our
	// domain-level error. Setting this value allows for unwrapping of the
	// domain-level error to this error, and also affects the output of the
	// Error's Error method.
	Cause error

	// Message is a user-facing message associated with the error. Setting this
	// value affects the output of the Error's Message method. If not set, a
	// default message based on the error kind is used.
	Message string
}

// NewError constructs a new Error.
func NewError(opts ErrOpts) *Error {
	e := &Error{
		kind: opts.Kind,    // If empty, the Kind method will handle the default.
		err:  opts.Cause,   // Nil is a valid value.
		msg:  opts.Message, // If empty, the Message method will handle the default.
	}

	// Add trace information if we can.
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return e
	}
	if f := runtime.FuncForPC(pc); f != nil {
		e.op = f.Name()
	}

	return e
}

// Code returns the HTTP status code associated with the error's kind.
func (e Error) Code() int {
	switch e.kind {
	case ErrKindBadRequest, ErrKindConflict, ErrKindForbidden, ErrKindInvalid, ErrKindNotFound:
		return http.StatusBadRequest
	case ErrKindAuthenticate:
		return http.StatusUnauthorized
	case ErrKindTodo:
		return http.StatusNotImplemented
	case ErrKindExternal:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

// Kind returns the error's kind.
func (e Error) Kind() ErrorKind {
	if e.kind != "" {
		return e.kind
	}

	return ErrKindInternal
}

// Message returns the user-facing message associated with the error.
func (e Error) Message() string {
	if e.msg != "" {
		return e.msg
	}

	switch e.kind {
	case ErrKindBadRequest, ErrKindConflict, ErrKindForbidden, ErrKindInvalid, ErrKindNotFound:
		return "We were unable to process your request."
	case ErrKindAuthenticate:
		return "An authentication error occurred."
	case ErrKindTodo:
		return "This feature is not yet ready."
	case ErrKindExternal:
		return "An unexpected external error occurred."
	default:
		return "An unexpected error occurred."
	}
}

// Error implements error.
func (e Error) Error() string {
	// We can have up to 3 parts to our error.
	parts := make([]string, 0, 3)

	if e.op != "" {
		parts = append(parts, e.op)
	}

	// Kind is the only guaranteed part of our error, since we use the Kind
	// method instead of the raw value.
	parts = append(parts, string(e.Kind()))

	if e.err != nil {
		parts = append(parts, e.err.Error())
	}

	return strings.Join(parts, ": ")
}

// Unwrap returns the cause error if one was specified, or nil otherwise.
func (e Error) Unwrap() error {
	return e.err
}
