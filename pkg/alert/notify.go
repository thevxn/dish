package alert

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"savla-dish/pkg/config"
)

// returns true or fals whether it sends
func SendTelegram(rawMessage string) error {
	if rawMessage == "" && config.Verbose {
		panic("messager: no message given")
	}

	// escape dish report string for Telegram
	msg := url.QueryEscape(rawMessage)

	// form the Telegram URL
	telegramURL := "https://api.telegram.org/bot" + config.TelegramBotToken + "/sendMessage?chat_id=" + config.TelegramChatID + "&text="

	req, err := http.NewRequest(http.MethodGet, telegramURL+msg, nil)
	if err != nil {
		if config.Verbose {
			log.Println(err)
		}
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.Get(telegramURL + msg)
	if err != nil {
		if config.Verbose {
			log.Println(err)
		}
		return err
	}
	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if config.Verbose {
			log.Println(err)
		}
		return err
	}

	// write to console log if verbose flag set
	if config.Verbose {
		log.Println(telegramURL)
		log.Println(string(body))
	}

	return nil
}
