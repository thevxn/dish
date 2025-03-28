package alert

import (
	"log"
	"net/http"

	"go.vxn.dev/dish/pkg/config"
)

func HandleAlerts(messengerText string, results Results, failedCount int, config *config.Config) {
	notifier := NewNotifier(http.DefaultClient, config)
	if err := notifier.SendChatNotifications(messengerText, failedCount); err != nil {
		log.Printf("error sending chat notifications: %v", err)
	}
	if err := notifier.SendMachineNotifications(results, failedCount); err != nil {
		log.Printf("error sending machine notifications: %v", err)
	}
}
