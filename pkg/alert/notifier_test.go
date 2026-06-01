package alert

import (
	"flag"
	"testing"

	"go.vxn.dev/dish/pkg/config"
)

const (
	badURL = "0řžuničx://č/tmpě/test.hook\x09"
)

func TestNewNotifier_Nil(t *testing.T) {
	var (
		configBlank             = &config.Config{}
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ := config.NewConfig(flag.NewFlagSet("test", flag.ContinueOnError), []string{""})

	_, err := NewNotifier(nil, nil, nil)
	if err == nil {
		t.Error("expected error, got nil")
	}

	_, err = NewNotifier(nil, configBlank, nil)
	if err == nil {
		t.Error("expected error, got nil")
	}

	_, err = NewNotifier(&successStatusHTTPClient, configBlank, nil)
	if err == nil {
		t.Error("expected error, got nil")
	}

	_, err = NewNotifier(&successStatusHTTPClient, nil, mockLogger)
	if err == nil {
		t.Error("expected error, got nil")
	}

	_, err = NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

}

func TestNewNotifier_Telegram(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ := config.NewConfig(flag.NewFlagSet("test", flag.ContinueOnError), []string{""})
	configDefault.TelegramBotToken = "abc:2025062700"
	configDefault.TelegramChatID = "-10987654321"

	notifierTelegram, err := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if notifiersLen := len(notifierTelegram.chatNotifiers); notifiersLen == 0 {
		t.Errorf("expected 1 chatNotifier: got %d", notifiersLen)
	}
}

func TestNewNotifier_API(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ := config.NewConfig(flag.NewFlagSet("test", flag.ContinueOnError), []string{""})
	configDefault.ApiURL = "https://api.example.com/?test=true"

	notifierAPI, err := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if len(notifierAPI.machineNotifiers) != 1 {
		t.Errorf("expected 1 machineNotifier, got %d", len(notifierAPI.machineNotifiers))
	}

	configDefault.ApiURL = badURL

	notifierAPI, err = NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if len(notifierAPI.machineNotifiers) != 0 {
		t.Errorf("expected 0 machineNotifiers, got %d", len(notifierAPI.machineNotifiers))
	}
}

func TestNewNotifier_Webhook(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ := config.NewConfig(flag.NewFlagSet("test", flag.ContinueOnError), []string{""})
	configDefault.WebhookURL = "https://www.example.com/hooks/test-hook"

	notifierWebhook, err := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if len(notifierWebhook.machineNotifiers) != 1 {
		t.Errorf("expected 1 machineNotifier, got %d", len(notifierWebhook.machineNotifiers))
	}

	configDefault.WebhookURL = badURL

	notifierWebhook, err = NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if len(notifierWebhook.machineNotifiers) != 0 {
		t.Errorf("expected 0 machineNotifiers, got %d", len(notifierWebhook.machineNotifiers))
	}
}

func TestNewNotifier_Pushgateway(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ := config.NewConfig(flag.NewFlagSet("test", flag.ContinueOnError), []string{""})
	configDefault.PushgatewayURL = "https://pgw.example.com/push/"

	notifierPushgateway, err := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if len(notifierPushgateway.machineNotifiers) != 1 {
		t.Errorf("expected 1 machineNotifier, got %d", len(notifierPushgateway.machineNotifiers))
	}

	configDefault.PushgatewayURL = badURL

	notifierPushgateway, err = NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if len(notifierPushgateway.machineNotifiers) != 0 {
		t.Errorf("expected 0 machineNotifiers, got %d", len(notifierPushgateway.machineNotifiers))
	}
}

func TestNewNotifier_Discord(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ := config.NewConfig(flag.NewFlagSet("test", flag.ContinueOnError), []string{""})
	configDefault.DiscordBotToken = "test"
	configDefault.DiscordChannelID = "-123"

	notifierDiscord, err := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if notifiersLen := len(notifierDiscord.chatNotifiers); notifiersLen == 0 {
		t.Errorf("expected 1 chatNotifier: got %d", notifiersLen)
	}
}

func TestSendChatNotifications(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ := config.NewConfig(flag.NewFlagSet("test", flag.ContinueOnError), []string{""})
	configDefault.TelegramBotToken = ""
	configDefault.TelegramChatID = ""

	notifierTelegram, err := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if len(notifierTelegram.chatNotifiers) > 0 {
		t.Errorf("expected 0 chatNotifiers, got %d", len(notifierTelegram.chatNotifiers))
	}

	if err := notifierTelegram.SendChatNotifications("SendChatNotifications test", 0); err != nil {
		t.Error("unexpected error: ", err)
	}

	configDefault.TelegramBotToken = "abc:2025062700"
	configDefault.TelegramChatID = "-10987654321"

	notifierTelegram, err = NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if err := notifierTelegram.SendChatNotifications("SendChatNotifications test", 0); err != nil {
		t.Error("unexpected error: ", err)
	}

	mockTelegram := telegramSender{httpClient: &successStatusHTTPClient, logger: mockLogger, token: "$á+\x00"}
	notifierTelegram.chatNotifiers[0] = &mockTelegram

	if err := notifierTelegram.SendChatNotifications("", 20); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSendMachineNotifications(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)
	configDefault, _ := config.NewConfig(flag.NewFlagSet("test", flag.ContinueOnError), []string{""})
	configDefault.WebhookURL = ""

	notifierWebhook, err := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}

	if err := notifierWebhook.SendMachineNotifications(nil, 0); err != nil {
		t.Error("unexpected error: ", err)
	}

	configDefault.WebhookURL = "https://www.example.com/hooks/test-hook"

	notifierWebhook, err = NewNotifier(nil, configDefault, mockLogger)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
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
