package alert

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"go.vxn.dev/dish/pkg/config"
)

// returns an error if given message is an empty string or if a request cannot be sent
func SendTelegram(rawMessage string) error {
	if rawMessage == "" {
		return errors.New("empty message string given")
	}

	// form the Telegram URL
	telegramURL := "https://api.telegram.org/bot" + config.TelegramBotToken + "/sendMessage?chat_id=" + config.TelegramChatID + "&disable_web_page_preview=True&parse_mode=HTML&text="

	msg := "<b>dish run results</b>:\n\n" + rawMessage

	// escape dish report string for Telegram
	msg = url.QueryEscape(msg)

	req, err := http.NewRequest(http.MethodGet, telegramURL+msg, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.Get(telegramURL + msg)
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
	if config.Verbose {
		log.Println(telegramURL)
		log.Println(string(body))
	}

	return nil
}
