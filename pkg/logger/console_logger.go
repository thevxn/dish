package logger

import (
	"fmt"
	"log"
	"os"
)

// ConsoleLogger logs output to stderr.
type ConsoleLogger struct {
	stdLogger *log.Logger
	logLevel  LogLevel
}

// NewConsoleLogger creates a new ConsoleLogger instance,
// If verbose is true, log level is set to TRACE (otherwise to INFO).
func NewConsoleLogger(verbose bool) *ConsoleLogger {
	l := &ConsoleLogger{
		stdLogger: log.New(os.Stderr, "", log.LstdFlags),
	}

	l.logLevel = INFO
	if verbose {
		l.logLevel = TRACE
	}

	return l
}

// log prints a message if the current log level allows it,
// It adds passed prefix and formats output if a format string is passed.
func (l *ConsoleLogger) log(level LogLevel, prefix string, format string, v ...any) {
	if l.logLevel > level {
		return
	}

	msg := prefix + " " + fmt.Sprint(v...)
	if format != "" {
		msg = prefix + " " + fmt.Sprintf(format, v...)
	}

	l.stdLogger.Print(msg)

	if level == PANIC {
		panic(msg)
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

func (l *ConsoleLogger) Panic(v ...any) {
	l.log(PANIC, "PANIC:", "", v...)
}

func (l *ConsoleLogger) Panicf(f string, v ...any) {
	l.log(PANIC, "PANIC:", f, v...)
}
