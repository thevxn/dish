package alert

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
)

type Results struct {
	Map map[string]bool `json:"dish_results"`
}

type ChatNotifier interface {
	send(string, int) error
}

type MachineNotifier interface {
	send(*Results, int) error
}

type notifier struct {
	verbose          bool
	chatNotifiers    []ChatNotifier
	machineNotifiers []MachineNotifier
	logger           logger.Logger
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

// NewNotifier creates a new instance of notifier. Based on the flags used, it spawns new instances of ChatNotifiers (e.g. Telegram) and MachineNotifiers (e.g. Webhooks) and stores them on the notifier struct to be used for alert notifications.
func NewNotifier(httpClient HTTPClient, config *config.Config, logger logger.Logger) *notifier {
	if logger == nil {
		return nil
	}

	if config == nil {
		logger.Error("nil pointer to config")
		return nil
	}

	// Set chat integrations to be notified (e.g. Telegram)
	notificationSenders := make([]ChatNotifier, 0)

	// Telegram
	if config.TelegramBotToken != "" && config.TelegramChatID != "" {
		notificationSenders = append(notificationSenders, NewTelegramSender(httpClient, config, logger))
	}

	// Set machine interface integrations to be notified (e.g. Webhooks)
	payloadSenders := make([]MachineNotifier, 0)

	// Remote API
	if config.ApiURL != "" {
		apiSender, err := NewAPISender(httpClient, config, logger)
		if err != nil {
			logger.Error("error creating new remote API sender: ", err)
		} else {
			payloadSenders = append(payloadSenders, apiSender)
		}
	}

	// Webhooks
	if config.WebhookURL != "" {
		webhookSender, err := NewWebhookSender(httpClient, config, logger)
		if err != nil {
			logger.Error("error creating new webhook sender: ", err)
		} else {
			payloadSenders = append(payloadSenders, webhookSender)
		}
	}

	// Pushgateway
	if config.PushgatewayURL != "" {
		pgwSender, err := NewPushgatewaySender(httpClient, config, logger)
		if err != nil {
			logger.Error("error creating new Pushgateway sender:", err)
		} else {
			payloadSenders = append(payloadSenders, pgwSender)
		}
	}

	// Discord
	fmt.Println(config)
	fmt.Println("DIscord channel ID:", config.DiscordChannelID)
	fmt.Println("Token ID: ", config.DiscordBotToken)
	if config.DiscordChannelID != "" && config.DiscordBotToken != "" {
		fmt.Println("here")
		notificationSenders = append(notificationSenders, NewDiscordSender(httpClient, config, logger))
	}

	return &notifier{
		verbose:          config.Verbose,
		chatNotifiers:    notificationSenders,
		machineNotifiers: payloadSenders,
		logger:           logger,
	}
}

func (n *notifier) SendChatNotifications(m string, failedCount int) error {
	var errs []error

	if len(n.chatNotifiers) == 0 {
		n.logger.Debug("no chat notification receivers configured, no notifications will be sent")

		return nil
	}

	for _, sender := range n.chatNotifiers {
		if err := sender.send(m, failedCount); err != nil {
			n.logger.Errorf("failed to send notification using %T: %v", sender, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (n *notifier) SendMachineNotifications(m *Results, failedCount int) error {
	var errs []error

	if len(n.machineNotifiers) == 0 {
		n.logger.Debug("no machine interface payload receivers configured, no notifications will be sent")

		return nil
	}

	for _, sender := range n.machineNotifiers {
		if err := sender.send(m, failedCount); err != nil {
			n.logger.Errorf("failed to send notification using %T: %v", sender, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
