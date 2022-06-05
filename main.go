package main

import (
	"flag"
	"fmt"
	"regexp"

	"savla-dish/messenger"
	"savla-dish/runner"
	"savla-dish/socket"
)

func main() {
	// predefine flags --- flag returns a pointer!
	sourceFlag := flag.String("source", "demo_sockets.json", "a string, path to/URL JSON socket list")
	verboseFlag := flag.Bool("verbose", false, "a bool, console stdout logging toggle")

	// telegram provider flags
	messenger.UseTelegram = flag.Bool("telegram", false, "a bool, Telegram provider usage toggle")
	messenger.TelegramBotToken = flag.String("telegramBotToken", "", "a string, Telegram bot private token")
	messenger.TelegramChatID = flag.String("telegramChatID", "", "a string/signet int, Telegram chat/channel ID")

	flag.Parse()

	// load init config/socket list to run tests on --- external file!
	sockets := socket.FetchSocketList(*sourceFlag, *verboseFlag)

	// final report header
	messengerText := fmt.Sprintln("[ savla-dish run results (failed) ]")
	var failedCount int8 = 0

	// iterate over given/loaded sockets
	for _, socket := range sockets.Sockets {
		// http/https app protocol patterns check
		match, _ := regexp.MatchString("^(http|https)://", socket.Host)
		if match {
			// here, 'status' should contain HTTP code if >0
			status := runner.CheckSite(socket, *verboseFlag)
			if status != 0 {
				messengerText += fmt.Sprintln(socket.Host, ":", socket.Port, socket.PathHttp, "--", status)
				failedCount++
			}
			continue
		}

		// testing raw host and port (tcp), report only unsuccessful tests; exclusively non-HTTP/S sockets
		status, _ := runner.RawConnect(socket, *verboseFlag)
		if status > 0 {
			messengerText += fmt.Sprintln(socket.Host, ":", socket.Port, "-- timeout")
			failedCount++
		}
	}

	// reporting threshold
	if failedCount > 0 {
		// send alert message
		if *messenger.UseTelegram {
			messenger.SendMsg(messengerText, *verboseFlag)
		}

		// final report output to stdout/console/docker logs
		fmt.Printf(messengerText)
		return
	}

	fmt.Println("dish run: all tests ok")
}
