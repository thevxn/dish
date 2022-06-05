//go:build dev

package telnet

import (
	/*"fmt"*/
	"log"

	"github.com/reiver/go-telnet"
)

func TestDial(endpoint string, port int) (status int) {
	//var caller telnet.Caller = telnet.StandardCaller

	conn, err := telnet.DialTo(endpoint + ":" + string(port))

	if err != nil {
		log.Print(err)
		return 1
	}

	var data []byte
	log.Print(conn.Read(data))
	conn.Close()

	//return resp.StatusCode
	return 0
}
