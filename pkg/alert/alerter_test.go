package alert

import (
	"testing"

	"go.vxn.dev/dish/pkg/config"
)

func TestNewAlerter(t *testing.T) {
	var (
		mockLogger = MockLogger{}
	)

	if alerterNil := NewAlerter(nil); alerterNil != nil {
		t.Error("expected nil, got alerter")
	}

	if alerter := NewAlerter(&mockLogger); alerter == nil {
		t.Error("expected alerter, got nil")
	}
}

func TestHandleAlerts(t *testing.T) {
	var (
		mockConfig  = config.Config{}
		mockLogger  = MockLogger{}
		mockResults = Results{}
	)

	alerter := NewAlerter(&mockLogger)
	if alerter == nil {
		t.Error("expected alerter, got nil")
	}

	alerter.HandleAlerts("", nil, 0, nil)

	alerter.HandleAlerts("HandleAlerts test", &mockResults, 20, &mockConfig)
}
