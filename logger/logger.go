package logger

import (
	"context"
)

type LoggerInterface interface {
	// AddAttrToCtx
	//
	// If `ctx` is `nil`, `context.Background()` is used.
	AddAttrToCtx(ctx context.Context, key string, val string) context.Context

	// DebugFf
	//
	// If `ctx` is `nil`, `context.Background()` is used.
	//
	// Uses fmt.Sprintf() to make a message.
	DebugFf(ctx context.Context, msg string, msgArgs ...any)

	// InfoFf
	//
	// If `ctx` is `nil`, `context.Background()` is used.
	//
	// Uses fmt.Sprintf() to make a message.
	InfoFf(ctx context.Context, msg string, msgArgs ...any)

	// WarnFf
	//
	// If `ctx` is `nil`, `context.Background()` is used.
	//
	// Uses fmt.Sprintf() to make a message.
	WarnFf(ctx context.Context, msg string, msgArgs ...any)

	// ErrorFf
	//
	// If `ctx` is `nil`, `context.Background()` is used.
	//
	// Uses fmt.Sprintf() to make a message.
	ErrorFf(ctx context.Context, msg string, msgArgs ...any)

	// PanicFf
	//
	// Extra method to simplify logging of the panic cases.
	//
	// If `ctx` is `nil`, `context.Background()` is used.
	//
	// Uses fmt.Sprintf() to make a message.
	PanicFf(ctx context.Context, panicValue any, debugStack []byte, msg string, msgArgs ...any)
}
