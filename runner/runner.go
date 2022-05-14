// +build dev

// Runner package for custom socket test execution
package runner

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

const DevMode = true

// getEndpoint private (unexported) macro
func getEndpoint(host string, port int) string {
	// reformat int port to number string
	portString := strconv.Itoa(port)

	return net.JoinHostPort(host, portString)
}

// RawConnect function to direct host:port socket check 
func RawConnect(protocol string, host string, port int) (status int, err error) {
	endpoint := getEndpoint(host, port)
	timeout := time.Duration(5 * time.Second)

	// console debug
	if DevMode {
		log.Println("runner: rawconnect: " + endpoint)
	}

	// open the socket
	//conn, err := net.DialTimeout("tcp", endpoint, timeout)
	conn, err := net.DialTimeout(protocol, endpoint, timeout)

	// close open conn after 5 seconds
	//conn.SetReadDeadline(time.Second*5)

	// prolly more possible to get not-nil err, than not-nil conn 
	// see https://stackoverflow.com/a/56336811
	if err != nil {
		log.Println("runner: rawconnect: conn error: " + endpoint)
		log.Println(err)
		return 1, err
	}
	if conn != nil {
		conn.Close()
		return 0, nil
	}

	// unexpected error
	return 1, nil
}

// checkHttpCode function for response and expected HTTP codes comparsion
func checkHttpCode(responseCode int, expectedCodes []int) (status int) {
	for _, code := range expectedCodes {
		if responseCode == code {
			// site is OK! do not report ok sites?
			//&msgText += fmt.Sprintf("%s:%d %d %s", h, p, status, newLine)
			return 0
			break
		}
	}

	return 1
}

// CheckSite executes test over HTTP/S endpoints exclusively
func CheckSite(host string, port int, path string, expectedCodes []int) (status int) {
	// config http client
	var netClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	url := host + ":" + strconv.Itoa(port) + path

	// console debug
	if DevMode {
		log.Println("runner: checksite: " + url)
	}

	// open socket --- give Head
	resp, err := netClient.Head(url)
	if err != nil {
		// this construct prolly halts the main process, close socket, not the whole executable...
		//log.Fatalln(err)
		log.Println(err)
		return 1
	}

	// fetch StatusCode for HTTP expected code comparsion
	if resp != nil {
		//defer resp.Body.Close()
		log.Print(resp.StatusCode)
		//return resp.StatusCode
		return checkHttpCode(resp.StatusCode, expectedCodes)
	}

	return 1


	//
	// LEGACY, to be deleted
	//

	/*
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		log.Fatalln("lmao")
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalln("lol")
	}

	defer resp.Body.Close()

	/*switch resp.StatusCode {
		case 200:
			fmt.Println("ok, 200")
			break
		case 404:
			fmt.Println("nok, 404")
			break
		default:
			fmt.Println(resp.Status)
	}*/

	//return resp.StatusCode
}
