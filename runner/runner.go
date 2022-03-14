package runner

import (
	/*"fmt"*/
	"log"
	"net/http"
)

func Run(endpoint string, port int) (status int) {
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		// status = err
		// continue 
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

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
