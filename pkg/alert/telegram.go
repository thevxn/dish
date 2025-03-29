package alert

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const baseURL = "https://api.telegram.org"

type telegramSender struct {
	httpClient *http.Client
	chatID     string
	token      string
	verbose    bool
	failedOnly bool
}

func NewTelegramSender(httpClient *http.Client, chatID string, token string, verbose bool, failedOnly bool) *telegramSender {
	return &telegramSender{
		httpClient,
		chatID,
		token,
		verbose,
		failedOnly,
	}
}

func (s *telegramSender) send(rawMessage string, failedCount int) error {
	if rawMessage == "" {
		return errors.New("empty message string given")
	}

	// If there are no failed sockets and we only wish to be notified when they fail, there is nothing to do
	if failedCount == 0 && s.failedOnly {
		log.Printf("%T: no failed sockets and failedOnly == true, nothing will be sent", s)
		return nil
	}

	// Construct the Telegram URL with params and the message
	telegramURL := fmt.Sprintf("%s/bot%s/sendMessage", baseURL, s.token)

	params := url.Values{}
	params.Set("chat_id", s.chatID)
	params.Set("disable_web_page_preview", "true")
	params.Set("parse_mode", "HTML")
	params.Set("text", "<b>dish run results</b>:\n\n"+rawMessage)

	fullURL := telegramURL + "?" + params.Encode()

	// Send the message
	resp, err := s.httpClient.Get(fullURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response code received from Telegram (expected: %d, got: %d)", http.StatusOK, resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Writethe body to console if verbose flag set
	if s.verbose {
		log.Println("telegram response:", string(body))
	}

	log.Println("telegram alert sent")

	return nil
}
