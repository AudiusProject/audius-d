package logger

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
)

/**
Simplistic logger for audius-d.
*/

var (
	stdoutLogger *log.Logger
	stderrLogger *slog.Logger
)

func init() {
	// configure loggers
	stdoutLogger = log.New(os.Stdout, "", 0)
	stderrLogger = slog.Default()
}

func Out(v ...any) {
	stdoutLogger.Print(v...)
}

func Info(msg string, v ...any) {
	stderrLogger.Info(msg, v...)
}

func Infof(format string, v ...any) {
	fmsg := fmt.Sprintf(format, v...)
	stderrLogger.Info(fmsg)
}

func Debug(msg string, v ...any) {
	stderrLogger.Debug(msg, v...)
}

func Warn(msg string, v ...any) {
	stderrLogger.Warn(msg, v...)
}

func ErrorF(format string, v ...any) error {
	emsg := fmt.Sprintf(format, v...)
	return Error(emsg)
}

// you can return this log as well
// to get log.Fatal effects
func Error(msg any, v ...any) error {
	switch m := msg.(type) {
	case string:
		stderrLogger.Error(m, v...)
		return errors.New(m)
	case error:
		stderrLogger.Error(m.Error(), v...)
		return m
	default:
		msgs := []any{m}
		vs := append(msgs, v...)
		stderrLogger.Error("", vs...)
		return errors.New("unexpected error")
	}
}
