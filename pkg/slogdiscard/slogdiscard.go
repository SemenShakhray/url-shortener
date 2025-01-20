package slogdiscard

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{}

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler)
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	//ignore the log entry
	return nil
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	//returns false because the log entry is ignored
	return false
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	//returns the same handler, since there are no attributes to save.
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	//returns the same handler, since there is no group to save.
	return h
}
