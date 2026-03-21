package logger

import "context"

func ctxOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}
