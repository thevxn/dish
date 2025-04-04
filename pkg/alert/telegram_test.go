package alert

import (
	"reflect"
	"testing"
)

func TestNewTelegramSender(t *testing.T) {
	mockHTTPClient := &successStatusHTTPClient{}

	chatID := "-123"
	token := "abc123"
	verbose := false
	notifySuccess := false

	expected := &telegramSender{
		httpClient:    mockHTTPClient,
		chatID:        chatID,
		token:         token,
		verbose:       verbose,
		notifySuccess: notifySuccess,
	}

	actual := NewTelegramSender(mockHTTPClient, chatID, token, verbose, notifySuccess)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestSend(t *testing.T) {
	tests := []struct {
		name          string
		client        HTTPClient
		rawMessage    string
		failedCount   int
		notifySuccess bool
		verbose       bool
		wantErr       bool
	}{
		{
			name:          "Failed Sockets",
			client:        &successStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Failed Sockets - Verbose",
			client:        &successStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets with notifySuccess",
			client:        &successStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   0,
			notifySuccess: true,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets without notifySuccess",
			client:        &successStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   0,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets without notifySuccess - Verbose",
			client:        &successStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   0,
			notifySuccess: false,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "Network Error When Sending Telegram Message",
			client:        &failureHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Unexpected Response Code From Telegram",
			client:        &errorStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Error Reading Response Body From Telegram",
			client:        &invalidResponseBodyHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       true,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := NewTelegramSender(tt.client, "-123", "abc123", tt.verbose, tt.notifySuccess)

			err := sender.send(tt.rawMessage, tt.failedCount)

			if tt.wantErr != (err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
