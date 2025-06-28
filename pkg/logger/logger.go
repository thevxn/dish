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

var logColors = map[logLevel]string{
	TRACE: "\033[34m", // Blue
	DEBUG: "\033[36m", // Cyan
	INFO:  "\033[32m", // Green
	WARN:  "\033[33m", // Yellow
	ERROR: "\033[31m", // Red
	PANIC: "\033[35m", // Magenta
}

var logPrefixes = map[logLevel]string{
	TRACE: "[ TRACE ]: ",
	DEBUG: "[ DEBUG ]: ",
	INFO:  "[ INFO ]: ",
	WARN:  "[ WARN ]: ",
	ERROR: "[ ERROR ]: ",
	PANIC: "[ PANIC ]: ",
}

func (l logLevel) Color() string {
	if color, exists := logColors[l]; exists {
		return color
	}
	return "\033[0m" // Default color (reset)
}

func (l logLevel) Prefix(withColor bool) string {
	prefix, prefixExists := logPrefixes[l]

	if !prefixExists {
		return "[ UNKNOWN ]: "
	}

	if withColor {
		return l.Color() + prefix + "\033[0m"
	}

	return prefix
}

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
