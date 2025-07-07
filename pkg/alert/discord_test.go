package alert

import (
	"reflect"
	"testing"

	"go.vxn.dev/dish/pkg/config"
)

func TestNewDiscordSender(t *testing.T) {
	mockHTTPClient := &SuccessStatusHTTPClient{}
	mockLogger := &MockLogger{}

	tests := []struct {
		name           string
		botToken       string
		channelID      string
		notifySuccess  bool
		expectErr      bool
		expectedSender *discordSender
	}{
		{
			name:          "successful sender creation",
			botToken:      "test",
			channelID:     "-123",
			notifySuccess: false,
			expectErr:     false,
			expectedSender: &discordSender{
				botToken:      "test",
				channelID:     "-123",
				httpClient:    mockHTTPClient,
				logger:        mockLogger,
				notifySuccess: false,
				url:           "https://discord.com/api/v10/channels/-123/messages",
			},
		},
		{
			name:           "invalid channel ID returns error",
			botToken:       "test",
			channelID:      "1%ZZ",
			notifySuccess:  false,
			expectErr:      true,
			expectedSender: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DiscordBotToken:   tt.botToken,
				DiscordChannelID:  tt.channelID,
				TextNotifySuccess: tt.notifySuccess,
			}

			actual, err := NewDiscordSender(mockHTTPClient, cfg, mockLogger)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tt.expectedSender, actual) {
				t.Errorf("expected %+v, got %+v", tt.expectedSender, actual)
			}
		})
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
			sender, _ := NewDiscordSender(tt.client, newConfig("test", "-123", tt.notifySuccess), &MockLogger{})
			err := sender.send(tt.rawMessage, tt.failedCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
