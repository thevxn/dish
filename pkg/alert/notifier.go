package alert

import (
	"errors"
	"io"
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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

// NewNotifier creates a new instance of notifier. Based on the flags used, it spawns new instances of ChatNotifiers (e.g. Telegram) and MachineNotifiers (e.g. Webhooks) and stores them on the notifier struct to be used for alert notifications.
func NewNotifier(httpClient HTTPClient, config *config.Config) *notifier {
	// Set chat integrations to be notified (e.g. Telegram)
	notificationSenders := make([]ChatNotifier, 0)

	// Telegram
	if config.TelegramBotToken != "" && config.TelegramChatID != "" {
		notificationSenders = append(notificationSenders, NewTelegramSender(httpClient, config.TelegramChatID, config.TelegramBotToken, config.Verbose, config.TextNotifySuccess))
	}

	// Set machine interface integrations to be notified (e.g. Webhooks)
	payloadSenders := make([]MachineNotifier, 0)

	// Remote API
	if config.ApiURL != "" {
		apiSender, err := NewAPISender(httpClient, config.ApiURL, config.ApiHeaderName, config.ApiHeaderValue, config.Verbose, config.MachineNotifySuccess)
		if err != nil {
			log.Println("error creating new remote API sender:", err)
		} else {
			payloadSenders = append(payloadSenders, apiSender)
		}
	}

	// Webhooks
	if config.WebhookURL != "" {
		webhookSender, err := NewWebhookSender(httpClient, config.WebhookURL, config.Verbose, config.MachineNotifySuccess)
		if err != nil {
			log.Println("error creating new webhook sender:", err)
		} else {
			payloadSenders = append(payloadSenders, webhookSender)
		}
	}

	// Pushgateway
	if config.PushgatewayURL != "" {
		pgwSender, err := NewPushgatewaySender(httpClient, config.PushgatewayURL, config.InstanceName, config.Verbose, config.MachineNotifySuccess)
		if err != nil {
			log.Println("error creating new Pushgateway sender:", err)
		} else {
			payloadSenders = append(payloadSenders, pgwSender)
		}
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
		log.Println("no chat notification receivers configured, no notifications will be sent")

		return nil
	}

	for _, sender := range n.chatNotifiers {
		if err := sender.send(m, failedCount); err != nil {
			log.Printf("failed to send notification using %T: %v", sender, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (n *notifier) SendMachineNotifications(m Results, failedCount int) error {
	var errs []error

	if len(n.machineNotifiers) == 0 {
		log.Println("no machine interface payload receivers configured, no notifications will be sent")

		return nil
	}
	for _, sender := range n.machineNotifiers {
		if err := sender.send(m, failedCount); err != nil {
			log.Printf("failed to send notification using %T: %v", sender, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
