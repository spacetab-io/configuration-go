package config

import (
	"context"
)

type Logger interface {
	Debug(context.Context, string, ...any)
	Warn(context.Context, string, ...any)
	Error(context.Context, string, ...any)
}

var _ Logger = (*noOpLogger)(nil)

type noOpLogger struct{}

func (n noOpLogger) Debug(ctx context.Context, msg string, args ...any) {
	return
}

func (n noOpLogger) Error(ctx context.Context, msg string, args ...any) {
	return
}

func (n noOpLogger) Warn(ctx context.Context, msg string, args ...any) {
	return
}
