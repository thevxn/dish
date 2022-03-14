package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Sockets struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	Endpoint string `json:"endpoint`
	Port     int    `json:"port"`
}

func main() {
	
	jsonFile, err := os.Open("sockets.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened sockets.json")
	
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var sockets Sockets

	json.Unmarshal(byteValue, &sockets)

	for i := 0; i < len(sockets.Sockets); i++ {
		fmt.Println("Endpoint: " + sockets.Sockets[i].Endpoint)
		fmt.Print("Port: ")
		fmt.Print(sockets.Sockets[i].Port)
	}

}
