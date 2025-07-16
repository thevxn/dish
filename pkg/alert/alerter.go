// Package alert provides functionality to handle alert and result submission
// to different text (e.g. Telegram) and machine (e.g. webhooks) integration channels.
package alert

import (
	"net/http"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
)

// alerter provides a centralized method of alerting the configured channels with the results of the performed checks
// while hiding implementation details of the channels.
type alerter struct {
	logger logger.Logger
}

// NewAlerter returns a new instance of alerter using the provided logger.
func NewAlerter(l logger.Logger) *alerter {
	if l == nil {
		return nil
	}

	return &alerter{
		logger: l,
	}
}

// HandleAlerts notifies all configured channels with either the provided message (if text channel) or the structured results (if machine channel).
func (a *alerter) HandleAlerts(messengerText string, results *Results, failedCount int, config *config.Config) {
	if results == nil || config == nil {
		return
	}

	notifier := NewNotifier(http.DefaultClient, config, a.logger)
	if err := notifier.SendChatNotifications(messengerText, failedCount); err != nil {
		a.logger.Errorf("some error(s) encountered when sending chat notifications: \n%v", err)
	}
	if err := notifier.SendMachineNotifications(results, failedCount); err != nil {
		a.logger.Errorf("some error(s) encountered when sending machine notifications: \n%v", err)
	}
}
