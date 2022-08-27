package netrunner

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
func RawConnect(socket *socket.Socket) error {
	endpoint := net.JoinHostPort(socket.Host, socket.Port)
	timeout := time.Duration(time.Second * 5)

	if config.Verbose {
		log.Println("runner: rawconnect: " + endpoint)
	}

	// open the socket
	conn, err := net.DialTimeout("tcp", endpoint, timeout)
	if err != nil {
		if config.Verbose {
			log.Println("runner: rawconnect: conn error:", endpoint)
			log.Println(err)
		}
		socket.Results.Error = err
		return err
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	return nil
}

// checkHTTPCode function for response and expected HTTP codes comparison
// panics if it fails to convert expected code to int
func checkHTTPCode(responseCode int, expectedCodes []string) bool {
	for _, code := range expectedCodes {
		code, err := strconv.Atoi(code)
		if err != nil {
			panic(err)
		}

		if responseCode != code {
			return false
		}
	}
	return true
}

// CheckSite executes test over HTTP/S endpoints exclusively
func CheckSite(socket *socket.Socket) bool {
	// config http client
	netClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	url := socket.Host + ":" + socket.Port + socket.PathHTTP

	if config.Verbose {
		log.Println("runner: checksite:", url)
	}

	// open socket --- Head to url
	resp, err := netClient.Head(url)
	if err != nil {
		if config.Verbose {
			log.Println(err)
		}

		socket.Results.Error = err
		return true
	}

	// fetch StatusCode for HTTP expected code comparison
	if resp != nil {
		defer resp.Body.Close()
		socket.Results.HTTPCode = resp.StatusCode
		return checkHTTPCode(resp.StatusCode, socket.ExpectedHTTPCodes)
	}

	return false
}
