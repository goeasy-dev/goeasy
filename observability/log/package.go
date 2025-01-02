package log

import (
	"os"

	"goeasy.dev/observability/log/loglevel"
	"goeasy.dev/observability/log/sinks/writersink"
)

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(writersink.NewWriterSink(os.Stdout), loglevel.INFO)
}

func Trace(args ...interface{}) {
	defaultLogger.Trace(args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}
