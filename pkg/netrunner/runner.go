package netrunner

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"sync"
	"time"

	"go.vxn.dev/dish/pkg/socket"
)

const agentVersion = "1.9"

func TestSocket(sock socket.Socket, channel chan<- socket.Result, wg *sync.WaitGroup, timeoutSeconds uint, verbose bool) {
	defer wg.Done()

	// TODO: move out?
	regex, err := regexp.Compile("^(http|https)://")
	if err != nil {
		log.Printf("regex compilation failed: %v", err)

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
		result.Error = rawConnect(sock, timeoutSeconds, verbose)
		result.Passed = result.Error == nil

		sendResult(channel, result)
		return
	}

	result.Passed, result.ResponseCode, result.Error = checkSite(sock, timeoutSeconds, verbose)
	sendResult(channel, result)
}

// rawConnect performs a direct host:port socket check
func rawConnect(sock socket.Socket, timeoutSeconds uint, verbose bool) error {
	endpoint := net.JoinHostPort(sock.Host, strconv.Itoa(sock.Port))
	timeout := time.Duration(time.Second * time.Duration(timeoutSeconds))

	if verbose {
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

// checkSite executes test over HTTP/S endpoints exclusively
func checkSite(socket socket.Socket, timeoutSeconds uint, verbose bool) (bool, int, error) {
	// Configure HTTP client
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	url := socket.Host + ":" + strconv.Itoa(socket.Port) + socket.PathHTTP

	if verbose {
		log.Println("runner: checksite:", url)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, 0, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("dish/%s", agentVersion))

	// open socket --- Head to url
	resp, err := client.Do(req)
	if err != nil {
		return false, 0, err
	}

	// fetch StatusCode for HTTP expected code comparison
	if resp != nil {
		defer resp.Body.Close()

		if !slices.Contains(socket.ExpectedHTTPCodes, resp.StatusCode) {
			err = fmt.Errorf("expected codes: %v, got %d", socket.ExpectedHTTPCodes, resp.StatusCode)
		}

		return slices.Contains(socket.ExpectedHTTPCodes, resp.StatusCode), resp.StatusCode, err
	}

	return true, 0, nil
}

// sendResult sends the result of a check to the result channel and closes it
func sendResult(channel chan<- socket.Result, result socket.Result) {
	if channel != nil {
		channel <- result
		close(channel)
	}
}
