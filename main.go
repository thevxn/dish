// savla-dish executable -- providing a simple remote socket testing
package main

import (
	"fmt"

	"savla-dish/messenger"
	"savla-dish/runner"
	//"savla-dish/telnet"
	"savla-dish/zasuvka"
)

const (
	devMode = true
	// could be a character/type byte too maybe
	newLine string = "%0A"
	socketsListFile string = "demo_sockets.json"
)

func main() {
	// load init config/socket list to run tests on
	sockets := zasuvka.GibPole(socketsListFile, false)

	// header
	msgText := fmt.Sprintf("savla-dish run results: %s", newLine)

	// iterate over given/loaded sockets
	for i := 0; i < len(sockets.Sockets); i++ {
	//for _, sock := range sockets.Sockets {
		//h := sock.Endpoint
		//p := sock.Port
		h := sockets.Sockets[i].Host
		p := sockets.Sockets[i].Port
		//e := net.JoinHostPort(h, p)

		// testing paradigmas 
		//status := runner.CheckSite(h, p)
		//status := telnet.TestDial(h, p)
		status, _ := runner.RawConnect(h, p)

		msgText += fmt.Sprintf("%s %d %d %s", h, p, status, newLine)

		fmt.Println(h, p, status)
	}

	// mute dish messenger if needed in a custom build/env
	if !devMode {
		messenger.SendMsg(msgText)
	}

	// final report output to stdout/console/logs
	fmt.Printf(msgText)

	//fmt.Println( resp.Status )
	//fmt.Println( resp.StatusCode )
	//fmt.Println( resp.Proto )
}
