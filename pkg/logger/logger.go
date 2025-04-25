package logger

// LogLevel specifies a level from which logs are printed.
type LogLevel int32

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

// Logger interface defines methods for logging at various levels.
type Logger interface {
	SetLogLevel(level LogLevel)
	GetLogLevel() LogLevel

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
	Fatal(v ...any)
	Fatalf(format string, v ...any)
}
