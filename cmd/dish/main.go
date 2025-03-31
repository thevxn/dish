package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"go.vxn.dev/dish/pkg/alert"
	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/netrunner"
	"go.vxn.dev/dish/pkg/socket"
)

func printHelp() {
	fmt.Print("Usage: dish [FLAGS] SOURCE\n\n")
	fmt.Print("A lightweight, one-shot socket checker\n\n")
	fmt.Println("SOURCE must be a file path leading to a JSON file with a list of sockets to be checked or a URL leading to a remote JSON API from which the list of sockets can be retrieved")
	fmt.Println("Use the `-h` flag for a list of available flags")
}

// fanInChannels collects results from multiple goroutines
func fanInChannels(channels ...chan socket.Result) <-chan socket.Result {
	var wg sync.WaitGroup
	out := make(chan socket.Result)

	// Start a goroutine for each channel
	for _, channel := range channels {
		wg.Add(1)
		go func(ch <-chan socket.Result) {
			defer wg.Done()
			for result := range ch {
				// Forward the result to the output channel
				out <- result
			}
		}(channel)
	}

	// Close the output channel once all workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	cfg, err := config.NewConfig(flag.CommandLine, os.Args[1:])
	if err != nil {
		// If the error is caused due to no source being provided, print help
		if errors.Is(err, config.ErrNoSourceProvided) {
			printHelp()
			os.Exit(1)
		}
		// Otherwise, print the errror
		log.Print("error loading config: ", err)
		return
	}

	log.Println("dish run: started")

	// Load socket list to run tests on
	list, err := socket.FetchSocketList(cfg.Source, cfg.ApiHeaderName, cfg.ApiHeaderValue, cfg.Verbose)
	if err != nil {
		log.Print("error loading socket list: ", err)
		return
	}

	var (
		messengerText string
		resultsToPush = alert.Results{Map: make(map[string]bool)}
		failedCount   int

		// A slice of channels needs to be used here so that each goroutine has its own channel which it then closes upon performing the socket check. One shared channel for all goroutines would not work as it would not be clear which goroutine should close the channel.
		channels = make([](chan socket.Result), len(list.Sockets))

		wg sync.WaitGroup
		i  int
	)

	// Start goroutines for each socket test
	for _, sock := range list.Sockets {
		wg.Add(1)
		channels[i] = make(chan socket.Result)

		go netrunner.TestSocket(sock, channels[i], &wg, cfg.TimeoutSeconds, cfg.Verbose)
		i++
	}

	// Merge channels into one
	results := fanInChannels(channels...)
	wg.Wait()

	// Collect results
	for result := range results {
		if !result.Passed || result.Error != nil {
			failedCount++
		}
		if !result.Passed || cfg.TextNotifySuccess {
			messengerText += alert.FormatMessengerText(result)
		}
		resultsToPush.Map[result.Socket.ID] = result.Passed
	}

	alert.HandleAlerts(messengerText, resultsToPush, failedCount, cfg)

	if failedCount > 0 {
		log.Println("dish run: some tests failed")
		log.Print("\n", messengerText)
		return
	}

	log.Println("dish run: all tests ok")
}
