package app

import (
	"fmt"
	"log"
	"regexp"

	"dish/pkg/alert"
	"dish/pkg/config"
	"dish/pkg/message"
	"dish/pkg/netrunner"
	"dish/pkg/socket"
)

func Run() {

	// load socket list to run tests on --- external file!
	list := socket.FetchSocketList(config.Source)

	messengerText := "[ dish run results (failed) ]\n"
	results := message.Results{Map: make(map[string]bool)}
	failedCount := 0

	regex, err := regexp.Compile("^(http|https)://")
	if err != nil {
		log.Println("Failed to create new regex object")
		return
	}

	// iterate over given/loaded sockets
	for _, socket := range list.Sockets {
		// http/https app protocol patterns check
		match := regex.MatchString(socket.Host)
		if !match {
			// testing raw host and port (tcp), report only unsuccessful tests; exclusively non-HTTP/S sockets
			rawErr := netrunner.RawConnect(socket)
			if rawErr != nil {
				messengerText += fmt.Sprintln(socket.Host, ":", socket.Port, rawErr.Error())
				failedCount++
			}
			results.Map[socket.Name] = (rawErr == nil)
			continue
		}

		ok, responseCode, httpErr := netrunner.CheckSite(socket)
		if !ok {
			failedCount++
			if httpErr != nil {
				messengerText += fmt.Sprintln(socket.Host, ":", socket.Port, socket.PathHTTP, "--", httpErr.Error())
			}
			messengerText += fmt.Sprintln(socket.Host, ":", socket.Port, socket.PathHTTP, "--", "Did not match expected response codes: got ", responseCode)
			results.Map[socket.Name] = (httpErr == nil)
		}
	}

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
