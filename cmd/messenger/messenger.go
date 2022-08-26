package messenger

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	TelegramBotToken *string
	TelegramChatID   *string
	telegramURL      string
	UseTelegram      *bool
	Verbose          *bool
)

// returns a status
func SendTelegram(rawMessage string) int {
	verbose := *Verbose

	if rawMessage == "" && verbose {
		log.Println("messager: no message given")
		return 1
	}

	// escape dish report string for Telegram
	msg := url.QueryEscape(rawMessage)

	// form the Telegram URL
	telegramURL = "https://api.telegram.org/bot" + *TelegramBotToken + "/sendMessage?chat_id=" + *TelegramChatID + "&text="

	req, err := http.NewRequest(http.MethodGet, telegramURL+msg, nil)
	if err != nil && verbose {
		log.Println(err)
		return 1
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.Get(telegramURL + msg)
	if err != nil && verbose {
		log.Println(err)
		return 1
	}

	defer resp.Body.Close()

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && verbose {
		log.Println(err)
		return 1
	}

	// write to console log if verbose flag set
	if verbose {
		log.Println(telegramURL)
		log.Println(string(body))
	}

	return 0
}
