package loglevel

//go:generate stringer -type=LogLevel
type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)
