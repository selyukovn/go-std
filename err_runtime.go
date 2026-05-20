package std

import (
	"errors"
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

// ErrorRuntime
//
// Client code should be able to distinguish business-logic errors from others.
// Moreover, each type of business-logic error may be handled differently.
// Runtime errors are typically handled uniformly (e.g: logged, reported, translated into an HTTP 500 response).
//
// Go does not allow us to mark business-logic errors and then check that mark in the client code in a simple way,
// like some other OOP languages do (e.g., through inheritance and try‑catch statements),
// but it allows to get the same result in another way --
// by wrapping all non-business-logic errors into a specific error type.
//
// I.e:
//
//	try { ... }
//	catch (BusinessLogicException) { /* ... */ }
//	catch (Exception) { /* any other error -- in this case non-business-logic */ }
//
// -->
//
//	switch err.(type) {
//	case nil:
//	case std.ErrorRuntime: /* any other error -- in this case non-business-logic */
//	default: /* ... */
//	}
//
// --
//
// For some cases it may be important to distinguish temporary runtime errors from permanent (e.g. for retry strategies).
// ErrorRuntime has a special flag for this -- IsTemporary() / IsPermanent().
//
// --
//
// Usage example:
//
//	package example
//
//	func SomeFunc() error {
//		v, err := ...
//		if err != nil {
//			return std.ErrToRuntime(err, "example", "SomeFunc")
//		}
//		...
//	}
//
//	type MyType struct { ... }
//
//	func (m *MyType) SomeMethod() (..., error) {
//		...
//		if err != nil {
//			return std.ErrToRuntime(err, m, "SomeMethod")
//		}
//		...
//	}
type ErrorRuntime struct {
	err error
	tmp bool
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewErrorRuntimeFf
//
// DEPRECATED: use NewErrorRuntimePermanentFf() or NewErrorRuntimeTemporaryFf() instead.
func NewErrorRuntimeFf(msg string, msgArgs ...any) ErrorRuntime {
	return NewErrorRuntimePermanentFf(msg, msgArgs...)
}

// NewErrorRuntimePermanentFf
//
// See ErrorRuntime for details.
func NewErrorRuntimePermanentFf(msg string, msgArgs ...any) ErrorRuntime {
	return newErrorRuntimeFf(false, msg, msgArgs...)
}

// NewErrorRuntimeTemporaryFf
//
// See ErrorRuntime for details.
func NewErrorRuntimeTemporaryFf(msg string, msgArgs ...any) ErrorRuntime {
	return newErrorRuntimeFf(true, msg, msgArgs...)
}

func newErrorRuntimeFf(isTmp bool, msg string, msgArgs ...any) ErrorRuntime {
	if msg == "" {
		panic("`msg` must not be empty")
	}

	return ErrorRuntime{
		err: fmt.Errorf(msg, msgArgs...),
		tmp: isTmp,
	}
}

// --

// WrapErrorToRuntime
//
// DEPRECATED: use ErrToRuntime() instead.
func WrapErrorToRuntime(err error, methodOwner any, methodName string, methodInfo ...string) ErrorRuntime {
	return ErrToRuntime(err, methodOwner, methodName, methodInfo...)
}

// ErrToRuntime
//
// Wraps provided error into ErrorRuntime.
//
// Adds this kind of prefix to the error to imitate stack trace:
// "{{`methodOwner`}}.{{`methodName`}}/{{joined by "/" elements of the `methodInfo`}}".
//
// If `err` is a temporary ErrorRuntime (wrapped or not -- errors.As() is used as a checker),
// result error also will be marked as temporary, otherwise -- permanent.
// To mark result error by force use ErrToRuntimePermanent() or ErrToRuntimeTemporary().
//
// Panics in case of empty arguments:
//   - err
//   - methodOwner
//   - methodName
//
// Arguments:
//   - `methodOwner` -- string or any method owner instance.
//     {{`methodOwner`}} will be replaced by the string value or method owner instance type accordingly.
//     String value is allowed, for example, to accept package names, when a regular function is called, not a method.
//   - `methodName` -- name of the method / function, that was called.
//   - `methodInfo` -- any additional info.
func ErrToRuntime(err error, methodOwner any, methodName string, methodInfo ...string) ErrorRuntime {
	isTmp := false

	var vErr ErrorRuntime
	if errors.As(err, &vErr) {
		isTmp = vErr.tmp
	}

	return errorToRuntime(isTmp, err, methodOwner, methodName, methodInfo...)
}

// ErrToRuntimePermanent
//
// Same as ErrToRuntime, but marks result error as permanent, even if source error was temporary.
func ErrToRuntimePermanent(err error, methodOwner any, methodName string, methodInfo ...string) ErrorRuntime {
	return errorToRuntime(false, err, methodOwner, methodName, methodInfo...)
}

// ErrToRuntimeTemporary
//
// Same as ErrToRuntime, but marks result error as temporary, even if source error was permanent.
func ErrToRuntimeTemporary(err error, methodOwner any, methodName string, methodInfo ...string) ErrorRuntime {
	return errorToRuntime(true, err, methodOwner, methodName, methodInfo...)
}

func errorToRuntime(isTmp bool, err error, methodOwner any, methodName string, methodInfo ...string) ErrorRuntime {
	if err == nil {
		panic("`err` must not be nil")
	}

	methodOwnerStr := ""
	if mos, ok := methodOwner.(string); ok {
		if mos == "" {
			panic("`methodOwner` must not be empty")
		}
		methodOwnerStr = mos
	} else {
		if methodOwner == nil {
			panic("`methodOwner` must not be nil")
		}
		methodOwnerStr = fmt.Sprintf("%T", methodOwner)
	}

	// --

	messagePrefix := fmt.Sprintf(
		"%s%s/%s",
		methodOwnerStr+Ternary[string](methodOwnerStr == "", "", "."),
		methodName,
		strings.Join(methodInfo, "/"),
	)

	return ErrorRuntime{
		err: fmt.Errorf("%s: %w", messagePrefix, err),
		tmp: isTmp,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

// IsPermanent
//
// Same as "not IsTemporary()".
//
// See NewErrorRuntimePermanentFf().
// See ErrToRuntimePermanent().
func (e ErrorRuntime) IsPermanent() bool {
	return !e.tmp
}

// IsTemporary
//
// See NewErrorRuntimeTemporaryFf().
// See ErrToRuntimeTemporary().
func (e ErrorRuntime) IsTemporary() bool {
	return e.tmp
}

func (e ErrorRuntime) Error() string {
	return e.err.Error()
}

func (e ErrorRuntime) Unwrap() error {
	return e.err
}

// ---------------------------------------------------------------------------------------------------------------------
