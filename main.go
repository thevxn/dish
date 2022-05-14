// go:build ignore

// savla-dish executable -- providing a simple remote socket testing
package main

import (
	"fmt"
	"regexp"

	"savla-dish/messenger"
	"savla-dish/runner"
	//"savla-dish/telnet"
	"savla-dish/zasuvka"
)

const (
	// if false, messenger sends telegrams
	DevMode = false
	// could be a character/type byte too maybe
	newLine string = "%0A"
	//socketListSource string = "./demo_sockets.json"
	socketListSource string = "http://swapi.savla.su:80/dish/sockets"
)

func main() {
	// load init config/socket list to run tests on --- external file!
	sockets := zasuvka.GibPole(socketListSource)

	// final report header
	msgText := fmt.Sprintf("savla-dish run results (failed): %s", newLine)
	failedCount := 0

	// iterate over given/loaded sockets
	for _, socket := range sockets.Sockets {
		host := socket.Host
		port := socket.Port
		
		// http/https app protocol patterns check
		match, _ := regexp.MatchString("^(http|https)://", host); if match {
			// compare HTTP response codes 
			expectedCodes := socket.ExpectedHttpCodes
			path := socket.PathHttp

			// here, 'status' should contain HTTP code if >0
			status := runner.CheckSite(host, port, path, expectedCodes); if status != 0 {
				msgText += fmt.Sprintf("%s:%d %d %s", host, port, status, newLine)
				failedCount++
			}
			continue
		}

		// testing raw host and port (tcp), report only unsuccessful tests
		status, _ := runner.RawConnect("tcp", host, port); if status > 0 {
			msgText += fmt.Sprintf("%s:%d %s %s", host, port, "timeout", newLine)
			failedCount++
			//msgText += fmt.Sprintf("%s:%d %d %s", host, port, status, newLine)
		}
	}

	// mute dish messenger if needed in a custom build/env
	if !DevMode && failedCount > 0 {
		messenger.SendMsg(msgText)
		//message.Send()
	}

	// final report output to stdout/console/logs
	fmt.Printf(msgText)

	//fmt.Println( resp.Status )
	//fmt.Println( resp.StatusCode )
	//fmt.Println( resp.Proto )
}
