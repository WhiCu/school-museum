package logger

import (
	"log/slog"

	"github.com/WhiCu/school-museum/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func GetLogger(cfg *config.LoggerConfig) *slog.Logger {
	h := make([]slog.Handler, 0, 2)

	if cfg.Path != "" {
		logFile := &lumberjack.Logger{
			Filename:   cfg.Path,
			MaxSize:    cfg.Size,
			MaxAge:     cfg.Age,
			MaxBackups: cfg.Backups,
			LocalTime:  true,
			Compress:   cfg.Compress,
		}

		h = append(h, slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	h = append(h, MustInitLogger(cfg.Level))

	return slog.New(slog.NewMultiHandler(h...))
}
