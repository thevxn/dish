// savla-dish executable -- providing a simple remote socket testing
package main

import (
	"fmt"
	"savla-dish/runner"
	"savla-dish/zasuvka"
)

func main() {
	endpointFile := "demo_sockets.json"
	//endpoint := "http://docs.savla.su/kus/hovna"

	sockets := zasuvka.GibPole(endpointFile, false)

	for i := 0; i < len(sockets.Sockets); i++ {
		e := sockets.Sockets[i].Endpoint
		p := sockets.Sockets[i].Port

		status := runner.Run(e, p)
		fmt.Println(e, p, status)
	}

	//fmt.Println( resp.Status )
	//fmt.Println( resp.StatusCode )
	//fmt.Println( resp.Proto )
}
