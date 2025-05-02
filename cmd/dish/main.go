package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"go.vxn.dev/dish/pkg/alert"
	"go.vxn.dev/dish/pkg/config"
)

func main() {
	cfg, err := config.NewConfig(flag.CommandLine, os.Args[1:])
	if err != nil {
		// If the error is caused due to no source being provided, print help
		if errors.Is(err, config.ErrNoSourceProvided) {
			printHelp()
			os.Exit(1)
		}
		// Otherwise, print the error
		log.Print("error loading config: ", err)
		return
	}

	log.Println("dish run: started")

	// Run tests on sockets
	res, err := runTests(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	// Submit results and alerts
	alert.HandleAlerts(res.messengerText, res.results, res.failedCount, cfg)

	if res.failedCount > 0 {
		log.Println("dish run: some tests failed:\n", res.messengerText)
		return
	}

	log.Println("dish run: all tests ok")
}
