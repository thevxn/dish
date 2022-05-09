package messenger

import (
	"io/ioutil"
	"log"
	"net/http"
)

func SendMsg(msg string) {
	req := "https://api.telegram.org/bot5226521972:AAEqJJYsnBbI3umEEOtEfoHFpnPtxRzXRiM/sendMessage?chat_id=-1001512138288&text=" + msg
	resp, err := http.Get(req)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
}
