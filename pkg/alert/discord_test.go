package alert

import (
	"reflect"
	"testing"

	"go.vxn.dev/dish/pkg/config"
)

func TestNewDiscordSender(t *testing.T) {
	mockHTTPClient := &SuccessStatusHTTPClient{}
	mockLogger := &MockLogger{}

	testBotToken := "test"
	testChannelID := "-123"

	expected := &discordSender{
		botToken:      testBotToken,
		channelID:     testChannelID,
		httpClient:    mockHTTPClient,
		logger:        mockLogger,
		notifySuccess: false,
	}

	cfg := &config.Config{
		DiscordBotToken:   testBotToken,
		DiscordChannelID:  testChannelID,
		TextNotifySuccess: false,
	}

	actual := NewDiscordSender(mockHTTPClient, cfg, mockLogger)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestSend_Discord(t *testing.T) {
	newConfig := func(botToken, channelID string, notifySuccess bool) *config.Config {
		return &config.Config{
			DiscordBotToken:   botToken,
			DiscordChannelID:  channelID,
			TextNotifySuccess: notifySuccess,
		}
	}

	tests := []struct {
		name          string
		client        HTTPClient
		rawMessage    string
		failedCount   int
		notifySuccess bool
		wantErr       bool
	}{
		{
			name:          "success with failed checks",
			client:        &SuccessStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			wantErr:       false,
		},
		{
			name:          "success with failed checks",
			client:        &SuccessStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: true,
			wantErr:       false,
		},
		{
			name:          "success with failed checks",
			client:        &SuccessStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   0,
			notifySuccess: false,
			wantErr:       false,
		},
		{
			name:          "success with no failed checks and notify success enabled",
			client:        &SuccessStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   0,
			notifySuccess: true,
			wantErr:       false,
		},
		{
			name:          "Network Error When Sending Discord Message",
			client:        &FailureHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			wantErr:       true,
		},
		{
			name:          "Network Error When Sending Discord Message",
			client:        &FailureHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			wantErr:       true,
		},
		{
			name:          "Unexpected Response Code From Discord",
			client:        &ErrorStatusHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			wantErr:       true,
		},
		{
			name:          "Error Reading Response Body From Discord",
			client:        &InvalidResponseBodyHTTPClient{},
			rawMessage:    "Test message",
			failedCount:   1,
			notifySuccess: false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := NewDiscordSender(tt.client, newConfig("test", "-123", tt.notifySuccess), &MockLogger{})
			err := sender.send(tt.rawMessage, tt.failedCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
