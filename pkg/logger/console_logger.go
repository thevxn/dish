package logger

import (
	"fmt"
	"log"
	"os"
)

// consoleLogger logs output to stderr.
type consoleLogger struct {
	stdLogger *log.Logger
	logLevel  LogLevel
}

// NewConsoleLogger creates a new ConsoleLogger instance,
// If verbose is true, log level is set to TRACE (otherwise to INFO).
func NewConsoleLogger(verbose bool) *consoleLogger {
	l := &consoleLogger{
		stdLogger: log.New(os.Stderr, "", log.LstdFlags),
	}

	l.logLevel = INFO
	if verbose {
		l.logLevel = TRACE
	}

	return l
}

// log prints a message if the current log level allows it.
// It adds the passed prefix and formats the output if a format string is passed.
func (l *consoleLogger) log(level LogLevel, prefix string, format string, v ...any) {
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

func (l *consoleLogger) Trace(v ...any) {
	l.log(TRACE, "TRACE:", "", v...)
}

func (l *consoleLogger) Tracef(f string, v ...any) {
	l.log(TRACE, "TRACE:", f, v...)
}

func (l *consoleLogger) Debug(v ...any) {
	l.log(DEBUG, "DEBUG:", "", v...)
}

func (l *consoleLogger) Debugf(f string, v ...any) {
	l.log(DEBUG, "DEBUG:", f, v...)
}

func (l *consoleLogger) Info(v ...any) {
	l.log(INFO, "INFO:", "", v...)
}

func (l *consoleLogger) Infof(f string, v ...any) {
	l.log(INFO, "INFO:", f, v...)
}

func (l *consoleLogger) Warn(v ...any) {
	l.log(WARN, "WARN:", "", v...)
}

func (l *consoleLogger) Warnf(f string, v ...any) {
	l.log(WARN, "WARN:", f, v...)
}

func (l *consoleLogger) Error(v ...any) {
	l.log(ERROR, "ERROR:", "", v...)
}

func (l *consoleLogger) Errorf(f string, v ...any) {
	l.log(ERROR, "ERROR:", f, v...)
}

func (l *consoleLogger) Panic(v ...any) {
	l.log(PANIC, "PANIC:", "", v...)
}

func (l *consoleLogger) Panicf(f string, v ...any) {
	l.log(PANIC, "PANIC:", f, v...)
}
