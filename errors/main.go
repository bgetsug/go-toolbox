package errors

import (
	"fmt"

	"github.com/go-stack/stack"
)

type Code string

const (
	BindFailedCode          Code = "BIND_FAILED"
	JSONUnmarshalFailedCode Code = "JSON_UNMARSHAL_FAILED"
	JSONSyntaxErrorCode     Code = "JSON_SYNTAX_ERROR"

	UnknownCode Code = "UNKNOWN"

	ValidationFailedCode Code = "VALIDATION_FAILED"
	PanicRecoveryCode    Code = "PANIC_RECOVERY"
)

type Error struct {
	err     error
	code    Code
	message string
	stack   CallStack
}

func (e Error) Error() string {
	// annotate
	if e.message != "" && e.err != nil {
		return fmt.Sprintf("%s: %s", e.message, e.err.Error())
	}

	// fallback
	if e.err == nil {
		return e.message
	}

	// default to wrapped error message
	return e.err.Error()
}

// Gets the error contained inside the Error.
func (e Error) Err() error {
	return e.err
}

// Gets the callstack for the Error.
func (e Error) Stack() CallStack {
	return e.stack
}

// Gets the code from the Error's CodedError.
func (e Error) Code() Code {
	return e.code
}

// Creates a new Error.
func New(code Code, msg string) error {
	callStack := stack.Trace()
	return Error{nil, code, msg, CallStack(callStack.TrimBelow(callStack[1]).TrimRuntime())}
}

// Wraps an existing error so that we can provide a callstack and our own CodedError with it.
func Wrap(err error, code Code, msg string) error {
	callStack := stack.Trace()
	return Error{err, code, msg, CallStack(callStack.TrimBelow(callStack[1]).TrimRuntime())}
}

// Gets the error that is contained inside the given error. If err is not a Error, err itself is returned.
func WrappedError(err error) error {
	if withCode, ok := err.(Error); ok {
		return withCode.err
	}

	return err
}

// Gets the callstack associated with the error, if any.
func Stack(err error) CallStack {
	if withCode, ok := err.(Error); ok {
		return withCode.stack
	}

	return nil
}

// Gets the error code associated with the error. Returns UnknownCode by default.
func CodeOrDefault(err error) Code {
	if withCode, ok := err.(Error); ok {
		return withCode.code
	}

	return UnknownCode
}
