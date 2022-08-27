package socket

import (
	"encoding/json"
	"log"
	"savla-dish/pkg/config"
)

type SocketList struct {
	Sockets []Socket
}

type Socket struct {
	Name              string   `json:"socket_name"`
	Host              string   `json:"host_name"`
	Port              string   `json:"port_tcp"`
	ExpectedHTTPCodes []string `json:"expected_http_code_array"`
	PathHTTP          string   `json:"path_http"`

	// can be blank, dish name here is meant as socket list owner/target from remote RESTful API server
	DishName string  `json:"dish_list"`
	Results  Results `json:"dish_results"`
}

type Results struct {
	HTTPCode      int   `json:"http_response_code"`
	SocketReached bool  `json:"socket_reached" default:"true"`
	Error         error `json:"error_message"`
}

// 'input' should be a string like '/path/filename.json', or a HTTP URL string
func FetchSocketList(input string) (list SocketList) {
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

	return list
}
