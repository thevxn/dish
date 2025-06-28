package alert

import (
	"flag"
	"testing"

	"go.vxn.dev/dish/pkg/config"
)

const (
	badURL = "0řžuničx://č/tmpě/test.hook\x09"
)

var (
	configDefault *config.Config
)

func TestNewNotifier(t *testing.T) {
	var (
		configBlank             = &config.Config{}
		mockLogger              = &MockLogger{}
		successStatushTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ = config.NewConfig(flag.CommandLine, []string{""})

	if notifier := NewNotifier(nil, nil, nil); notifier != nil {
		t.Error("unexpected behaviour, should be nil")
	}

	if notifier := NewNotifier(nil, configBlank, nil); notifier != nil {
		t.Error("unexpected behaviour, should be nil")
	}

	if notifier := NewNotifier(&successStatushTTPClient, configBlank, nil); notifier != nil {
		t.Error("expected nil, got notifier (nil logger)")
	}

	if notifier := NewNotifier(&successStatushTTPClient, nil, mockLogger); notifier != nil {
		t.Error("expected nil, got notifier (nil config)")
	}

	if notifier := NewNotifier(&successStatushTTPClient, configDefault, mockLogger); notifier == nil {
		t.Error("unexpected nil on output")
	}

	// Telegram tests

	configDefault.TelegramBotToken = "abc:2025062700"
	configDefault.TelegramChatID = "-10987654321"

	notifierTelegram := NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if notifierTelegram == nil {
		t.Fatal("unexpected nil on output (Telegram*)")
	}

	if notifiersLen := len(notifierTelegram.chatNotifiers); notifiersLen == 0 {
		t.Errorf("expected 1 chatNotifier: got %d", notifiersLen)
	}

	// Remote API tests

	configDefault.ApiURL = "https://api.example.com/?test=true"

	notifierAPI := NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if notifierAPI == nil {
		t.Fatal("unexpected nil on output (ApiURL)")
	}

	if len(notifierAPI.machineNotifiers) != 1 {
		t.Errorf("expected 1 machineNotifier, got %d", len(notifierAPI.machineNotifiers))
	}

	configDefault.ApiURL = badURL

	notifierAPI = NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if notifierAPI == nil {
		t.Fatal("unexpected nil on output (ApiURL)")
	}

	if len(notifierAPI.machineNotifiers) != 0 {
		t.Errorf("expected 0 machineNotifiers, got %d", len(notifierAPI.machineNotifiers))
	}

	// Webhooks tests

	configDefault.WebhookURL = "https://www.example.com/hooks/test-hook"

	notifierWebhook := NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if notifierWebhook == nil {
		t.Fatal("unexpected nil on output (Webhooks)")
	}

	if len(notifierWebhook.machineNotifiers) != 1 {
		t.Errorf("expected 1 machineNotifier, got %d", len(notifierWebhook.machineNotifiers))
	}

	configDefault.WebhookURL = badURL

	notifierWebhook = NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if notifierWebhook == nil {
		t.Fatal("unexpected nil on output (Webhooks)")
	}

	if len(notifierWebhook.machineNotifiers) != 0 {
		t.Errorf("expected 0 machineNotifiers, got %d", len(notifierWebhook.machineNotifiers))
	}

	// Pushgateway tests

	configDefault.PushgatewayURL = "https://pgw.example.com/push/"

	notifierPushgateway := NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if notifierPushgateway == nil {
		t.Fatal("unexpected nil on output (Pushgateway)")
	}

	if len(notifierPushgateway.machineNotifiers) != 1 {
		t.Errorf("expected 1 machineNotifier, got %d", len(notifierPushgateway.machineNotifiers))
	}

	configDefault.PushgatewayURL = badURL

	notifierPushgateway = NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if notifierPushgateway == nil {
		t.Fatal("unexpected nil on output (Pushgateway)")
	}

	if len(notifierPushgateway.machineNotifiers) != 0 {
		t.Errorf("expected 0 machineNotifiers, got %d", len(notifierPushgateway.machineNotifiers))
	}
}

func TestSendChatNotifications(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatushTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault.TelegramBotToken = ""
	configDefault.TelegramChatID = ""

	notifierTelegram := NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if len(notifierTelegram.chatNotifiers) > 0 {
		t.Errorf("expected 0 chatNotifiers, got %d", len(notifierTelegram.chatNotifiers))
	}

	if err := notifierTelegram.SendChatNotifications("SendChatNotifications test", 0); err != nil {
		t.Error("unexpected error: ", err)
	}

	configDefault.TelegramBotToken = "abc:2025062700"
	configDefault.TelegramChatID = "-10987654321"

	notifierTelegram = NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if notifierTelegram == nil {
		t.Fatal("unexpected nil on output")
	}

	if err := notifierTelegram.SendChatNotifications("SendChatNotifications test", 0); err != nil {
		t.Error("unexpected error: ", err)
	}

	mockTelegram := telegramSender{httpClient: &successStatushTTPClient, logger: mockLogger, token: "$á+\x00"}
	notifierTelegram.chatNotifiers[0] = &mockTelegram

	if err := notifierTelegram.SendChatNotifications("", 20); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSendMachineNotifications(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatushTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault.WebhookURL = ""

	notifierWebhook := NewNotifier(&successStatushTTPClient, configDefault, mockLogger)
	if notifierWebhook == nil {
		t.Error("unexpected nil on output (Webhooks)")
	}

	if err := notifierWebhook.SendMachineNotifications(nil, 0); err != nil {
		t.Error("unexpected error: ", err)
	}

	configDefault.WebhookURL = "https://www.example.com/hooks/test-hook"

	notifierWebhook = NewNotifier(nil, configDefault, mockLogger)
	if notifierWebhook == nil {
		t.Fatal("unexpected nil on output (Webhooks)")
	}

	if err := notifierWebhook.SendMachineNotifications(nil, 0); err != nil {
		t.Error("unexpected error: ", err)
	}

	mockWebhook := webhookSender{httpClient: nil, url: badURL, logger: mockLogger}
	notifierWebhook.machineNotifiers[0] = &mockWebhook

	if err := notifierWebhook.SendMachineNotifications(nil, 20); err == nil {
		t.Error("expected error, got nil")
	}
}
