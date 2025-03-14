package main

import (
	"fmt"
	"log"
	"sync"

	"go.vxn.dev/dish/pkg/alert"
	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/message"
	"go.vxn.dev/dish/pkg/netrunner"
	"go.vxn.dev/dish/pkg/socket"
)

func main() {
	// Load socket list to run tests on --- external file!
	list := socket.FetchSocketList(config.Source)

	var (
		messengerText string
		resultsToPush = message.Results{Map: make(map[string]bool)}
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

		go netrunner.TestSocket(sock, channels[i], &wg)
		i++
	}

	// Merge channels into one
	results := fanInChannels(channels...)
	wg.Wait()

	// Collect results
	for result := range results {
		if !result.Passed || result.Error != nil {
			failedCount++
			messengerText += formatMessengerText(result)
		}
		resultsToPush.Map[result.Socket.ID] = result.Passed
	}

	handlePushgateway(failedCount)
	handleStateUpdate(resultsToPush)

	if failedCount > 0 {
		handleAlerts(messengerText, resultsToPush)
		log.Println(messengerText)
		return
	}

	log.Println("dish run: all tests ok")
}

// Fan-in function that collects results from multiple workers
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

func formatMessengerText(result socket.Result) string {
	// Hotfix unsupported <nil> tag by TG
	if result.Error == nil {
		result.Error = fmt.Errorf("")
	}

	if result.Socket.PathHTTP != "" {
		return fmt.Sprintf("• %s:%d%s (code %d) -- %v\n",
			result.Socket.Host, result.Socket.Port, result.Socket.PathHTTP, result.ResponseCode, result.Error)
	}
	return fmt.Sprintf("• %s:%d -- %v\n", result.Socket.Host, result.Socket.Port, result.Error)
}

func handlePushgateway(failedCount int) {
	if config.UsePushgateway && config.TargetURL != "" {
		msg := message.Make(failedCount)
		if err := msg.PushDishResults(); err != nil {
			log.Printf("Failed to push dish results: %v", err)
		}
	}
}

func handleStateUpdate(results message.Results) {
	if config.UpdateStates {
		if err := message.UpdateSocketStates(results); err != nil {
			log.Printf("Failed to update socket states: %v", err)
		}
	}
}

func handleAlerts(messengerText string, results message.Results) {
	if config.UseTelegram {
		err := alert.SendTelegram(messengerText)
		if err != nil {
			log.Printf("Error sending Telegram notification: %v", err)
		}
	}

	if config.UseWebhooks {
		err := alert.SendWebhooks(&results)
		if err != nil {
			log.Printf("Error sending webhook notification: %v", err)
		}
	}
}
