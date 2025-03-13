package socket

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"go.vxn.dev/dish/pkg/config"
)

// fetchRemoteStream sends a GET HTTP request to remote RESTful API endpoint, returns response body
// 'url' argument should be a full-quality URL to remote http server, e.g. http://api.example.com:5569/stream?query=variable
func fetchRemoteStream(url string) (io.ReadCloser, error) {
	// try URL
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")

	if config.HeaderName != "" && config.HeaderValue != "" {
		req.Header.Set(config.HeaderName, config.HeaderValue)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Fatal("not a HTTP/200 status code on socket list fetch --- got " + strconv.Itoa(resp.StatusCode))
	}

	body := resp.Body

	return body, nil
}

func fetchFileStream(input string) (io.ReadCloser, error) {
	jsonFile, err := os.Open(input)
	if err != nil {
		return nil, err
	}

	return jsonFile, nil
}

// getStreamFromPath tries to load data from given input. It checks whether input is a file path or url
func getStreamFromPath(input string) (io.ReadCloser, error) {
	// Check if input is an url
	match, err := regexp.MatchString("^(http|https)://", input)
	if err != nil {
		return nil, err
	}

	if !match {
		reader, err := fetchFileStream(input)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}

	reader, err := fetchRemoteStream(input)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
