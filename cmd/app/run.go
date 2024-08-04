package app

import (
	"fmt"
	"log"
	"sync"

	"dish/pkg/alert"
	"dish/pkg/config"
	"dish/pkg/message"
	"dish/pkg/netrunner"
	"dish/pkg/socket"
)

func Run() {
	// load socket list to run tests on --- external file!
	list := socket.FetchSocketList(config.Source)

	//messengerText := "[ dish run results (failed) ]\n"
	messengerText := ""
	results := message.Results{Map: make(map[string]bool)}
	failedCount := 0

	socketChan := make(chan socket.Result)
	var wg sync.WaitGroup

	// iterate over given/loaded sockets --> start goroutines
	for _, socket := range list.Sockets {
		wg.Add(1)
		go netrunner.TestSocket(socket, socketChan, &wg)
	}

	// iterate again to fetch results from the channel
	for _, _ = range list.Sockets {
		result := <-socketChan

		if !result.Passed || result.Error != nil {
			failedCount++
			if result.Error != nil {
				if result.Socket.PathHTTP != "" {
					messengerText += fmt.Sprintln(result.Socket.Host, ":", result.Socket.Port, result.Socket.PathHTTP, " (code ", result.ResponseCode, " )", "--", result.Error.Error())
				} else {
					messengerText += fmt.Sprintln(result.Socket.Host, ":", result.Socket.Port, result.Error.Error())
				}
			}
		}
		results.Map[result.Socket.ID] = (result.Error == nil)
	}

	wg.Wait()
	close(socketChan)

	// report failedCount to pushgateway
	if config.UsePushgateway && config.TargetURL != "" {
		msg := message.Make(failedCount)
		pushErr := msg.PushDishResults()
		if pushErr != nil {
			log.Println("Failed to push dish results, err: " + pushErr.Error())
		}
	}

	if config.UpdateStates {
		updateErr := message.UpdateSocketStates(results)
		if updateErr != nil {
			log.Println("Failed to update socket states, err: " + updateErr.Error())
		}
	}

	if failedCount > 0 {
		// send alert message
		if config.UseTelegram {
			alert.SendTelegram(messengerText)
		}

		// final report output to stdout/console/docker logs
		log.Println(messengerText)
		return
	}

	log.Println("dish run: all tests ok")
}
