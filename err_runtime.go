package std

import (
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
// Usage example:
//
//	package example
//
//	func SomeFunc() error {
//		v, err := ...
//		if err != nil {
//			return std.WrapErrorToRuntime(err, "example", "SomeFunc")
//		}
//		...
//	}
//
//	type MyType struct { ... }
//
//	func (m *MyType) SomeMethod() (..., error) {
//		...
//		if err != nil {
//			return std.WrapErrorToRuntime(err, m, "SomeMethod")
//		}
//		...
//	}
//
// See WrapErrorToRuntime.
type ErrorRuntime struct {
	err error
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewErrorRuntimeFf
//
// Panics in case of empty arguments:
//   - msg
func NewErrorRuntimeFf(msg string, msgArgs ...any) ErrorRuntime {
	if msg == "" {
		panic("`msg` must not be empty")
	}

	return ErrorRuntime{err: fmt.Errorf(msg, msgArgs...)}
}

// WrapErrorToRuntime
//
// Adds this kind of prefix to the error to imitate stack trace:
// "{{`methodOwner`}}.{{`methodName`}}/{{joined by "/" elements of the `methodInfo`}}".
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
func WrapErrorToRuntime(err error, methodOwner any, methodName string, methodInfo ...string) ErrorRuntime {
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
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e ErrorRuntime) Error() string {
	return e.err.Error()
}

func (e ErrorRuntime) Unwrap() error {
	return e.err
}

// ---------------------------------------------------------------------------------------------------------------------
