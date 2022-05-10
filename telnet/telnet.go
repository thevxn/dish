package telnet

import (
	/*"fmt"*/
	"log"
	"net/http"

	"github.com/reiver/go-telnet"
)

func CheckSite(endpoint string, port int) (status int) {

	var caller telnet.Caller = telnet.StandardCaller

	telnet.DialToAndCall(endpoint + ":" + string(port), caller)
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		// status = err
		// continue
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
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

	return resp.StatusCode
}
