// Package alert provides functionality to handle alert and result submission
// to different text (e.g. Telegram) and machine (e.g. webhooks) integration channels.
package alert

import (
	"log"
	"net/http"

	"go.vxn.dev/dish/pkg/config"
)

func HandleAlerts(messengerText string, results *Results, failedCount int, config *config.Config) {
	notifier := NewNotifier(http.DefaultClient, config)
	if err := notifier.SendChatNotifications(messengerText, failedCount); err != nil {
		log.Printf("some error(s) encountered when sending chat notifications: \n%v", err)
	}
	if err := notifier.SendMachineNotifications(results, failedCount); err != nil {
		log.Printf("some error(s) encountered when sending machine notifications: \n%v", err)
	}
}
