package main

import (
	"fmt"
	"sync"

	"go.vxn.dev/dish/pkg/alert"
	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/netrunner"
	"go.vxn.dev/dish/pkg/socket"
)

// testResults holds the overall results of all socket checks combined.
type testResults struct {
	messengerText string
	results       *alert.Results
	failedCount   int
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

// runTests orchestrates the process of checking of a list of sockets. It fetches the socket list, runs socket checks, collects results and returns them.
func runTests(cfg *config.Config) (*testResults, error) {
	// Load socket list to run tests on
	list, err := socket.FetchSocketList(cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading socket list: %w", err)
	}

	// Print loaded sockets if flag is set in cfg
	if cfg.Verbose {
		socket.PrintSockets(list)
	}

	testResults := &testResults{
		messengerText: "",
		results:       &alert.Results{Map: make(map[string]bool)},
		failedCount:   0,
	}

	var (
		// A slice of channels needs to be used here so that each goroutine has its own channel which it then closes upon performing the socket check. One shared channel for all goroutines would not work as it would not be clear which goroutine should close the channel.
		channels = make([](chan socket.Result), len(list.Sockets))

		wg sync.WaitGroup
		i  int
	)

	// Start goroutines for each socket test
	for _, sock := range list.Sockets {
		wg.Add(1)
		channels[i] = make(chan socket.Result)

		go netrunner.RunSocketTest(sock, channels[i], &wg, cfg)
		i++
	}

	// Merge channels into one
	results := fanInChannels(channels...)
	wg.Wait()

	// Collect results
	for result := range results {
		if !result.Passed || result.Error != nil {
			testResults.failedCount++
		}
		if !result.Passed || cfg.TextNotifySuccess {
			testResults.messengerText += alert.FormatMessengerText(result)
		}
		testResults.results.Map[result.Socket.ID] = result.Passed
	}

	return testResults, nil
}
