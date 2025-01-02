package log

import (
	"context"
	"fmt"
	"strings"
	"time"

	"goeasy.dev/observability/log/loglevel"
)

type Sink interface {
	Write(p []byte) (n int, err error)
	Flush() error
}

type Logger struct {
	sink     Sink
	logLevel loglevel.LogLevel
}

const txtLogFormat = `time=%s level=%s message="%s"`

var messageReplacer = strings.NewReplacer("\n", "\\n", "\t", "\\t", "\r", "\\r", "\"", "\\\"")

func NewLogger(sink Sink, logLevel loglevel.LogLevel) *Logger {
	return &Logger{
		sink:     sink,
		logLevel: logLevel,
	}
}

func (l *Logger) Trace(args ...interface{}) {
	l.Log(loglevel.TRACE, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Log(loglevel.DEBUG, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Log(loglevel.INFO, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Log(loglevel.WARN, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Log(loglevel.ERROR, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	panic(l.Log(loglevel.FATAL, args...))
}

func (l *Logger) Log(level loglevel.LogLevel, args ...interface{}) string {
	if level < l.logLevel {
		return ""
	}

	if len(args) == 0 {
		return ""
	}

	if _, ok := args[0].(context.Context); ok {
		args = args[1:]
	}

	var msg string
	if len(args) == 1 {
		msg = args[0].(string)
	} else if len(args) > 1 {
		msg = fmt.Sprintf(args[0].(string), args[1:]...)
	}

	fmt.Fprintf(l.sink, txtLogFormat+"\n", time.Now().UTC().Format(time.RFC3339), level.String(), messageReplacer.Replace(msg))

	return msg
}
