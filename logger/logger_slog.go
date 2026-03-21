package logger

import (
	"context"
	"fmt"
	"log/slog"
)

// #####################################################################################################################
// LOGGER
// #####################################################################################################################

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type SlogLogger struct {
	l *slog.Logger
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewSlogLogger(h slog.Handler) SlogLogger {
	h = slogLogger_NewAttributedSlogHandler(h)

	return SlogLogger{
		l: slog.New(h),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// AddAttrToCtx
//
// If `ctx` is `nil`, `context.Background()` is used.
func (l SlogLogger) AddAttrToCtx(ctx context.Context, key string, val string) context.Context {
	ctx = ctxOrBackground(ctx)
	return slogLogger_AddAttrToCtx(ctx, key, val)
}

// DebugFf
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func (l SlogLogger) DebugFf(ctx context.Context, msg string, msgArgs ...any) {
	ctx = ctxOrBackground(ctx)
	l.l.DebugContext(ctx, fmt.Sprintf(msg, msgArgs...))
}

// InfoFf
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func (l SlogLogger) InfoFf(ctx context.Context, msg string, msgArgs ...any) {
	ctx = ctxOrBackground(ctx)
	l.l.InfoContext(ctx, fmt.Sprintf(msg, msgArgs...))
}

// WarnFf
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func (l SlogLogger) WarnFf(ctx context.Context, msg string, msgArgs ...any) {
	ctx = ctxOrBackground(ctx)
	l.l.WarnContext(ctx, fmt.Sprintf(msg, msgArgs...))
}

// ErrorFf
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func (l SlogLogger) ErrorFf(ctx context.Context, msg string, msgArgs ...any) {
	ctx = ctxOrBackground(ctx)
	l.l.ErrorContext(ctx, fmt.Sprintf(msg, msgArgs...))
}

// PanicFf
//
// Extra method to simplify logging of the panic cases.
//
// If `ctx` is `nil`, `context.Background()` is used.
//
// Uses fmt.Sprintf() to make a message.
func (l SlogLogger) PanicFf(ctx context.Context, panicValue any, debugStack []byte, msg string, msgArgs ...any) {
	ctx = ctxOrBackground(ctx)

	pvStr := ""
	switch pv := panicValue.(type) {
	case error:
		pvStr = pv.Error()
	case string:
		pvStr = pv
	case fmt.Stringer:
		pvStr = pv.String()
	default:
		pvStr = fmt.Sprintf("%#v", pv)
	}

	l.l.ErrorContext(
		ctx,
		fmt.Sprintf(msg, msgArgs...),
		slog.String("stack", string(debugStack)),
		slog.String("panic", pvStr),
	)
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

// SlogLogger
//
// Returns a pointer to the internal `slog.Logger`.
//
// Can be used, for example, to `slog.SetDefault()`.
func (l SlogLogger) SlogLogger() *slog.Logger {
	return l.l
}

// #####################################################################################################################
// SLOG HANDLER
// #####################################################################################################################

const slogLogger_AttrsCtxKey = "logger.slogLoggerInternalAttrs"

type slogLogger_AttributedSlogHandler struct {
	slog.Handler
}

func slogLogger_NewAttributedSlogHandler(h slog.Handler) *slogLogger_AttributedSlogHandler {
	return &slogLogger_AttributedSlogHandler{Handler: h}
}

func slogLogger_AddAttrToCtx(ctx context.Context, key string, val string) context.Context {
	attrs, ok := ctx.Value(slogLogger_AttrsCtxKey).([]slog.Attr)
	if attrs == nil || !ok {
		attrs = make([]slog.Attr, 0)
	}

	attrs = append(attrs, slog.String(key, val))

	return context.WithValue(ctx, slogLogger_AttrsCtxKey, attrs)
}

func (h *slogLogger_AttributedSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	if extraAttrs, ok := ctx.Value(slogLogger_AttrsCtxKey).([]slog.Attr); ok {
		r.AddAttrs(extraAttrs...)
	}

	return h.Handler.Handle(ctx, r)
}

func (h *slogLogger_AttributedSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return slogLogger_NewAttributedSlogHandler(h.Handler.WithAttrs(attrs))
}

func (h *slogLogger_AttributedSlogHandler) WithGroup(name string) slog.Handler {
	return slogLogger_NewAttributedSlogHandler(h.Handler.WithGroup(name))
}

// #####################################################################################################################
