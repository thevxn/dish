package messenger

import (
	"io/ioutil"
	"log"
	"net/http"
)

const (
	bot_token string = "5226521972:AAEqJJYsnBbI3umEEOtEfoHFpnPtxRzXRiM"
	chat_id string = "-1001248157564"
	t_endpoint string = "https://api.telegram.org/bot" + bot_token + "/sendMessage?chat_id=" + chat_id + "&text="
)

func SendMsg(msg string) {
	if msg == "" {
		log.Print("no message given")
		return
	}

	req, err := http.NewRequest("GET", t_endpoint + msg, nil)
	if err != nil {
		// status = err
		// continue
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	//resp, err := http.Get(t_endpoint + msg)
	if err != nil {
		log.Fatal(err)
	}

	// We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Convert the body to type string
	sb := string(body)
	log.Printf(sb)
}
