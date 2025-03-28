<h1 align="left">
<img alt="dish_logo" src="https://vxn.dev/logos/dish.svg" width="90" height="90">
dish
</h1>

[![PkgGoDev](https://pkg.go.dev/badge/go.vxn.dev/dish)](https://pkg.go.dev/go.vxn.dev/dish)
[![Go Report Card](https://goreportcard.com/badge/go.vxn.dev/dish)](https://goreportcard.com/report/go.vxn.dev/dish)
[![libs.tech recommends](https://libs.tech/project/468033120/badge.svg)](https://libs.tech/project/468033120/dish)

+ __tiny__ one-shot monitoring service
+ __remote__ configuration of independent 'dish network' (via `-source ${REMOTE_JSON_API_URL}` flag)
+ __fast__ concurrent testing, low overall execution time, 10-sec timeout per socket by default
+ __0__ dependencies

## Install

```shell
go install go.vxn.dev/dish/cmd/dish@latest
```

## Usage

```
dish -h
Usage of dish:
  -failedOnly
        a bool, specifies whether only failed checks should be reported (default true)
  -hname string
        a string, custom additional header name
  -hvalue string
        a string, custom additional header value
  -name string
        a string, dish instance name (default "generic-dish")
  -source string
        a string, path to/URL JSON socket list (default "./configs/demo_sockets.json")
  -target string
        a string, result update path/URL to pushgateway, plaintext/byte output
  -telegramBotToken string
        a string, Telegram bot private token
  -telegramChatID string
        a string, Telegram chat/channel ID
  -timeout uint
        an int, timeout in seconds for http and tcp calls (default 10)
  -updateURL string
        a string, URL of the source api instance
  -verbose
        a bool, console stdout logging toggle
  -webhookURL string
        a string, URL of webhook endpoint
```

### Socket List

The list of sockets can be provided via a local JSON-formated file (e.g. `demo_sockets.json` file in the CWD), or via a remote REST/RESTful JSON API.

```bash
# local JSON file
dish -source /opt/dish/sockets.json

# remote JSON API source
dish -source http://restapi.example.com/dish/sockets/:instance
```

### Alerting

When a socket test fails, it's always good to be notified. For this purpose, dish provides 4 different ways of doing so (can be combined):

+ test results upload to a remote JSON API (via `-updateURL` flag)
+ failed sockets list as the Telegram message body (via Telegram-related flags, see the help output above)
+ failed count and last test timestamp update to Pushgateway for Prometheus (via the `-target` flag)
+ test results push to a webhook URL (via the `webhookURL` flag)

![telegram-alerting](/.github/dish-telegram.png)

(The screenshot above shows the Telegram alerting as of `v1.5.0`.)

### Examples

One way to run dish is to build and install a binary executable.

```shell
# Fetch and install the specific version
go install go.vxn.dev/dish/cmd/dish@latest

export PATH=$PATH:~/go/bin

# Load sockets from sockets.json file, and use Telegram 
# provider for alerting
dish -source sockets.json -telegram -telegramChatID "-123456789" \
 -telegramBotToken "123:AAAbcD_ef"

# Use remote JSON API service as socket source, and push
# the results to Pushgateway
dish -source https://api.example.com/dish/sockets -pushgw \
 -target https://pushgw.example.com/
```

#### Using Docker

```shell
# Copy, and/or edit dot-env file (optional)
cp .env.example .env
vi .env

# Build a Docker image
make build

# Run using docker compose stack
make run

# Run using native docker run
docker run --rm \
 dish:1.7.1-go1.23 \
 -verbose \
 -source https://api.example.com \
 -pushgw \
 -target https://pushgateway.example.com
```

#### Bash script and cronjob

Create a bash script to easily deploy dish and update its settings:

```shell
vi tiny-dish-run.sh
```

```shell
#!/bin/bash

TELEGRAM_TOKEN="123:AAAbcD_ef"
TELEGRAM_CHATID="-123456789"

SOURCE_URL=https://api.example.com/dish/sockets
UPDATE_URL=https://api.example.com/dish/sockets/results
TARGET_URL=https://pushgw.example.com

DISH_TAG=dish:1.6.0-go1.22
INSTANCE_NAME=tiny-dish

SWAPI_TOKEN=AbCd

docker run --rm \
        ${DISH_TAG} \
        -name ${INSTANCE_NAME} \
        -source ${SOURCE_URL} \
        -hvalue ${SWAPI_TOKEN} \
        -hname X-Auth-Token \
        -target ${TARGET_URL} \
        -updateURL ${UPDATE_URL} \
        -telegramBotToken ${TELEGRAM_TOKEN} \
        -telegramChatID ${TELEGRAM_CHATID} \
        -timeout 15 \
        -verbose
```

Make it an executable:

```shell
chmod +x tiny-dish-run.sh
```

##### Cronjob to run periodically

```shell
crontab -e
```

```shell
# m h  dom mon dow   command
MAILTO=monitoring@example.com

*/2 * * * * /home/user/tiny-dish-run.sh
```

## History

[dish history article](https://krusty.space/projects/dish/)

## Use Cases

The idea of a tiny one-shot service comes with the need for a quick monitoring service implementation to test HTTP/S and generic TCP endpoints (or just sockets in general = hosts and their ports).
