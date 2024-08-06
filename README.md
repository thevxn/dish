# dish

[![PkgGoDev](https://pkg.go.dev/badge/github.com/thevxn/dish)](https://pkg.go.dev/github.com/thevxn/dish)
[![Go Report Card](http://goreportcard.com/badge/github.com/thevxn/dish)](https://goreportcard.com/report/github.com/thevxn/dish)

+ __tiny__ one-shot monitoring service
+ __remote__ configuration of independent 'dish network' (via `-source ${REMOTE_JSON_API_URL}` flag)
+ __fast__ parallel testing, low overall execution time, 10-sec timeout per socket by default

```shell
$ go install go.vxn.dev/dish@1.6.0

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
  -update
    	a bool, switch for socket's last state batch upload to the source swis-api instance
  -updateURL string
    	a string, URL of the source swis-api instance
  -verbose
    	a bool, console stdout logging toggle (default true)
```

[dish history article](https://krusty.space/projects/savla-dish/)

## use-cases

The idea of a tiny one-shot service comes with the need for a quick monitoring service implementation to test HTTP/S and generic TCP endpoints (or just sockets in general = hosts and their ports).

### socket list

The list of sockets can be provided via a local JSON-formated file (e.g. `demo_sockets.json` file in the CWD), or via a remote REST/RESTful JSON API.

```bash
# local JSON file
dish -source=/opt/dish/sockets.json

# remote JSON API source
dish -source=http://restapi.example.com/dish/sockets/:instance
```

### alerting

When a socket test fails, it's always good to be notified. For this purpose, dish provides three different ways of doing so (can be combined):

+ test results upload to a remote JSON API (via `-updateURL` flag)
+ failed sockets list as the Telegram message body (via Telegram-related flags, see the help output above)
+ failed count and last test timestamp update to Pushgateway for Prometheus (via `-pushgw` and `-target` flags)

![telegram-alerting](/.github/savla-dish-telegram.png)

(The screenshot above shows the Telegram alerting as of `v1.5.0`.)

## examples

One way to run dish is to build and install a binary executable.

```shell
# Fetch and install the specific version
go install go.vxn.dev/dish@1.6.0

# Load sockets from sockets.json file, and use Telegram 
# provider for alerting
dish -source sockets.json -telegram -telegramChatID "-123456789" \
	-telegramBotToken "123:AAAbcD_ef"

# Use remote JSON API service as socket source, and push
# the results to Pushgateway
dish -source https://api.example.com/dish/sockets -pushgw \
	-target https://pushgw.example.com/
```

### using Docker

```shell
# Copy, and/or edit dot-env file (optional)
cp .env.example .env
vi .env

# Build a Docker image
make build

# Run
docker run --rm \
	dish:1.6.0-go1.22 \
	-verbose \
	-source https://api.example.com \
	-pushgw \
	-target https://pushgateway.example.com
```

### bash script and cronjob

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
        -pushgw \
        -target ${TARGET_URL} \
        -update \
        -updateURL ${UPDATE_URL} \
        -telegram \
        -telegramBotToken ${TELEGRAM_TOKEN} \
        -telegramChatID ${TELEGRAM_CHATID} \
        -timeout 15 \
        -verbose
```

Make it an executable:

```shell
chmod +x tiny-dish-run.sh
```

#### cronjob to run periodically

```shell
crontab -e
```

```shell
# m h  dom mon dow   command
MAILTO=monitoring@example.com

*/2 * * * * /home/user/tiny-dish-run.sh
```

