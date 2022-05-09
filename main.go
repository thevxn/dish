// savla-dish executable -- providing a simple remote socket testing
package main

import (
	"fmt"
	"savla-dish/messenger"
	"savla-dish/runner"
	"savla-dish/zasuvka"
)

func main() {
	endpointFile := "demo_sockets.json"

	sockets := zasuvka.GibPole(endpointFile, false)

	msgText := ""
	newLine := "%0A"

	for i := 0; i < len(sockets.Sockets); i++ {
		e := sockets.Sockets[i].Endpoint
		p := sockets.Sockets[i].Port

		status := runner.CheckSite(e, p)

		msgText += fmt.Sprintf("%s %d %d %s", e, p, status, newLine)

		fmt.Println(e, p, status)
	}

	messenger.SendMsg(msgText)

	//fmt.Println( resp.Status )
	//fmt.Println( resp.StatusCode )
	//fmt.Println( resp.Proto )
}
