package alert

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

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

	// form the Telegram URL
	telegramURL := "https://api.telegram.org/bot" + s.token + "/sendMessage?chat_id=" + s.chatID + "&disable_web_page_preview=True&parse_mode=HTML&text="

	msg := "<b>dish run results</b>:\n\n" + rawMessage

	// escape dish report string for Telegram
	msg = url.QueryEscape(msg)

	resp, err := s.httpClient.Get(telegramURL + msg)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// write to console log if verbose flag set
	if s.verbose {
		log.Println(telegramURL)
		log.Println(string(body))
	}

	return nil
}
