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
	DevMode = true
	// could be a character/type byte too maybe
	newLine string = "%0A"
	socketListFile string = "demo_sockets.json"
)

func main() {
	// load init config/socket list to run tests on --- external file!
	sockets := zasuvka.GibPole(socketListFile)

	// final report header
	msgText := fmt.Sprintf("savla-dish run results (failed): %s", newLine)

	// iterate over given/loaded sockets
	for _, socket := range sockets.Sockets {
		host := socket.Host
		port := socket.Port
		
		// http/https app protocol patterns check
		match, _ := regexp.MatchString("^(http|https)://", host); if match {
			// compare HTTP response codes 
			expectedCodes := socket.ExpectedHttpCodes
			path := socket.PathHttp

			status := runner.CheckSite(host, port, path, expectedCodes); if status != 0 {
				msgText += fmt.Sprintf("%s:%d %d %s", host, port, status, newLine)
			}
			continue
		}

		// testing raw host and port (tcp), report only unsuccessful connects?
		status, _ := runner.RawConnect("tcp", host, port); if status > 0 {
			msgText += fmt.Sprintf("%s:%d %d %s", host, port, status, newLine)
		}
	}

	// mute dish messenger if needed in a custom build/env
	if !DevMode {
		messenger.SendMsg(msgText)
		//message.Send()
	}

	// final report output to stdout/console/logs
	fmt.Printf(msgText)

	//fmt.Println( resp.Status )
	//fmt.Println( resp.StatusCode )
	//fmt.Println( resp.Proto )
}
