package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Users struct which contains
// an array of users
type Sockets struct {
	Sockets []Socket `json:"sockets"`
}

// User struct which contains a name
// a type and a list of social links
type Socket struct {
	Endpoint string `json:"endpoint`
	Port     int    `json:"port"`
}

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("sockets.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened sockets.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var sockets Sockets

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &sockets)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(sockets.Sockets); i++ {
		fmt.Println("Endpoint: " + sockets.Sockets[i].Endpoint)
		fmt.Print("Port: ")
		fmt.Print(sockets.Sockets[i].Port)
	}

}
