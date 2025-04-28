package logger

import (
	"bytes"
	"log"
	"testing"
)

func TestNewConsoleLogger(t *testing.T) {
	t.Run("verbose mode on", func(t *testing.T) {
		logger := NewConsoleLogger(true)
		if logger.logLevel != TRACE {
			t.Errorf("expected loglevel to be %d, got %d", TRACE, logger.logLevel)
		}
	})

	t.Run("verbose mode off", func(t *testing.T) {
		logger := NewConsoleLogger(false)
		if logger.logLevel != INFO {
			t.Errorf("expected loglevel to be %d, got %d", INFO, logger.logLevel)
		}
	})
}

func TestConsoleLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := &ConsoleLogger{
		stdLogger: log.New(&buf, "", 0),
		logLevel:  TRACE,
	}

	tests := []struct {
		name     string
		logFunc  func()
		expected string
	}{
		{
			name: "Info",
			logFunc: func() {
				logger.Info("hello")
			},
			expected: "INFO: hello\n",
		},
		{
			name: "Infof",
			logFunc: func() {
				logger.Infof("hello %s", "dish")
			},
			expected: "INFO: hello dish\n",
		},
	}

	for _, tt := range tests {
		buf.Reset()

		tt.logFunc()

		output := buf.String()

		if output != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, output)
		}
	}
}

func TestConsoleLogger_Panic(t *testing.T) {
	logger := NewConsoleLogger(true)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic but did not get one")
		}

		expected := "fatal error: dish"
		if r != expected {
			t.Fatalf("expected panic message %s, got %s", expected, r)
		}
	}()

	logger.Fatalf("fatal error: %s", "dish")
}
