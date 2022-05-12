package runner

import (
	"log"
	"net"
	"net/http"
	"strconv"
	//"time"
)

// getEndpoint private (unexported) macro
func getEndpoint(host string, port int) string {
	// reformat int port to number string
	portString := strconv.Itoa(port)

	return net.JoinHostPort(host, portString)
}


func RawConnect(protocol string, host string, port int) (status int, err error) {
	// vars
	endpoint := getEndpoint(host, port)
	//timeout := time.Second

	// console debug
	log.Println("runner: rawconnect: " + endpoint)

	// open the socket
	//conn, err := net.DialTimeout("tcp", endpoint, timeout)
	conn, err := net.Dial(protocol, endpoint)

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
		defer conn.Close()
		//log.Println("[-] Opened: " + conn.Read())
		return 0, nil
	}

	// unexpected error
	return 1, nil
}

// CheckSite executes test over HTTP/S endpoints exclusively
func CheckSite(host string, port int) (status int) {
	var netClient = &http.Client{}
	url := getEndpoint(host, port)

	// console debug (should be toggable to increase speed)
	log.Println("runner: checksite: " + url)

	// open socket
	resp, err := netClient.Get(url)
	if err != nil {
		log.Fatalln(err)
		return 1
	}

	log.Print(resp)
	return 0

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
