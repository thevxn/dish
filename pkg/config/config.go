package config

import (
	"flag"
)

var (
	InstanceName     string
	HeaderName       string
	HeaderValue      string
	Source           string
	Verbose          bool
	TargetURL        string
	UsePushgateway   bool
	UseTelegram      bool
	TelegramBotToken string
	TelegramChatID   string
	Timeout          int // In seconds
	UpdateStates	 bool
	UpdateURL	 string
)

func init() {
	instanceName := flag.String("name", "generic-dish", "a string, dish instance name")
	timeoutFlag := flag.Int("timeout", 10, "an int, timeout in seconds for http and tcp calls")
	verboseFlag := flag.Bool("verbose", true, "a bool, console stdout logging toggle")

	sourceFlag := flag.String("source", "demo_sockets.json", "a string, path to/URL JSON socket list")
	sourceHeaderName := flag.String("hname", "", "a string, custom additional header name")
	sourceHeaderValue := flag.String("hvalue", "", "a string, custom additional header value")

	usePushgatewayFlag := flag.Bool("pushgw", false, "a bool, enable reporter module to post dish results to pushgateway")
	targetURLFlag := flag.String("target", "", "a string, result update path/URL, plaintext/byte output")

	// telegram provider flags
	useTelegramFlag := flag.Bool("telegram", false, "a bool, Telegram provider usage toggle")
	telegramBotTokenFlag := flag.String("telegramBotToken", "", "a string, Telegram bot private token")
	telegramChatIDFlag := flag.String("telegramChatID", "", "a string/signet int, Telegram chat/channel ID")

	updateStateFlag := flag.Bool("update", false, "a bool, switch for socket's last state batch upload to the source swis-api instance")
	updateURLFlag := flag.String("updateURL", "", "a string, URL of the source swis-api instance")
	
	flag.Parse()

	// system vars
	InstanceName = *instanceName
	Timeout = *timeoutFlag
	Verbose = *verboseFlag

	// source vars
	Source = *sourceFlag
	HeaderName = *sourceHeaderName
	HeaderValue = *sourceHeaderValue

	// target vars
	UsePushgateway = *usePushgatewayFlag
	TargetURL = *targetURLFlag

	// telegram vars
	UseTelegram = *useTelegramFlag
	TelegramBotToken = *telegramBotTokenFlag
	TelegramChatID = *telegramChatIDFlag

	UpdateStates = *updateStateFlag
	UpdateURL = *updateURLFlag
}
