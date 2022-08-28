package config

import (
	"flag"
)

var (
	Source           string
	Verbose          bool
	TargetURL        string
	UsePushgateway   bool
	UseTelegram      bool
	TelegramBotToken string
	TelegramChatID   string
	Timeout          int // In seconds
)

// Gets called before main()
func init() {

	sourceFlag := flag.String("source", "demo_sockets.json", "a string, path to/URL JSON socket list")
	verboseFlag := flag.Bool("verbose", true, "a bool, console stdout logging toggle")
	timeoutFlag := flag.Int("timeout", 10, "a int, timeout in seconds for http and tcp calls")
	targetURLFlag := flag.String("target", "", "a string, result update path/URL, plaintext/byte output")
	usePushgatewayFlag := flag.Bool("pushgw", false, "a bool, enable reporter module to post dish results to pushgateway")

	// telegram provider flags
	useTelegramFlag := flag.Bool("telegram", false, "a bool, Telegram provider usage toggle")
	telegramBotTokenFlag := flag.String("telegramBotToken", "", "a string, Telegram bot private token")
	telegramChatIDFlag := flag.String("telegramChatID", "", "a string/signet int, Telegram chat/channel ID")

	flag.Parse()

	Source = *sourceFlag
	Verbose = *verboseFlag
	TargetURL = *targetURLFlag
	UsePushgateway = *usePushgatewayFlag
	UseTelegram = *useTelegramFlag
	TelegramBotToken = *telegramBotTokenFlag
	TelegramChatID = *telegramChatIDFlag
	Timeout = *timeoutFlag

}
