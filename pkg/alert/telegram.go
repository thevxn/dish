package alert

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const baseURL = "https://api.telegram.org"
const messageTitle = "\U0001F4E1 <b>dish run results</b>:"

type telegramSender struct {
	httpClient    *http.Client
	chatID        string
	token         string
	verbose       bool
	notifySuccess bool
}

func NewTelegramSender(httpClient *http.Client, chatID string, token string, verbose bool, notifySuccess bool) *telegramSender {
	return &telegramSender{
		httpClient,
		chatID,
		token,
		verbose,
		notifySuccess,
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

	// Send the message
	res, err := s.httpClient.Get(fullURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response code received from Telegram (expected: %d, got: %d)", http.StatusOK, res.StatusCode)
	}

	// Write the body to console if verbose flag set
	if s.verbose {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		log.Println("telegram response:", string(body))
	}

	log.Println("telegram alert sent")

	return nil
}
