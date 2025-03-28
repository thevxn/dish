package alert

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"go.vxn.dev/dish/pkg/config"
)

type Results struct {
	Map map[string]bool `json:"dish_results"`
}

type ChatNotifier interface {
	send(string, int) error
}
type MachineNotifier interface {
	send(Results, int) error
}

type notifier struct {
	verbose          bool
	chatNotifiers    []ChatNotifier
	machineNotifiers []MachineNotifier
}

func NewNotifier(httpClient *http.Client, config *config.Config) *notifier {
	// Set chat integrations to be notified (e.g. Telegram)
	notificationSenders := make([]ChatNotifier, 0)
	if config.TelegramBotToken != "" && config.TelegramChatID != "" {
		notificationSenders = append(notificationSenders, NewTelegramSender(httpClient, config.TelegramChatID, config.TelegramBotToken, config.Verbose, config.FailedOnly))
	}

	// Set machine interface integrations to be notified (e.g. Webhooks)
	payloadSenders := make([]MachineNotifier, 0)
	if config.ApiURL != "" {
		payloadSenders = append(payloadSenders, NewApiSender(httpClient, config.ApiURL, config.ApiHeaderName, config.ApiHeaderValue))
	}
	if config.WebhookURL != "" {
		payloadSenders = append(payloadSenders, NewWebhookSender(httpClient, config.WebhookURL, config.Verbose, config.FailedOnly))
	}
	if config.PushgatewayURL != "" {
		payloadSenders = append(payloadSenders, NewPushgatewaySender(httpClient, config.PushgatewayURL, config.InstanceName))
	}

	return &notifier{
		verbose:          config.Verbose,
		chatNotifiers:    notificationSenders,
		machineNotifiers: payloadSenders,
	}
}

func (n *notifier) SendChatNotifications(m string, failedCount int) error {
	var errs []error

	if len(n.chatNotifiers) == 0 {
		if n.verbose {
			log.Println("no chat notification receivers configured, no notifications sent.")
		}
		return nil
	}

	for _, sender := range n.chatNotifiers {
		if err := sender.send(m, failedCount); err != nil {
			log.Printf("failed to send notification using %T: %v", sender, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send chat notifications: %w", errors.Join(errs...))
	}

	return nil
}

func (n *notifier) SendMachineNotifications(m Results, failedCount int) error {
	var errs []error

	if len(n.machineNotifiers) == 0 {
		if n.verbose {
			log.Println("no machine interface payload receivers configured, no notifications sent.")
		}
		return nil
	}
	for _, sender := range n.machineNotifiers {
		if err := sender.send(m, failedCount); err != nil {
			log.Printf("failed to send notification using %T: %v", sender, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send machine notifications: %w", errors.Join(errs...))
	}

	return nil
}
