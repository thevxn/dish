package app

import (
	"fmt"
	"log"
	"regexp"
	"savla-dish/cmd/messenger"
	"savla-dish/cmd/runner"
	"savla-dish/pkg/config"
	"savla-dish/pkg/reporter"
	"savla-dish/pkg/socket"
)

func Run() {

	// load init config/socket list to run tests on --- external file!
	list := socket.FetchSocketList(config.Source)

	// final report header
	messengerText := fmt.Sprintln("[ savla-dish run results (failed) ]")
	var failedCount int8 = 0

	// iterate over given/loaded sockets
	for _, socket := range list.Sockets {
		// http/https app protocol patterns check
		_, err := regexp.Compile("^(http|https)://" + socket.Host)
		if err == nil {
			// here, 'status' should contain HTTP code if >0
			status := runner.CheckSite(socket)
			if status != 0 {
				messengerText += fmt.Sprintln(socket.Host, ":", socket.Port, socket.PathHTTP, "--", status)
				socket.Results.SocketReached = false
				failedCount++
			}
			continue
		}

		// testing raw host and port (tcp), report only unsuccessful tests; exclusively non-HTTP/S sockets
		status, _ := runner.RawConnect(socket)
		if status > 0 {
			messengerText += fmt.Sprintln(socket.Host, ":", socket.Port, "-- timeout")
			socket.Results.SocketReached = false
			failedCount++
		}
	}

	// report failedCount to pushgateway
	if config.UsePushgateway && config.TargetURL != "" {
		reporter.Reporter.FailedCount = failedCount
		reporter.PushDishResults()
	}

	// messenger threshold
	if failedCount > 0 {
		// send alert message
		if config.UseTelegram {
			messenger.SendTelegram(messengerText)
		}

		// final report output to stdout/console/docker logs
		log.Println(messengerText)
		return
	}

	log.Println("dish run: all tests ok")
}
