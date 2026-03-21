package logger

import "context"

// ---------------------------------------------------------------------------------------------------------------------

const defaultLoggerIsNotSetPanicMessage = "default logger is not set"

var defaultLogger LoggerInterface = nil

// SetDefault
//
// Sets the global logger for package-level shortcuts usage.
//
// Panics in case of empty arguments.
func SetDefault(l LoggerInterface) {
	if l == nil {
		panic("`logger.SetDefault` expects non-nil arguments")
	}

	defaultLogger = l
}

// ---------------------------------------------------------------------------------------------------------------------

// AddAttrToCtx
//
// Panics if default logger is not set -- see SetDefault.
//
// If `ctx` is `nil`, `context.Background()` is used.
func AddAttrToCtx(ctx context.Context, key string, val string) context.Context {
	if defaultLogger == nil {
		panic(defaultLoggerIsNotSetPanicMessage)
	}

	return defaultLogger.AddAttrToCtx(ctx, key, val)
}

// DebugFf
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func DebugFf(ctx context.Context, msg string, msgArgs ...any) {
	if defaultLogger == nil {
		panic(defaultLoggerIsNotSetPanicMessage)
	}

	defaultLogger.DebugFf(ctx, msg, msgArgs...)
}

// InfoFf
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func InfoFf(ctx context.Context, msg string, msgArgs ...any) {
	if defaultLogger == nil {
		panic(defaultLoggerIsNotSetPanicMessage)
	}

	defaultLogger.InfoFf(ctx, msg, msgArgs...)
}

// WarnFf
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func WarnFf(ctx context.Context, msg string, msgArgs ...any) {
	if defaultLogger == nil {
		panic(defaultLoggerIsNotSetPanicMessage)
	}

	defaultLogger.WarnFf(ctx, msg, msgArgs...)
}

// ErrorFf
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func ErrorFf(ctx context.Context, msg string, msgArgs ...any) {
	if defaultLogger == nil {
		panic(defaultLoggerIsNotSetPanicMessage)
	}

	defaultLogger.ErrorFf(ctx, msg, msgArgs...)
}

// PanicFf
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func PanicFf(ctx context.Context, panicValue any, debugStack []byte, msg string, msgArgs ...any) {
	if defaultLogger == nil {
		panic(defaultLoggerIsNotSetPanicMessage)
	}

	defaultLogger.PanicFf(ctx, panicValue, debugStack, msg, msgArgs...)
}
