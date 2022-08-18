package socket

import (
	"encoding/json"
	"log"
)

type Sockets struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	Name              string `json:"socket_name"`
	Host              string `json:"host_name"`
	Port              int    `json:"port_tcp"`
	ExpectedHTTPCodes []int  `json:"expected_http_code_array"`
	PathHTTP          string `json:"path_http"`

	// can be blank, dish name here is meant as socket list owner/target from remote RESTful API server
	DishName string  `json:"dish_list"`
	Results  Results `json:"dish_results"`
}

type Results struct {
	HTTPCode      int   `json:"http_response_code"`
	SocketReached bool  `json:"socket_reached" default:true`
	Error         error `json:"error_message"`
}

var Verbose *bool

// FetchSocketList method ...
// 'input' should be a string like '/path/filename.json', or a HTTP URL string
func FetchSocketList(input string) (socketsPointer *Sockets) {
	// fetch JSON byte stream from input URL/path
	stream := getStreamFromInput(input)
	if stream == nil {
		log.Fatalln("socket: fatal: no JSON stream to get socket list")
		return nil
	}

	// got stream, load struct Sockets
	var sockets Sockets
	json.Unmarshal(*stream, &sockets)

	// write JSON data to console
	if *Verbose {
		for _, socket := range sockets.Sockets {
			log.Println("socket: Host:", socket.Host)
			log.Println("socket: Port:", socket.Port)
			log.Println("socket: ExpectedHTTPCodes:", socket.ExpectedHTTPCodes)
		}
	}

	return &sockets
}
