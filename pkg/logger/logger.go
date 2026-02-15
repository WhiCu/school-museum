package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gitlab.com/greyxor/slogor"
)

func MustInitLogger(level string) slog.Handler {
	const op = "logger.MustInitLogger"

	var log slog.Handler

	out := os.Stdout

	switch level {
	case "debug":
		log = slogor.NewHandler(out, slogor.SetLevel(slog.LevelDebug), slogor.SetTimeFormat(time.ANSIC))
	case "info":
		log = slogor.NewHandler(out, slogor.SetLevel(slog.LevelInfo), slogor.SetTimeFormat(time.ANSIC))
	case "warn":
		log = slogor.NewHandler(out, slogor.SetLevel(slog.LevelWarn), slogor.SetTimeFormat(time.ANSIC))
	case "error":
		log = slogor.NewHandler(out, slogor.SetLevel(slog.LevelError), slogor.SetTimeFormat(time.ANSIC))
	default:
		panic(fmt.Sprintf("%s: level - {%s} not in {debug, info, warn, error}", op, level))
	}
	return log
}
