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

	for i := 0; i < len(sockets.Sockets); i++ {
		e := sockets.Sockets[i].Endpoint
		p := sockets.Sockets[i].Port

		status := runner.CheckSite(e, p)
		fmt.Println(e, p, status)
	}

	messenger.SendMsg("testikl")

	//fmt.Println( resp.Status )
	//fmt.Println( resp.StatusCode )
	//fmt.Println( resp.Proto )
}
