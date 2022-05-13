# savla-dish

tiny monitoring service, clusterable (peer-to-peer) and fast

## dev environment

package `messenger` contains sensitive data (token, chat_id etc), therefore it is secured/exclusively excluded in build time (may raise an error)

```
# build binary/module with all packages
go build -tags dev savla-dev
```
