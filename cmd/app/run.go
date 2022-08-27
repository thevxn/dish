package app

import (
	"fmt"
	"log"
	"regexp"
	"savla-dish/pkg/alert"
	"savla-dish/pkg/config"
	"savla-dish/pkg/netrunner"
	"savla-dish/pkg/reporter"
	"savla-dish/pkg/socket"
)

func Run() {

	// load init config/socket list to run tests on --- external file!
	list := socket.FetchSocketList(config.Source)

	// final report header
	messengerText := "[ savla-dish run results (failed) ]"
	failedCount := 0

	// iterate over given/loaded sockets
	for _, socket := range list.Sockets {
		// http/https app protocol patterns check
		_, err := regexp.Compile("^(http|https)://" + socket.Host)
		if err == nil {
			// here, 'ok' should contain HTTP code if >0
			ok := netrunner.CheckSite(&socket)
			if !ok {
				messengerText += fmt.Sprintln(socket.Host, ":", socket.Port, socket.PathHTTP, "--", ok)
				failedCount++
			}
		}

		// testing raw host and port (tcp), report only unsuccessful tests; exclusively non-HTTP/S sockets
		err = netrunner.RawConnect(&socket)
		if err != nil {
			messengerText += fmt.Sprintln(socket.Host, ":", socket.Port, "-- timeout")
			failedCount++
			panic(err)
		}
	}

	// report failedCount to pushgateway
	if config.UsePushgateway && config.TargetURL != "" {
		msg := reporter.MakeMessage(failedCount)
		err := msg.PushDishResults()
		if err != nil {
			panic(err)
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
