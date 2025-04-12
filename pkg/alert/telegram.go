package alert

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"go.vxn.dev/dish/pkg/config"
)

const (
	baseURL      = "https://api.telegram.org"
	messageTitle = "\U0001F4E1 <b>dish run results</b>:" // 📡
)

type telegramSender struct {
	httpClient    HTTPClient
	chatID        string
	token         string
	verbose       bool
	notifySuccess bool
}

func NewTelegramSender(httpClient HTTPClient, config *config.Config) *telegramSender {
	return &telegramSender{
		httpClient,
		config.TelegramChatID,
		config.TelegramBotToken,
		config.Verbose,
		config.TextNotifySuccess,
	}
}

func (s *telegramSender) send(rawMessage string, failedCount int) error {
	// If no checks failed and success should not be notified, there is nothing to send
	if failedCount == 0 && !s.notifySuccess {
		if s.verbose {
			log.Printf("no sockets failed, nothing will be sent to Telegram")
		}
		return nil
	}

	// Construct the Telegram URL with params and the message
	telegramURL := fmt.Sprintf("%s/bot%s/sendMessage", baseURL, s.token)

	params := url.Values{}
	params.Set("chat_id", s.chatID)
	params.Set("disable_web_page_preview", "true")
	params.Set("parse_mode", "HTML")
	params.Set("text", messageTitle+"\n\n"+rawMessage)

	fullURL := telegramURL + "?" + params.Encode()

	err := handleSubmit(s.httpClient, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("error submitting Telegram alert: %w", err)
	}

	log.Println("Telegram alert sent")

	return nil
}
