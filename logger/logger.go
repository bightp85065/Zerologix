package logger

import (
	"context"
	"fmt"
	"sync"
	"time"
	"zerologix/constant"

	"golang.org/x/exp/slog"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

var getXlogWithCtxOnce sync.Once

// Even ctx is nil, it should be safe because log.WithCtx have nil guard.
type log struct {
	ctx context.Context
}

func Log(ctx context.Context) *log {
	now := time.Now()
	// Background context (context.Background()) shares the same instance, so need to "split" it.
	if ctx == nil || ctx == context.Background() {
		ctx = context.WithValue(context.Background(), constant.ContextKeyTraceID, fmt.Sprintf("local-%v", now.UnixMilli()))
	}

	if t := ctx.Value(constant.ContextKeyTraceID); t != nil {
		return &log{ctx: ctx}
	} else {
		t = fmt.Sprintf("local-%v", now.UnixMilli())
		return &log{context.WithValue(ctx, constant.ContextKeyTraceID, t)}
	}
}

func (l *log) logInternal(message string, level LogLevel) {
	fmt.Printf("[%s] [%s] %s\n", time.Now(), l.levelToString(level), message)
}

func (l *log) levelToString(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARNING"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (l *log) Debugf(msg string, args ...interface{}) {
	l.logInternal(fmt.Sprintf(msg, args...), LogLevelDebug)
}

func (l *log) Infof(msg string, args ...interface{}) {
	l.logInternal(fmt.Sprintf(msg, args...), LogLevelInfo)
}

func (l *log) Warnf(msg string, args ...interface{}) {
	l.logInternal(fmt.Sprintf(msg, args...), LogLevelWarning)
}

func (l *log) Errorf(msg string, args ...interface{}) {
	l.logInternal(fmt.Sprintf(msg, args...), LogLevelError)
}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, slog.Logger{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	return ctx.Value(slog.Logger{}).(*slog.Logger)
}
