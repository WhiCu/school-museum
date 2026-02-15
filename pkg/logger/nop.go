package logger

import (
	"context"
	"log/slog"
)

type NOPHandler struct {
}

func (m NOPHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (m NOPHandler) Handle(_ context.Context, r slog.Record) error {
	return nil
}

func (m NOPHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return m
}

func (m NOPHandler) WithGroup(_ string) slog.Handler {
	return m
}

func NewNOPHandler() NOPHandler {
	return NOPHandler{}
}

func NewNOPSlog() *slog.Logger {
	return slog.New(NewNOPHandler())
}
