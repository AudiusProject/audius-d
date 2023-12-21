package logger

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

/**
Simplistic logger for audius-d.
*/

var (
	stdoutLogger *log.Logger
	stderrLogger *slog.Logger

	cliHandlerOpts = slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				return a
			} else { // Remove everything but the message
				return slog.Attr{}
			}
		},
	}
)

func init() {
	// configure loggers
	stdoutLogger = log.New(os.Stdout, "", 0)
	stderrLogger = slog.New(
		NewCliHandler(slog.NewTextHandler(os.Stderr, &cliHandlerOpts)),
	)
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
func Error(values ...any) error {
	messages := []string{}
	for _, v := range values {
		switch t := v.(type) {
		case string:
			messages = append(messages, t)
		case error:
			messages = append(messages, t.Error())
		default:
		}
	}
	message := strings.Join(messages, " ")
	stderrLogger.Error(message)
	return errors.New(message)
}
