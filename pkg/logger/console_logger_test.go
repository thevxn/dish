package logger

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestNewConsoleLogger(t *testing.T) {
	t.Run("verbose mode on", func(t *testing.T) {
		logger := NewConsoleLogger(true, nil)
		if logger.logLevel != TRACE {
			t.Errorf("expected loglevel %d, got %d", TRACE, logger.logLevel)
		}
	})

	t.Run("verbose mode off", func(t *testing.T) {
		logger := NewConsoleLogger(false, nil)
		if logger.logLevel != INFO {
			t.Errorf("expected loglevel %d, got %d", INFO, logger.logLevel)
		}
	})

	t.Run("default out", func(t *testing.T) {
		origStderr := os.Stderr
		defer func() { os.Stderr = origStderr }()

		r, w, _ := os.Pipe()
		os.Stderr = w

		logger := &consoleLogger{
			stdLogger: log.New(w, "", 0),
			logLevel:  TRACE,
		}
		logger.Info("hello stderr")

		w.Close()
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r)
		if err != nil {
			t.Fatalf("failed to read from pipe: %v", err)
		}

		expected := INFO.Prefix(false) + "hello stderr\n"
		actual := buf.String()
		if actual != expected {
			t.Fatalf("expected %q in stderr, got %q", expected, actual)
		}
	})

	t.Run("provided out", func(t *testing.T) {
		var buf bytes.Buffer

		logger := &consoleLogger{
			stdLogger: log.New(&buf, "", 0),
			logLevel:  TRACE,
		}

		logger.Info("output test")

		expected := INFO.Prefix(false) + "output test\n"
		actual := buf.String()

		if actual != expected {
			t.Fatalf("expected %s, got %s", expected, actual)
		}

	})

	t.Run("with colors when verbose and no env set", func(t *testing.T) {
		logger := NewConsoleLogger(true, nil)
		if !logger.withColors {
			t.Error("expected logger to have colors enabled")
		}
	})

	t.Run("without colors when env is set", func(t *testing.T) {
		_ = os.Setenv("NO_COLOR", "true")
		logger := NewConsoleLogger(true, nil)
		if logger.withColors {
			t.Error("expected logger to have colors disabled")
		}
		_ = os.Setenv("NO_COLOR", "false")
	})

	t.Run("without colors when not verbose is not set", func(t *testing.T) {
		logger := NewConsoleLogger(false, nil)
		if logger.withColors {
			t.Error("expected logger to have colors disabled")
		}
	})
}

func TestConsoleLogger_log(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		name     string
		logFunc  func(*consoleLogger)
		logger   *consoleLogger
		expected string
	}{
		{
			name: "Info adds INFO prefix and joins arguments with spaces",
			logFunc: func(logger *consoleLogger) {
				logger.Info("hello", 123, 321)
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  TRACE,
			},
			expected: INFO.Prefix(false) + "hello123 321\n",
		},
		{
			name: "Info adds INFO prefix with colors",
			logFunc: func(logger *consoleLogger) {
				logger.Info("hello")
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   TRACE,
				withColors: true,
			},
			expected: INFO.Prefix(true) + "hello\n",
		},
		{
			name: "Infof adds INFO prefix and formats string correctly",
			logFunc: func(logger *consoleLogger) {
				logger.Infof("hello %s !", "dish")
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  TRACE,
			},
			expected: INFO.Prefix(false) + "hello dish !\n",
		},
		{
			name: "Infof adds INFO prefix with color",
			logFunc: func(logger *consoleLogger) {
				logger.Infof("hello %s !", "dish")
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   TRACE,
				withColors: true,
			},
			expected: INFO.Prefix(true) + "hello dish !\n",
		},
		{
			name: "Debug does not print if logLevel is INFO",
			logFunc: func(logger *consoleLogger) {
				logger.Debug("should not print")
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  INFO,
			},
			expected: "",
		},
		{
			name: "Debug adds DEBUG prefix",
			logFunc: func(logger *consoleLogger) {
				logger.Debug("debug")
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  DEBUG,
			},
			expected: DEBUG.Prefix(false) + "debug\n",
		},
		{
			name: "Debug adds DEBUG prefix with color",
			logFunc: func(logger *consoleLogger) {
				logger.Debug("debug")
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   DEBUG,
				withColors: true,
			},
			expected: DEBUG.Prefix(true) + "debug\n",
		},
		{
			name: "Debugf adds DEBUG prefix and formats string correctly",
			logFunc: func(logger *consoleLogger) {
				logger.Debugf("debug %d", 1)
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  DEBUG,
			},
			expected: DEBUG.Prefix(false) + "debug 1\n",
		},
		{
			name: "Debugf adds DEBUG prefix with color and formats string correctly",
			logFunc: func(logger *consoleLogger) {
				logger.Debugf("debug %d", 1)
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   DEBUG,
				withColors: true,
			},
			expected: DEBUG.Prefix(true) + "debug 1\n",
		},
		{
			name: "Warn prints with WARN prefix",
			logFunc: func(logger *consoleLogger) {
				logger.Warn("warn message")
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  TRACE,
			},
			expected: WARN.Prefix(false) + "warn message\n",
		},
		{
			name: "Warn prints with WARN prefix with color",
			logFunc: func(logger *consoleLogger) {
				logger.Warn("warn message")
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   TRACE,
				withColors: true,
			},
			expected: WARN.Prefix(true) + "warn message\n",
		},
		{
			name: "Warnf prints formatted WARN message",
			logFunc: func(logger *consoleLogger) {
				logger.Warnf("warn %d", 42)
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  TRACE,
			},
			expected: WARN.Prefix(false) + "warn 42\n",
		},
		{
			name: "Warnf prints formatted WARN message with color",
			logFunc: func(logger *consoleLogger) {
				logger.Warnf("warn %d", 42)
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   TRACE,
				withColors: true,
			},
			expected: WARN.Prefix(true) + "warn 42\n",
		},
		{
			name: "Error prints with ERROR prefix",
			logFunc: func(logger *consoleLogger) {
				logger.Error("error")
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  TRACE,
			},
			expected: ERROR.Prefix(false) + "error\n",
		},
		{
			name: "Error prints with ERROR prefix with color",
			logFunc: func(logger *consoleLogger) {
				logger.Error("error")
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   TRACE,
				withColors: true,
			},
			expected: ERROR.Prefix(true) + "error\n",
		},
		{
			name: "Errorf prints formatted ERROR message",
			logFunc: func(logger *consoleLogger) {
				logger.Errorf("fail %s", "here")
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  TRACE,
			},
			expected: ERROR.Prefix(false) + "fail here\n",
		},
		{
			name: "Errorf prints formatted ERROR message with color",
			logFunc: func(logger *consoleLogger) {
				logger.Errorf("fail %s", "here")
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   TRACE,
				withColors: true,
			},
			expected: ERROR.Prefix(true) + "fail here\n",
		},
		{
			name: "Trace prints with TRACE prefix",
			logFunc: func(logger *consoleLogger) {
				logger.Trace("trace")
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  TRACE,
			},
			expected: TRACE.Prefix(false) + "trace\n",
		},
		{
			name: "Trace prints with TRACE prefix with color",
			logFunc: func(logger *consoleLogger) {
				logger.Trace("trace")
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   TRACE,
				withColors: true,
			},
			expected: TRACE.Prefix(true) + "trace\n",
		},
		{
			name: "Tracef prints formatted TRACE message",
			logFunc: func(logger *consoleLogger) {
				logger.Tracef("trace %d", 1)
			},
			logger: &consoleLogger{
				stdLogger: log.New(&buf, "", 0),
				logLevel:  TRACE,
			},
			expected: TRACE.Prefix(false) + "trace 1\n",
		},
		{
			name: "Tracef prints formatted TRACE message with color",
			logFunc: func(logger *consoleLogger) {
				logger.Tracef("trace %d", 1)
			},
			logger: &consoleLogger{
				stdLogger:  log.New(&buf, "", 0),
				logLevel:   TRACE,
				withColors: true,
			},
			expected: TRACE.Prefix(true) + "trace 1\n",
		},
	}

	for _, tt := range tests {
		buf.Reset()

		tt.logFunc(tt.logger)

		output := buf.String()

		if output != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, output)
		}
	}
}

func TestConsoleLogger_log_Panic(t *testing.T) {
	logger := NewConsoleLogger(true, nil)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic but did not get one")
		}

		expected := PANIC.Prefix(true) + "could not start dish"
		if r != expected {
			t.Fatalf("expected panic message %s, got %s", expected, r)
		}
	}()

	logger.Panic("could not start dish")
}

func TestConsoleLogger_log_Panicf(t *testing.T) {
	logger := NewConsoleLogger(true, nil)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic but did not get one")
		}

		expected := PANIC.Prefix(true) + "could not start dish"
		if r != expected {
			t.Fatalf("expected panic message %s, got %s", expected, r)
		}
	}()

	logger.Panicf("could not start %s", "dish")
}
