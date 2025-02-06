package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"
)

const (
	envLocal = "local"
	envProd  = "prod"

	layout = "2006-01-02 15:04:05.000"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true, ReplaceAttr: formatTime}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

// formatTime changes time representation in "time" record attribute
func formatTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == "time" {
		oldTimeVal := a.Value.String()[:len(layout)]

		oldTime, err := time.Parse(layout, oldTimeVal)
		if err != nil {
			log.Fatalf("cannot parse time attribute: %s", err.Error())
		}

		newTime := fmt.Sprintf(
			"{%d:%d:%d}:[%d:%d:%d.%d]",
			oldTime.Day(), oldTime.Month(), oldTime.Year(),
			oldTime.Hour(), oldTime.Minute(), oldTime.Second(), oldTime.Nanosecond()/1e7,
		)

		a.Value = slog.StringValue(newTime)
	}

	return a
}
