package socket

import (
	"encoding/json"
	"log"

	"dish/pkg/config"

	//"go.savla.dev/swapi/dish"
)

type SocketList struct {
	Sockets map[string]Socket `json:"items"`
	//Sockets map[string]dish.Socket `json:"items"`
}

type Socket struct {
	// ID is an unique identifier of such socket.
	ID string `json:"socket_id"`

	// Socket name, unique identificator, snake_cased.
	Name string `json:"socket_name"`

	// Remote endpoint hostname or URL.
	Host string `json:"host_name"`

	// Remote port to assemble a socket.
	Port int `json:"port_tcp"`

	// HTTP Status Codes expected when giving the endpoint a HEAD/GET request.
	ExpectedHTTPCodes []int `json:"expected_http_code_array"`

	// HTTP Path to test on Host.
	PathHTTP string `json:"path_http"`

	// Can be blank, dish name here is meant as socket list owner/target from remote RESTful API server.
	DishName string `json:"dish_list"`
}

// 'input' should be a string like '/path/filename.json', or a HTTP URL string
func FetchSocketList(input string) SocketList {
	var list = &SocketList{}

	// fetch JSON byte reader from input URL/path
	reader, err := getStreamFromPath(input)
	if err != nil {
		panic(err)
	}

	// got data, load struct Sockets
	err = json.NewDecoder(reader).Decode(&list)
	if err != nil {
		panic(err)
	}
	reader.Close()

	// write JSON data to console
	if config.Verbose {
		for _, socket := range list.Sockets {
			log.Println("socket: Host:", socket.Host)
			log.Println("socket: Port:", socket.Port)
			log.Println("socket: ExpectedHTTPCodes:", socket.ExpectedHTTPCodes)
		}
	}

	return *list
}
