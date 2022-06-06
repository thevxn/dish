# savla-dish (golang1.18)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/savla-dev/savla-dish)](https://pkg.go.dev/github.com/savla-dev/savla-dish)
[![Go Report Card](http://goreportcard.com/badge/github.com/savla-dev/savla-dish)](https://goreportcard.com/report/github.com/savla-dev/savla-dish)

  + __tiny__ monitoring one-shot service
  + __clusterable__ (peer-to-peer idea **TBD**), remote configuration of independent 'dish network' (`-source=$REMOTE_JSON_API_URL`)
  + __fast__ (quick load and exec time, 5s timeout per socket -- hardcoded), instant messaging connectors

```shell
$ go get github.com/savla-dev/savla-dish

$ savla-dish -h
Usage of go/bin/savla-dish:
  -source string
    	a string, path to/URL JSON socket list (default "demo_sockets.json")
  -telegram
    	a bool, Telegram provider usage toggle
  -telegramBotToken string
    	a string, Telegram bot private token
  -telegramChatID string
    	a string/signet int, Telegram chat/channel ID
  -verbose
    	a bool, console stdout logging toggle
```

## use-cases

the idea of a tiny one-shot service comes with the need of a quick monitoring service implementation over HTTP/S and generic TCP endpoints (more like 'sockets' = hosts and their ports)

### socket list

the list of sockets can be provided via a local JSON-formated file, or via remote REST/RESTful JSON-returning API (JSON structure has to be of the same structure anyway; see `demo_sockets.json`)

### alerting

as the alerting system (in case of socket test timeout threshold hit, or an unexpected HTTP response code) we provide a simple embedded `messenger` with Telegram IM implementation example (see `messenger/messenger.go`); since the Telegram bot token and the potential Telegram chat ID are considered as the **secrets**, we do recommend including these to the custom, local, binary executable instead of passing them into the CLI shell (security breach as secrets can then leak in process list)

## examples

```shell
# get the actual git version
go get github.com/savla-dev/savla-dish

# load sockets from demo_sockets.json file (by default) and use telegram provider for alerting (hardcoded token and chatID -- messenger/messenger.go)
savla-dish -source=demo_sockets.json -telegram

# use remote RESTful API service's socket list, use _explicit_ telegram bot and chat
savla-dish -source='https://api.example.com/dish/source' -telegram -telegramChatID=-123456789 -telegramBotToken='idk:00779988ddd'
```

### cronjob example

```shell
# non-root user!
crontab -e
```

```shell
# m h  dom mon dow   command
MAILTO=monitoring@example.com

TELEGRAM_TOKEN="000001:AFFDS45454d5ccfsadf34" 
TELEGRAM_CHATID="-12345678900"
DISH_EXECUTABLE_PATH=/home/dish/golang/bin/savla-dish
DISH_SOURCE=http://restapi.example.com/dish/sockets/${HOSTNAME}

*/1 * * * * ${DISH_EXECUTABLE_PATH} -source=${DISH_SOURCE} -telegram -telegramBotToken=${TELEGRAM_TOKEN} -telegramChatID=${TELEGRAM_CHATID}
```

Please note, that `savla-dish` executable returns "dish run: all tests ok" and exit code `0`, as soon as the execution ends (and no problems are present to report).
