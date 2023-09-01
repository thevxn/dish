# dish (golang1.21)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/savla-dev/savla-dish)](https://pkg.go.dev/github.com/savla-dev/savla-dish)
[![Go Report Card](http://goreportcard.com/badge/github.com/savla-dev/savla-dish)](https://goreportcard.com/report/github.com/savla-dev/savla-dish)

+ __tiny__ monitoring one-shot service
+ __remote__ configuration of independent 'dish network' (`-source=$REMOTE_JSON_API_URL`)
+ __fast__ (quick load and exec time, 10 sec timeout per socket by default), instant messenger connectors

```shell
$ go install go.savla.dev/dish@1.3.0

$ dish -h
Usage of ./dish:
  -hname string
     a string, custom additional header name
  -hvalue string
     a string, custom additional header value
  -name string
     a string, dish instance name (default "generic-dish")
  -pushgw
     a bool, enable reporter module to post dish results to pushgateway
  -source string
     a string, path to/URL JSON socket list (default "demo_sockets.json")
  -target string
     a string, result update path/URL, plaintext/byte output
  -telegram
     a bool, Telegram provider usage toggle
  -telegramBotToken string
     a string, Telegram bot private token
  -telegramChatID string
     a string/signet int, Telegram chat/channel ID
  -timeout int
     an int, timeout in seconds for http and tcp calls (default 10)
  -verbose
     a bool, console stdout logging toggle (default true)
```

## use-cases

[savla-dish history article](https://krusty.savla.dev/projects/savla-dish/)

the idea of a tiny one-shot service comes with the need of a quick monitoring service implementation over HTTP/S and generic TCP endpoints (more like 'sockets' = hosts and their ports)

it is not meant to be a competition with blackbox exporter, this is just another implementation approach

### socket list

the list of sockets can be provided via a local JSON-formated file, or via remote REST/RESTful JSON-returning API (JSON structure has to be of the same structure anyway; see `demo_sockets.json`)

```bash
./dish -source=http://restapi.example.com/dish/sockets/:instance
```

### alerting

as the alerting system (in case of socket test timeout threshold hit, or an unexpected HTTP response code) we provide a simple embedded `messenger` with Telegram IM implementation example (see `messenger/messenger.go`); since the Telegram bot token and the potential Telegram chat ID are considered as the __secrets__, we do recommend including these to the custom, local, binary executable instead of passing them into the CLI shell (security breach as secrets can then leak in process list --- to be reviewed)

![telegram-alerting](/.github/savla-dish-telegram.png)

### pushgateway

to keep dish simple and light, we decided not to import http server (even though net/http package is used) and use just its Client interface to push/post results to Pushgateway by Prometheus (TODO: insert pushgateway into docker-compose.yml config)

job name and instance name are hardcoded constants in the [reporter](/reporter/reporter.go) module source

[short article on motivation and history behind dish](https://krusty.savla.dev/projects/savla-dish/)

## examples

```shell
# get the actual git version
go install go.savla.dev/dish@latest

# load sockets from demo_sockets.json file (by default) and use telegram provider for alerting (hardcoded token and chatID -- messenger/messenger.go)
dish -source=demo_sockets.json -telegram

# use remote RESTful API service's socket list, use _explicit_ telegram bot and chat
dish -source='https://api.example.com/dish/source' -telegram -telegramChatID=-123456789 -telegramBotToken='idk:00779988ddd'
```

### docker it

we use `.env` and `gnumake` (Makefile) to simplify/semiautomate our development procedures, feel free to give it a try

```bash
# copy, and/or edit dot-env file
cp .env.example .env
vim .env

# build an image
make build

# run! (not the same as `make run`, but it should've been so)
docker run --rm -i savla-dish:golang-1.19 -verbose -pushgw -source=http://[...] -target=http://pushgateway.example.com
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
DISH_EXECUTABLE_PATH=/home/dish/golang/bin/dish
DISH_SOURCE=http://restapi.example.com/dish/sockets/${HOSTNAME}

*/1 * * * * ${DISH_EXECUTABLE_PATH} -source=${DISH_SOURCE} -telegram -telegramBotToken=${TELEGRAM_TOKEN} -telegramChatID=${TELEGRAM_CHATID}
```

Please note, that `dish` executable returns "dish run: all tests ok" and exit code `0`, as soon as the execution ends (and no problems are present to report).
