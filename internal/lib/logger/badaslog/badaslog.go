package badaslog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

const (
	reset = "\033[0m"

	// black        = 30
	// red          = 31
	// green        = 32
	// yellow       = 33
	// blue         = 34
	magenta   = 35
	cyan      = 36
	lightGray = 37
	darkGray  = 90
	lightRed  = 91
	// lightGreen   = 92
	lightYellow = 93
	// lightBlue    = 94
	// lightMagenta = 95
	// lightCyan    = 96
	white = 97

	timeFormat = "[15:04:05.000]"
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%dm%s%s", colorCode, v, reset)
}

// Handler is a wrap on slog default handler.
// It should implement slog.Handler interface.
type Handler struct {
	// a nested handler which we wrap
	h slog.Handler
	// captures the output of nested handler
	b *bytes.Buffer
	// guarantee of thread safe access
	m *sync.Mutex
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(darkGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	bytes, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	fmt.Print(
		colorize(magenta, r.Time.Format(timeFormat)),
		level,
		colorize(white, r.Message), " ",
		colorize(darkGray, string(bytes)),
	)

	return nil
}

func NewHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	b := &bytes.Buffer{}

	return &Handler{
		b: b,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: supressDefaults(opts.ReplaceAttr),
		}),
		m: &sync.Mutex{},
	}
}

// computeAttrs computes attributes of inner logger and deserialize it into map
// for further writing it as a colourful JSON.
func (h *Handler) computeAttrs(ctx context.Context, r slog.Record) ([]byte, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return []byte{}, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	bytesRec := h.b.Bytes()

	var obj map[string]json.RawMessage
	err := json.Unmarshal([]byte(bytesRec), &obj)
	if err != nil {
		return []byte{}, fmt.Errorf("error when unmarshallign json: %w", err)
	}

	// check whether json is empty. In this case we shouldn't print '{}'
	if len(obj) == 0 {
		return []byte("\n"), nil
	}

	var b bytes.Buffer
	err = json.Indent(&b, bytesRec, "", "  ")
	if err != nil {
		return []byte{}, fmt.Errorf("error when indenting json: %w", err)
	}

	return b.Bytes(), nil
}

// supressDefaults exclude time, level and message from nested logger. It works as a middleware
func supressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}

		if next == nil {
			return a
		}

		return next(groups, a)
	}
}
