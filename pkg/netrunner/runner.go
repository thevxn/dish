package netrunner

import (
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/socket"
)

func TestSocket(sock socket.Socket, channel chan<- socket.Result, wg *sync.WaitGroup) {
	defer wg.Done()

	regex, err := regexp.Compile("^(http|https)://")
	if err != nil {
		log.Println("Failed to create new regex object")

		if channel != nil {
			close(channel)
		}
		return
	}

	result := socket.Result{
		Socket: sock,
	}

	if !regex.MatchString(sock.Host) {
		// Testing raw host and port (tcp), report only unsuccessful tests; exclusively non-HTTP/S sockets
		result.Error = rawConnect(sock)
		result.Passed = result.Error == nil

		sendResult(channel, result)
		return
	}

	result.Passed, result.ResponseCode, result.Error = checkSite(sock)
	sendResult(channel, result)
}

// rawConnect function for direct host:port socket check
func rawConnect(sock socket.Socket) error {
	endpoint := net.JoinHostPort(sock.Host, strconv.Itoa(sock.Port))
	timeout := time.Duration(time.Second * time.Duration(config.Timeout))

	if config.Verbose {
		log.Println("runner: rawconnect: " + endpoint)
	}

	// open the socket
	conn, err := net.DialTimeout("tcp", endpoint, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

// checkHTTPCode function for response and expected HTTP codes comparison
// panics if it fails to convert expected code to int
func checkHTTPCode(responseCode int, expectedCodes []int) bool {
	for _, code := range expectedCodes {
		if responseCode == code {
			return true
		}
	}
	return false
}

// checkSite executes test over HTTP/S endpoints exclusively
func checkSite(socket socket.Socket) (bool, int, error) {
	// Configure HTTP client
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}
	url := socket.Host + ":" + strconv.Itoa(socket.Port) + socket.PathHTTP

	if config.Verbose {
		log.Println("runner: checksite:", url)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, 0, err
	}
	req.Header.Set("User-Agent", "savla-dish/1.6")

	// open socket --- Head to url
	resp, err := client.Do(req)
	if err != nil {
		return false, 0, err
	}

	// fetch StatusCode for HTTP expected code comparison
	if resp != nil {
		defer resp.Body.Close()
		return checkHTTPCode(resp.StatusCode, socket.ExpectedHTTPCodes), resp.StatusCode, nil
	}

	return true, resp.StatusCode, nil
}

func sendResult(channel chan<- socket.Result, result socket.Result) {
	if channel != nil {
		channel <- result
		close(channel)
	}
}
