package socket

import (
	//"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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


