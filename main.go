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
	socketsListFile string = "demo_sockets.json"
)

func main() {
	// load init config/socket list to run tests on --- external file!
	sockets := zasuvka.GibPole(socketsListFile)

	// final report header
	msgText := fmt.Sprintf("savla-dish run results (failed): %s", newLine)

	// iterate over given/loaded sockets
	for i := 0; i < len(sockets.Sockets); i++ {
		h := sockets.Sockets[i].Host
		p := sockets.Sockets[i].Port

		// http/https app protocol patterns check
		match, _ := regexp.MatchString("^(http|https)://", h); if match {
			// compare HTTP response codes 
			exp := sockets.Sockets[i].ExpectedHttpCodes
			status := runner.CheckSite(h, p, exp); if status != 0 {
				msgText += fmt.Sprintf("%s:%d %d %s", h, p, status, newLine)
			}
			continue
		}

		// testing raw host and port (tcp), report only unsuccessful connects?
		status, _ := runner.RawConnect("tcp", h, p); if status > 0 {
			msgText += fmt.Sprintf("%s:%d %d %s", h, p, status, newLine)
		}

		//fmt.Println(h, p, status)
	}

	// mute dish messenger if needed in a custom build/env
	if !DevMode {
		messenger.SendMsg(msgText)
	}

	// final report output to stdout/console/logs
	fmt.Printf(msgText)

	//fmt.Println( resp.Status )
	//fmt.Println( resp.StatusCode )
	//fmt.Println( resp.Proto )
}
