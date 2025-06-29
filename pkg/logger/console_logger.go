package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

// consoleLogger logs to the output provided when instantiating it via NewConsoleLogger.
type consoleLogger struct {
	stdLogger  *log.Logger
	logLevel   logLevel
	withColors bool
}

var defaultOut = os.Stderr

// NewConsoleLogger creates a new ConsoleLogger instance logging to the provided output.
// If the output is not specified (nil), it logs to stderr by default.
//
// If verbose is true, log level is set to TRACE (otherwise to INFO).
func NewConsoleLogger(verbose bool, out io.Writer) *consoleLogger {
	if out == nil {
		out = defaultOut
	}

	l := &consoleLogger{
		stdLogger:  log.New(out, "", log.LstdFlags),
		withColors: !(os.Getenv("NO_COLOR") == "true") && verbose,
	}

	l.logLevel = INFO
	if verbose {
		l.logLevel = TRACE
	}

	return l
}

// log prints a message if the current log level allows it.
// It adds the passed prefix and formats the output if a format string is passed.
func (l *consoleLogger) log(level logLevel, prefix string, format string, v ...any) {
	if l.logLevel > level {
		return
	}

	msg := prefix + fmt.Sprint(v...)
	if format != "" {
		msg = prefix + fmt.Sprintf(format, v...)
	}

	l.stdLogger.Print(msg)

	if level == PANIC {
		panic(msg)
	}
}

func (l *consoleLogger) Trace(v ...any) {
	l.log(TRACE, TRACE.Prefix(l.withColors), "", v...)
}

func (l *consoleLogger) Tracef(f string, v ...any) {
	l.log(TRACE, TRACE.Prefix(l.withColors), f, v...)
}

func (l *consoleLogger) Debug(v ...any) {
	l.log(DEBUG, DEBUG.Prefix(l.withColors), "", v...)
}

func (l *consoleLogger) Debugf(f string, v ...any) {
	l.log(DEBUG, DEBUG.Prefix(l.withColors), f, v...)
}

func (l *consoleLogger) Info(v ...any) {
	l.log(INFO, INFO.Prefix(l.withColors), "", v...)
}

func (l *consoleLogger) Infof(f string, v ...any) {
	l.log(INFO, INFO.Prefix(l.withColors), f, v...)
}

func (l *consoleLogger) Warn(v ...any) {
	l.log(WARN, WARN.Prefix(l.withColors), "", v...)
}

func (l *consoleLogger) Warnf(f string, v ...any) {
	l.log(WARN, WARN.Prefix(l.withColors), f, v...)
}

func (l *consoleLogger) Error(v ...any) {
	l.log(ERROR, ERROR.Prefix(l.withColors), "", v...)
}

func (l *consoleLogger) Errorf(f string, v ...any) {
	l.log(ERROR, ERROR.Prefix(l.withColors), f, v...)
}

func (l *consoleLogger) Panic(v ...any) {
	l.log(PANIC, PANIC.Prefix(l.withColors), "", v...)
}

func (l *consoleLogger) Panicf(f string, v ...any) {
	l.log(PANIC, PANIC.Prefix(l.withColors), f, v...)
}
