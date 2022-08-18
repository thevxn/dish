package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"savla-dish/messenger"
	"savla-dish/reporter"
	"savla-dish/runner"
	"savla-dish/socket"
	//"github.com/savla-dev/savla-dish/messenger"
	//"github.com/savla-dev/savla-dish/reporter"
	//"github.com/savla-dev/savla-dish/runner"
	//"github.com/savla-dev/savla-dish/socket"
)

func main() {
	// predefine flags --- flag returns a pointer! TODO: tidy this...
	sourceFlag := flag.String("source", "demo_sockets.json", "a string, path to/URL JSON socket list")
	runner.Verbose = flag.Bool("verbose", false, "a bool, console stdout logging toggle")
	socket.Verbose = runner.Verbose
	messenger.Verbose = runner.Verbose

	// reporter driver flags
	reporter.TargetURL = flag.String("target", "", "a string, result update path/URL, plaintext/byte output")
	reporter.UsePushgateway = flag.Bool("pushgw", false, "a bool, enable reporter module to post dish results to pushgateway")

	// telegram provider flags
	messenger.UseTelegram = flag.Bool("telegram", false, "a bool, Telegram provider usage toggle")
	messenger.TelegramBotToken = flag.String("telegramBotToken", "", "a string, Telegram bot private token")
	messenger.TelegramChatID = flag.String("telegramChatID", "", "a string/signet int, Telegram chat/channel ID")

	flag.Parse()

	// load init config/socket list to run tests on --- external file!
	sockets := socket.FetchSocketList(*sourceFlag)

	// final report header
	messengerText := fmt.Sprintln("[ savla-dish run results (failed) ]")
	var failedCount int8 = 0
	reporter.Reporter.FailedCount = 0

	// iterate over given/loaded sockets
	for _, socket := range sockets.Sockets {
		// http/https app protocol patterns check
		match, _ := regexp.MatchString("^(http|https)://", socket.Host)
		if match {
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
	if *reporter.UsePushgateway && *reporter.TargetURL != "" {
		reporter.Reporter.FailedCount = failedCount
		reporter.PushDishResults()
	}

	// messenger threshold
	if failedCount > 0 {
		// send alert message
		if *messenger.UseTelegram {
			messenger.SendTelegram(messengerText)
		}

		// final report output to stdout/console/docker logs
		log.Printf(messengerText)
		return
	}

	log.Println("dish run: all tests ok")
}
