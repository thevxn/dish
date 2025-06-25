package alert

import (
	"reflect"
	"testing"

	"go.vxn.dev/dish/pkg/config"
	testhelpers "go.vxn.dev/dish/pkg/testdata"
)

func TestNewTelegramSender(t *testing.T) {
	mockHTTPClient := &testhelpers.SuccessStatusHTTPClient{}
	mockLogger := &testhelpers.MockLogger{}

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
		logger:        mockLogger,
	}

	cfg := &config.Config{
		TelegramChatID:    chatID,
		TelegramBotToken:  token,
		Verbose:           verbose,
		TextNotifySuccess: notifySuccess,
	}

	actual := NewTelegramSender(mockHTTPClient, cfg, mockLogger)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestSend_Telegram(t *testing.T) {
	newConfig := func(chatID, token string, verbose, notifySuccess bool) *config.Config {
		return &config.Config{
			TelegramChatID:    chatID,
			TelegramBotToken:  token,
			Verbose:           verbose,
			TextNotifySuccess: notifySuccess,
		}
	}

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
			client:        &testhelpers.SuccessStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Failed Sockets - Verbose",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets with notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   0,
			notifySuccess: true,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets without notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   0,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets without notifySuccess - Verbose",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   0,
			notifySuccess: false,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "Network Error When Sending Telegram Message",
			client:        &testhelpers.FailureHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Unexpected Response Code From Telegram",
			client:        &testhelpers.ErrorStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Error Reading Response Body From Telegram",
			client:        &testhelpers.InvalidResponseBodyHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			verbose:       true,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig("-123", "abc123", tt.verbose, tt.notifySuccess)
			sender := NewTelegramSender(tt.client, cfg, &testhelpers.MockLogger{})

			err := sender.send(tt.rawMessage, tt.failedCount)

			if tt.wantErr != (err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
