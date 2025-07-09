package alert

import (
	"flag"
	"fmt"
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

	configDefault, _ := config.NewConfig(
		flag.NewFlagSet("test", flag.ContinueOnError),
		[]string{""},
	)

	if notifier := NewNotifier(nil, nil, nil); notifier != nil {
		t.Error("unexpected behaviour, should be nil")
	}

	if notifier := NewNotifier(nil, configBlank, nil); notifier != nil {
		t.Error("unexpected behaviour, should be nil")
	}

	if notifier := NewNotifier(&successStatusHTTPClient, configBlank, nil); notifier != nil {
		t.Error("expected nil, got notifier (nil logger)")
	}

	if notifier := NewNotifier(&successStatusHTTPClient, nil, mockLogger); notifier != nil {
		t.Error("expected nil, got notifier (nil config)")
	}

	if notifier := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger); notifier == nil {
		t.Error("unexpected nil on output")
	}
}

func TestNewNotifier_Telegram(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ := config.NewConfig(
		flag.NewFlagSet("test", flag.ContinueOnError),
		[]string{""},
	)
	configDefault.TelegramBotToken = "abc:2025062700"
	configDefault.TelegramChatID = "-10987654321"

	notifierTelegram := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierTelegram == nil {
		t.Fatal("unexpected nil on output (Telegram*)")
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

	configDefault, _ := config.NewConfig(
		flag.NewFlagSet("test", flag.ContinueOnError),
		[]string{""},
	)
	configDefault.ApiURL = "https://api.example.com/?test=true"

	notifierAPI := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierAPI == nil {
		t.Fatal("unexpected nil on output (ApiURL)")
	}

	if len(notifierAPI.machineNotifiers) != 1 {
		t.Errorf("expected 1 machineNotifier, got %d", len(notifierAPI.machineNotifiers))
	}

	configDefault.ApiURL = badURL

	notifierAPI = NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierAPI == nil {
		t.Fatal("unexpected nil on output (ApiURL)")
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

	configDefault, _ := config.NewConfig(
		flag.NewFlagSet("test", flag.ContinueOnError),
		[]string{""},
	)
	configDefault.WebhookURL = "https://www.example.com/hooks/test-hook"

	notifierWebhook := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierWebhook == nil {
		t.Fatal("unexpected nil on output (Webhooks)")
	}

	if len(notifierWebhook.machineNotifiers) != 1 {
		t.Errorf("expected 1 machineNotifier, got %d", len(notifierWebhook.machineNotifiers))
	}

	configDefault.WebhookURL = badURL

	notifierWebhook = NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierWebhook == nil {
		t.Fatal("unexpected nil on output (Webhooks)")
	}

	if len(notifierWebhook.machineNotifiers) != 0 {
		t.Errorf("expected 0 machineNotifiers, got %d", len(notifierWebhook.machineNotifiers))
	}
}

func TestNewNotifier_Pushgateawy(t *testing.T) {
	var (
		mockLogger              = &MockLogger{}
		successStatusHTTPClient = SuccessStatusHTTPClient{}
	)

	configDefault, _ := config.NewConfig(
		flag.NewFlagSet("test", flag.ContinueOnError),
		[]string{""},
	)
	configDefault.PushgatewayURL = "https://pgw.example.com/push/"

	notifierPushgateway := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierPushgateway == nil {
		t.Fatal("unexpected nil on output (Pushgateway)")
	}

	if len(notifierPushgateway.machineNotifiers) != 1 {
		t.Errorf("expected 1 machineNotifier, got %d", len(notifierPushgateway.machineNotifiers))
	}

	configDefault.PushgatewayURL = badURL

	notifierPushgateway = NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierPushgateway == nil {
		t.Fatal("unexpected nil on output (Pushgateway)")
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

	configDefault, _ := config.NewConfig(
		flag.NewFlagSet("test", flag.ContinueOnError),
		[]string{""},
	)
	configDefault.DiscordBotToken = "test"
	configDefault.DiscordChannelID = "-123"

	notifierDiscord := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierDiscord == nil {
		t.Fatal("unexpected nil on output (Discord*)")
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

	configDefault, _ := config.NewConfig(
		flag.NewFlagSet("test", flag.ContinueOnError),
		[]string{""},
	)
	configDefault.TelegramBotToken = ""
	configDefault.TelegramChatID = ""

	notifierTelegram := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	fmt.Println(notifierTelegram.chatNotifiers)
	if len(notifierTelegram.chatNotifiers) > 0 {
		t.Errorf("expected 0 chatNotifiers, got %d", len(notifierTelegram.chatNotifiers))
	}

	if err := notifierTelegram.SendChatNotifications("SendChatNotifications test", 0); err != nil {
		t.Error("unexpected error: ", err)
	}

	configDefault.TelegramBotToken = "abc:2025062700"
	configDefault.TelegramChatID = "-10987654321"

	notifierTelegram = NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierTelegram == nil {
		t.Fatal("unexpected nil on output")
	}

	if err := notifierTelegram.SendChatNotifications("SendChatNotifications test", 0); err != nil {
		t.Error("unexpected error: ", err)
	}

	mockTelegram := telegramSender{
		httpClient: &successStatusHTTPClient,
		logger:     mockLogger,
		token:      "$á+\x00",
	}
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
	configDefault, _ := config.NewConfig(
		flag.NewFlagSet("test", flag.ContinueOnError),
		[]string{""},
	)
	configDefault.WebhookURL = ""

	notifierWebhook := NewNotifier(&successStatusHTTPClient, configDefault, mockLogger)
	if notifierWebhook == nil {
		t.Fatal("unexpected nil on output (Webhooks)")
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
