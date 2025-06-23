// Package logger provides a logging interface and implementations for formatted log output.
package logger

// LogLevel specifies a level from which logs are printed.
type logLevel int32

const (
	TRACE logLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	PANIC
)

// Logger interface defines methods for logging at various levels.
type Logger interface {
	Trace(v ...any)
	Tracef(format string, v ...any)
	Debug(v ...any)
	Debugf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	Panic(v ...any)
	Panicf(format string, v ...any)
}
