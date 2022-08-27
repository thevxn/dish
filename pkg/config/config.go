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
)

// var Config Flags

// type Flags struct {
// 	Source           string
// 	Verbose          bool
// 	TargetURL        string
// 	UsePushgateway   bool
// 	UseTelegram      bool
// 	TelegramBotToken string
// 	TelegramChatID   string
// }

// Gets called before main()
func init() {

	sourceFlag := flag.String("source", "/home/geeko/savla/savla-dish/demo_sockets.json", "a string, path to/URL JSON socket list")
	verboseFlag := flag.Bool("verbose", true, "a bool, console stdout logging toggle")

	// reporter driver flags
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

	// Config = Flags{
	// 	Source:           *sourceFlag,
	// 	Verbose:          *verboseFlag,
	// 	TargetURL:        *targetURLFlag,
	// 	UsePushgateway:   *usePushgatewayFlag,
	// 	UseTelegram:      *useTelegramFlag,
	// 	TelegramBotToken: *telegramBotTokenFlag,
	// 	TelegramChatID:   *telegramChatIDFlag,
	// }

}
