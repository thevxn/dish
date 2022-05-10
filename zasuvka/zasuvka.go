// Zasuvka package to parse JSON input file
package zasuvka

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Sockets struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	Name         string `json:"socket_name"`
	Host	     string `json:"endpoint_url"`
	Port         int    `json:"port_tcp"`
	ExpectedCode int    `json:"expected_http_port"`
}

func GibPole(f string, debug bool) (s Sockets) {
	//jsonFile, err := os.Open("sockets.json")
	jsonFile, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	if debug {
		log.Println("Successfully Opened sockets.json")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var sockets Sockets

	json.Unmarshal(byteValue, &sockets)

	if debug {
		for i := 0; i < len(sockets.Sockets); i++ {
			log.Printf("zasuvka: Host: %s", sockets.Sockets[i].Host)
			log.Printf("zasuvka: Port: %d", sockets.Sockets[i].Port)
		}
	}

	return sockets
}
