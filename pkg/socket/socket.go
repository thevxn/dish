package socket

import (
	"encoding/json"
	"log"

	"go.vxn.dev/dish/pkg/config"
)

type Result struct {
	Socket       Socket
	Passed       bool
	ResponseCode int
	Error        error
}

type SocketList struct {
	Sockets map[string]Socket `json:"items"`
}

type Socket struct {
	// ID is an unique identifier of such socket.
	ID string `json:"id"`

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
}

// 'input' should be a string like '/path/filename.json', or a HTTP URL string
func FetchSocketList(input string) SocketList {
	var list = &SocketList{}

	// fetch JSON byte reader from input URL/path
	reader, err := getStreamFromPath(input)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// got data, load struct Sockets
	err = json.NewDecoder(reader).Decode(&list)
	if err != nil {
		panic(err)
	}

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
