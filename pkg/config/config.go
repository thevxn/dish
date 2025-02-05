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
	UpdateStates     bool
	UpdateURL        string
	UseWebhooks      bool
	WebhookURL       string
)

func init() {
	// system vars
	flag.StringVar(&InstanceName, "name", "generic-dish", "a string, dish instance name")
	flag.IntVar(&Timeout, "timeout", 10, "an int, timeout in seconds for http and tcp calls")
	flag.BoolVar(&Verbose, "verbose", false, "a bool, console stdout logging toggle")

	// source vars
	flag.StringVar(&Source, "source", "./configs/demo_sockets.json", "a string, path to/URL JSON socket list")
	flag.StringVar(&HeaderName, "hname", "", "a string, custom additional header name")
	flag.StringVar(&HeaderValue, "hvalue", "", "a string, custom additional header value")

	// target vars
	flag.BoolVar(&UsePushgateway, "pushgw", false, "a bool, enable reporter module to post dish results to pushgateway")
	flag.StringVar(&TargetURL, "target", "", "a string, result update path/URL, plaintext/byte output")

	// telegram vars
	flag.BoolVar(&UseTelegram, "telegram", false, "a bool, Telegram provider usage toggle")
	flag.StringVar(&TelegramBotToken, "telegramBotToken", "", "a string, Telegram bot private token")
	flag.StringVar(&TelegramChatID, "telegramChatID", "", "a string/signet int, Telegram chat/channel ID")

	// remote source vars
	flag.BoolVar(&UpdateStates, "update", false, "a bool, switch for socket's last state batch upload to the source swis-api instance")
	flag.StringVar(&UpdateURL, "updateURL", "", "a string, URL of the source swis-api instance")

	// webhook vars
	flag.BoolVar(&UseWebhooks, "webhooks", false, "a bool, Webhook usage toggle")
	flag.StringVar(&WebhookURL, "webhookURL", "", "a string, URL of webhook endpoint")

	flag.Parse()

}
