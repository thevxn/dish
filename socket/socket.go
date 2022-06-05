package socket

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Sockets struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	Name              string `json:"socket_name"`
	Host              string `json:"host_name"`
	Port              int    `json:"port_tcp"`
	ExpectedHttpCodes []int  `json:"expected_http_code_array"`
	PathHttp          string `json:"path_http"`
}

// fetchRemoteStream sends a GET HTTP request to remote RESTful API endpoint, returns JSON stream
// 'url' argument should be a full-quality URL to remote http server, e.g. http://api.example.com:5569/stream?query=variable
func fetchRemoteStream(url string) (byteStream *[]byte) {
	// try URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//log.Println(err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.Get(url)
	if err != nil {
		//log.Println(err)
		return nil
	}

	// read response body stream
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Println(err)
		return nil
	}

	// Convert the body to type string
	defer resp.Body.Close()
	return &body
}

// fetchFileStream
func fetchFileStream(input string) (byteStream *[]byte) {
	//jsonFile, err := os.Open("sockets.json")
	jsonFile, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer jsonFile.Close()

	// use local var as "buffer", then return pointer to data
	stream, _ := ioutil.ReadAll(jsonFile)
	return &stream
}

// getStreamFromInput metamethod ('case-like macro') tries to load data stream from given source; returns pointer to stream
func getStreamFromInput(input string) (byteStream *[]byte) {
	// try to open stream, if URL, else open file
	stream := fetchRemoteStream(input)
	if stream != nil {
		return stream
	}

	// use input string as path to a file
	stream = fetchFileStream(input)
	if stream != nil {
		return stream
	}

	return nil
}

// FetchSocketList method ...
// 'input' should be a string like '/path/filename.json', or a HTTP URL string
func FetchSocketList(input string, verbose bool) (socketsPointer *Sockets) {
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
	if verbose {
		for _, socket := range sockets.Sockets {
			log.Println("socket: Host:", socket.Host)
			log.Println("socket: Port:", socket.Port)
			log.Println("socket: ExpectedHttpCodes:", socket.ExpectedHttpCodes)
		}
	}

	return &sockets
}
