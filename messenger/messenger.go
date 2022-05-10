package messenger

import (
	"io/ioutil"
	"log"
	"net/http"
)

func SendMsg(msg string) {
	if msg == "" {
		log.Print("no message given")
		return
	}

	chat_id := "-1001248157564"
	t_endpoint := "https://api.telegram.org/bot5226521972:AAEqJJYsnBbI3umEEOtEfoHFpnPtxRzXRiM/sendMessage?chat_id=" + chat_id + "&text=" + msg
	req, err := http.NewRequest("GET", t_endpoint, nil)
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
