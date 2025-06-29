package logger

import (
	"fmt"
	"testing"
)

func TestLogLevel_Color(t *testing.T) {
	tests := []struct {
		level    logLevel
		expected string
	}{
		{TRACE, "\033[34m"},
		{DEBUG, "\033[36m"},
		{INFO, "\033[32m"},
		{WARN, "\033[33m"},
		{ERROR, "\033[31m"},
		{PANIC, "\033[35m"},
		{logLevel(999), "\033[0m"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%v", tt.level), func(t *testing.T) {
			got := tt.level.Color()
			if got != tt.expected {
				t.Errorf("color is: %q, but should be: %q", got, tt.expected)
			}
		})
	}
}

func TestLogLevel_Prefix(t *testing.T) {
	tests := []struct {
		level      logLevel
		withColor  bool
		expectPart string
	}{
		{TRACE, false, "[ TRACE ]: "},
		{DEBUG, false, "[ DEBUG ]: "},
		{INFO, false, "[ INFO ]: "},
		{WARN, false, "[ WARN ]: "},
		{ERROR, false, "[ ERROR ]: "},
		{PANIC, false, "[ PANIC ]: "},
		{logLevel(999), false, "[ UNKNOWN ]: "},
		{TRACE, true, "\033[34m[ TRACE ]\033[0m: "},
		{logLevel(999), true, "[ UNKNOWN ]: "},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%v_Color_%v", tt.level, tt.withColor), func(t *testing.T) {
			got := tt.level.Prefix(tt.withColor)
			if tt.expectPart != got {
				t.Errorf("prefix is: %q, but should be: %q", got, tt.expectPart)
			}
		})
	}
}
