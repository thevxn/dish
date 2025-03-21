package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"go.vxn.dev/dish/pkg/alert"
	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/netrunner"
	"go.vxn.dev/dish/pkg/socket"
)

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
	config, err := config.NewConfig(flag.CommandLine, os.Args[1:])
	if err != nil {
		log.Fatal("error loading config: ", err)
	}

	// Load socket list to run tests on --- external file!
	list := socket.FetchSocketList(config.Source, config.ApiHeaderName, config.ApiHeaderValue, config.Verbose)

	var (
		messengerText string
		resultsToPush = alert.Results{Map: make(map[string]bool)}
		failedCount   int
		// A slice of channels needs to be used here so that each goroutine has its own channel which it then closes upon performing the socket check. One shared channel for all goroutines would not work as it would not be clear which goroutine should close the channel.
		channels = make([](chan socket.Result), len(list.Sockets))
		wg       sync.WaitGroup
		i        int
	)

	// Start goroutines for each socket test
	for _, sock := range list.Sockets {
		wg.Add(1)
		channels[i] = make(chan socket.Result)

		go netrunner.TestSocket(sock, channels[i], &wg, config.TimeoutSeconds, config.Verbose)
		i++
	}

	// Merge channels into one
	results := fanInChannels(channels...)
	wg.Wait()

	// Collect results
	for result := range results {
		if !result.Passed || result.Error != nil {
			failedCount++
			messengerText += alert.FormatMessengerText(result)
		}
		resultsToPush.Map[result.Socket.ID] = result.Passed
	}

	alert.HandleAlerts(messengerText, resultsToPush, failedCount, config)

	if failedCount > 0 {
		log.Println(messengerText)
		return
	}

	log.Println("dish run: all tests ok")
}
