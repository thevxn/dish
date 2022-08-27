package runner

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"savla-dish/pkg/config"
	"savla-dish/pkg/socket"
)

// RawConnect function to direct host:port socket check
func RawConnect(socket socket.Socket) (int, error) {
	endpoint := net.JoinHostPort(socket.Host, socket.Port)
	timeout := time.Duration(5 * time.Second)

	if config.Verbose {
		log.Println("runner: rawconnect: " + endpoint)
	}

	// open the socket
	conn, err := net.DialTimeout("tcp", endpoint, timeout)

	// close open conn after 5 seconds
	//conn.SetReadDeadline(time.Second*5)

	// prolly more possible to get not-nil err, than not-nil conn
	// see --> https://stackoverflow.com/a/56336811
	if err != nil {
		if config.Verbose {
			log.Println("runner: rawconnect: conn error:", endpoint)
			log.Println(err)
		}
		socket.Results.Error = err
		return 1, err
	}

	if conn != nil {
		conn.Close()
		return 0, nil
	}

	// unexpected error
	return 2, nil
}

// checkHTTPCode function for response and expected HTTP codes comparison
func checkHTTPCode(responseCode int, expectedCodes []string) int {
	for _, code := range expectedCodes {
		if code, err := strconv.Atoi(code); responseCode == code {
			if err != nil {
				panic(err)
			}
			// site is OK! do not report ok sites?
			return 0
		}
	}
	return responseCode
}

// CheckSite executes test over HTTP/S endpoints exclusively
func CheckSite(socket socket.Socket) int {
	// config http client
	var netClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	url := socket.Host + ":" + socket.Port + socket.PathHTTP

	if config.Verbose {
		log.Println("runner: checksite:", url)
	}

	// open socket --- give Head
	resp, err := netClient.Head(url)
	if err != nil {
		if config.Verbose {
			log.Println(err)
		}

		socket.Results.Error = err
		return 0
	}

	// fetch StatusCode for HTTP expected code comparison
	if resp != nil {
		defer resp.Body.Close()
		socket.Results.HTTPCode = resp.StatusCode
		return checkHTTPCode(resp.StatusCode, socket.ExpectedHTTPCodes)
	}

	return 2
}
