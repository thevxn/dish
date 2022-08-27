package messenger

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"savla-dish/pkg/config"
)

var (
	telegramURL string
)

// returns a status
func SendTelegram(rawMessage string) int {
	if rawMessage == "" && config.Verbose {
		log.Println("messager: no message given")
		return 1
	}

	// escape dish report string for Telegram
	msg := url.QueryEscape(rawMessage)

	// form the Telegram URL
	telegramURL = "https://api.telegram.org/bot" + config.TelegramBotToken + "/sendMessage?chat_id=" + config.TelegramChatID + "&text="

	req, err := http.NewRequest(http.MethodGet, telegramURL+msg, nil)
	if err != nil && config.Verbose {
		log.Println(err)
		return 1
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.Get(telegramURL + msg)
	if err != nil && config.Verbose {
		log.Println(err)
		return 1
	}

	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil && config.Verbose {
		log.Println(err)
		return 1
	}

	// write to console log if verbose flag set
	if config.Verbose {
		log.Println(telegramURL)
		log.Println(string(body))
	}

	return 0
}
