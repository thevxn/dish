#
# savla-dish / Dockerfile
#

# https://hub.docker.com/_/golang

ARG ALPINE_VERSION
ARG GOLANG_VERSION
FROM golang:${GOLANG_VERSION}-alpine as dish-build

LABEL maintainer="krusty@savla.dev, tack@savla.dev, krixlion@savla.dev"

ARG APP_NAME

WORKDIR /go/src/${APP_NAME}
COPY . /go/src/${APP_NAME}
COPY .docker/resolv.conf /etc/resolv.conf

# build binary
RUN go install cmd/main.go 

FROM alpine:${ALPINE_VERSION} as dish-release

WORKDIR /usr/local/bin
COPY --from=dish-build /go/bin/main .

ENTRYPOINT [ "./main" ]
