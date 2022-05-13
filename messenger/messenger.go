// +build dev

package messenger

import (
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// see http://docs.savla.su/projects/telegram-bots
	// TODO: do not hardcode telegram bots!
	botToken string = "5226521972:AAEqJJYsnBbI3umEEOtEfoHFpnPtxRzXRiM"
	chatID string = "-1001248157564"
	telegramURL string = "https://api.telegram.org/bot" + botToken + "/sendMessage?chat_id=" + chatID + "&text="
)

func SendMsg(msg string) (status int) {
	if msg == "" {
		log.Println("messager: no message given")
		return 1
	}

	req, err := http.NewRequest("GET", telegramURL + msg, nil); if err != nil {
		log.Println(err)
		return 1
	}

	req.Header.Set("Content-Type", "application/json")
	//resp, err := http.DefaultClient.Do(req)
	resp, err := http.Get(telegramURL + msg); if err != nil {
		log.Println(err)
		return 1
	}

	// We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body); if err != nil {
		log.Println(err)
		return 1
	}

	// Convert the body to type string
	resp.Body.Close()
	sb := string(body)
	log.Printf(sb)

	return 0
}
