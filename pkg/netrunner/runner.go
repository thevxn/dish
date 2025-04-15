package netrunner

import (
	"context"
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

const agentVersion = "1.10"

// RunSocketTest is meant to be invoked in a separate goroutine.
// It runs a test for the given socket. The test result is sent through the given
// channel. If the test fails to start then the error is logged to stdout and no
// result is sent. When this func returns, it calls Done() on the WaitGroup and
// the channel is closed.
func RunSocketTest(sock socket.Socket, out chan<- socket.Result, wg *sync.WaitGroup, timeoutSeconds uint, verbose bool) {
	defer wg.Done()
	defer close(out)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	runner, err := NewNetRunner(sock, verbose)
	if err != nil {
		log.Printf("failed to test socket: %v", err.Error())
		return
	}

	out <- runner.RunTest(ctx, sock)
}

// NetRunner is used to run tests for a socket.
type NetRunner interface {
	RunTest(ctx context.Context, sock socket.Socket) socket.Result
}

// NewNetRunner determines the protocol used for the socket test and
// creates a new NetRunner for it.
func NewNetRunner(sock socket.Socket, verbose bool) (NetRunner, error) {
	exp, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return nil, fmt.Errorf("regex compilation failed: %w", err)
	}

	if exp.MatchString(sock.Host) {
		return httpRunner{client: &http.Client{}, verbose: verbose}, nil
	}

	return tcpRunner{verbose: verbose}, nil
}

type tcpRunner struct {
	verbose bool
}

// RunTest is used to test TCP sockets. It opens a TCP connection with the given socket.
// The test passes if the connection is successfully opened with no errors.
func (runner tcpRunner) RunTest(ctx context.Context, sock socket.Socket) socket.Result {
	endpoint := net.JoinHostPort(sock.Host, strconv.Itoa(sock.Port))

	if runner.verbose {
		log.Println("tcprunner: connect: " + endpoint)
	}

	d := net.Dialer{}

	conn, err := d.DialContext(ctx, "tcp", endpoint)
	if err != nil {
		return socket.Result{Socket: sock, Error: err, Passed: false}
	}
	defer conn.Close()

	return socket.Result{Socket: sock, Passed: true}
}

type httpRunner struct {
	client  *http.Client
	verbose bool
}

// RunTest is used to test HTTP/S endpoints exclusively. It executes a HTTP GET
// request to the given socket. The test passes if the request did not end with
// an error and the response status matches the expected HTTP codes.
func (runner httpRunner) RunTest(ctx context.Context, sock socket.Socket) socket.Result {
	url := sock.Host + ":" + strconv.Itoa(sock.Port) + sock.PathHTTP

	if runner.verbose {
		log.Println("httprunner: connect:", url)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return socket.Result{Socket: sock, Passed: false, Error: err}
	}
	req.Header.Set("User-Agent", fmt.Sprintf("dish/%s", agentVersion))

	resp, err := runner.client.Do(req)
	if err != nil {
		return socket.Result{Socket: sock, Passed: false, Error: err}
	}
	defer resp.Body.Close()

	if !slices.Contains(sock.ExpectedHTTPCodes, resp.StatusCode) {
		err = fmt.Errorf("expected codes: %v, got %d", sock.ExpectedHTTPCodes, resp.StatusCode)
	}

	return socket.Result{
		Socket:       sock,
		Passed:       slices.Contains(sock.ExpectedHTTPCodes, resp.StatusCode),
		ResponseCode: resp.StatusCode,
		Error:        err,
	}
}
