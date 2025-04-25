package logger

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
)

// ConsoleLogger logs output to stderr.
type ConsoleLogger struct {
	stdLogger *log.Logger
	logLevel  atomic.Int32
}

// NewConsoleLogger creates a new ConsoleLogger instance,
// If verbose is true, log level is set to TRACE (otherwise to INFO).
func NewConsoleLogger(verbose bool) *ConsoleLogger {
	l := &ConsoleLogger{
		stdLogger: log.New(os.Stderr, "", log.LstdFlags),
	}

	if verbose {
		l.SetLogLevel(TRACE)
	} else {
		l.SetLogLevel(INFO)
	}

	return l
}

// SetLogLevel safely sets logger's log level.
func (l *ConsoleLogger) SetLogLevel(level LogLevel) {
	l.logLevel.Store(int32(level))
}

// GetLogLevel safely returns current log level.
func (l *ConsoleLogger) GetLogLevel() LogLevel {
	return LogLevel(l.logLevel.Load())
}

// log prints a message if the current log level allows it,
// It adds passed prefix and formats output if a format string is passed.
func (l *ConsoleLogger) log(level LogLevel, prefix string, format string, v ...any) {
	if l.GetLogLevel() > level {
		return
	}

	if format == "" {
		l.stdLogger.Println(prefix, fmt.Sprint(v...))
	} else {
		l.stdLogger.Printf(prefix+" "+format, v...)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *ConsoleLogger) Trace(v ...any) {
	l.log(TRACE, "TRACE:", "", v...)
}

func (l *ConsoleLogger) Tracef(f string, v ...any) {
	l.log(TRACE, "TRACE:", f, v...)
}

func (l *ConsoleLogger) Debug(v ...any) {
	l.log(DEBUG, "DEBUG:", "", v...)
}

func (l *ConsoleLogger) Debugf(f string, v ...any) {
	l.log(DEBUG, "DEBUG:", f, v...)
}

func (l *ConsoleLogger) Info(v ...any) {
	l.log(INFO, "INFO:", "", v...)
}

func (l *ConsoleLogger) Infof(f string, v ...any) {
	l.log(INFO, "INFO:", f, v...)
}

func (l *ConsoleLogger) Warn(v ...any) {
	l.log(WARN, "WARN:", "", v...)
}

func (l *ConsoleLogger) Warnf(f string, v ...any) {
	l.log(WARN, "WARN:", f, v...)
}

func (l *ConsoleLogger) Error(v ...any) {
	l.log(ERROR, "ERROR:", "", v...)
}

func (l *ConsoleLogger) Errorf(f string, v ...any) {
	l.log(ERROR, "ERROR:", f, v...)
}

func (l *ConsoleLogger) Fatal(v ...any) {
	l.log(FATAL, "FATAL:", "", v...)
}

func (l *ConsoleLogger) Fatalf(f string, v ...any) {
	l.log(FATAL, "FATAL:", f, v...)
}
