// Zasuvka package to parse JSON input file
package zasuvka

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Sockets struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	Name         string `json:"socket_name"`
	Endpoint     string `json:"endpoint_url`
	Port         int    `json:"port_tcp"`
	ExpectedCode int    `json:"expected_http_port"`
}

func GibPole(f string, debug bool) (s Sockets) {
	//jsonFile, err := os.Open("sockets.json")
	jsonFile, err := os.Open(f)

	if err != nil {
		//fmt.Println(err)
		log.Fatal(err)
	}

	if debug {
		fmt.Println("Successfully Opened sockets.json")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var sockets Sockets

	json.Unmarshal(byteValue, &sockets)

	if debug {
		for i := 0; i < len(sockets.Sockets); i++ {
			fmt.Println("Endpoint: " + sockets.Sockets[i].Endpoint)
			fmt.Print("Port: ")
			fmt.Print(sockets.Sockets[i].Port)
		}
	}

	return sockets
}
