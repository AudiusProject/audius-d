package logger

import (
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

func Error(msg any, v ...any) {
	switch m := msg.(type) {
	case string:
		stderrLogger.Error(m, v...)
	case error:
		stderrLogger.Error(m.Error(), v...)
	default:
		msgs := []any{m}
		vs := append(msgs, v...)
		stderrLogger.Error("", vs...)
	}
}
