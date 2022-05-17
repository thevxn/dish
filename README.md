# savla-dish (golang1.17)

tiny monitoring one-shot service, clusterable (peer-to-peer idea TBD) and fast (quick load and exec time, 5s timeout per socket)

## use-cases

the idea of a tiny one-shot service comes with the need of a quick monitoring service implementation over HTTP/S and generic TCP endpoints (more like 'sockets' = hosts and their ports)

the list of sockets can be provided via a local JSON-formated file, or via remote REST/RESTful JSON-returning API (JSON structure has to be of the same structure anyway; see `demo_sockets.json`)

as the alerting system (in case of socket test timeout threshold hit, or an unexpected HTTP response code) we provide a simple embedded `messenger` with Telegram IM implementation example (see `messenger/messenger.go`); since the Telegram bot token and the potential Telegram chat ID are considered as the **secrets**, we do recommend including these to the custom, local, binary executable instead of passing them into the CLI shell (security breach as secrets can then leak in process list)

```
# get the actual git version
go get github.com/savla-dev/savla-dish

# load sockets from demo_sockets.json file (by default) and use telegram provider for alerting
savla-dish -source=demo_sockets.json -telegram
savla-dish -source=demo_sockets.json -telegram -telegram_chat_id=-123456789

# use remote RESTful API service's socket list, use _explicit_ telegram bot and chat
savla-dish -source='https://api.example.com/dish/source' -telegram -telegram_chat_id=-123456789 -telegram_bot_token='idk:00779988ddd'
```


## dev environment

~~package `messenger` contains sensitive data (token, chat_id etc), therefore it is secured/exclusively excluded in build time (may raise an error)~~

```
# build binary/module with all packages
go build -tags dev savla-dish
```

### environment variables

```
export $(cat .env | sed -e 'd/^#.*/' | xargs) 

./savla-dish
```

### cronjob example

```
# non-root user
cronjob -u
```

```
DISH_ENVIRONMENT=prod
DISH_SOCKET_SOURCE="http://restapi.endpoint.com/dish"
DISH_TELEGRAM_BOT_TOKEN="token:whatever"
DISH_TELEGRAM_CHAT_ID="-123456789"

*/1 * * * * /path/to/savla-dish -flags
```
