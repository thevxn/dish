package runner

import (
	"log"
	"net"
	"net/http"
	"strconv"
	//"time"
)


func RawConnect(host string, port int) (status int, err error) {
	// reformat int port to number string
	portString := strconv.Itoa(port)

	// vars
	endpoint := net.JoinHostPort(host, portString)
	//endpoint := host + ":" + portString
	//protocol := "tcp"
	//timeout := time.Second

	log.Println("runner: " + endpoint)

	// open the socket
	//conn, err := net.DialTimeout("tcp", endpoint, timeout)
	conn, err := net.Dial("tcp", endpoint)

	// Close open conn after 5 seconds
	//&conn.SetReadDeadline(time.Seconds(&timeout))

	// prolly more possible to get not-nil err, than not-nil conn 
	// https://stackoverflow.com/a/56336811
	if err != nil {
		log.Println("runner: conn error: " + endpoint)
		log.Println(err)
		return 1, err
	}
	if conn != nil {
		defer conn.Close()
		//log.Println("[-] Opened: " + conn.Read())
		return 0, nil
	}

	return 0, nil
}

func CheckSite(endpoint string, port int) (status int) {
	var netClient = &http.Client{}

	resp, err := netClient.Get("http://" + endpoint + "" + strconv.Itoa(port))

	if err != nil {
		log.Fatalln(err)
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
