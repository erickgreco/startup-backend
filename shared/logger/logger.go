package logger

import (
	"log/slog"
	"os"
)

var l *slog.Logger

func Init(level string) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo
	}

	l = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(l)
}

func Info(msg string, args ...any)  { l.Info(msg, args...) }
func Debug(msg string, args ...any) { l.Debug(msg, args...) }
func Warn(msg string, args ...any)  { l.Warn(msg, args...) }
func Error(msg string, err error, args ...any) {
	l.Error(msg, append([]any{"err", err}, args...)...)
}
