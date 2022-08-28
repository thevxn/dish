#
# savla-dish / Dockerfile
#

# https://hub.docker.com/_/golang

ARG GOLANG_VERSION
FROM golang:${GOLANG_VERSION}-alpine as dish-build

LABEL maintainer="krusty@savla.dev, tack@savla.dev"

ARG APP_NAME
ARG APP_FLAGS

ENV APP_FLAGS ${APP_FLAGS}}
ENV GOLANG_VERSION ${GOLANG_VERSION}
ENV APP_NAME ${APP_NAME}
ENV APP_VERSION ${APP_NAME}_${GOLANG_VERSION}

WORKDIR /go/src/${APP_NAME}
COPY . /go/src/${APP_NAME}
COPY .docker/resolv.conf /etc/resolv.conf

# build binary
RUN go build cmd/main.go 

ENTRYPOINT [ "./main", "${APP_FLAGS}" ]