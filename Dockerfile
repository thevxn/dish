#
# savla-dish / Dockerfile
#

# https://hub.docker.com/_/golang

FROM golang:1.17-alpine

MAINTAINER krusty@savla.dev
MAINTAINER tack@savla.dev

ARG APP_NAME
ARG DOCKER_DEV_PORT

ENV APP_NAME ${APP_NAME}
ENV DOCKER_DEV_PORT ${DOCKER_DEV_PORT}

WORKDIR /go/src/${APP_NAME}
COPY . /go/src/${APP_NAME}
COPY .docker/resolv.conf /etc/resolv.conf

RUN go mod init
RUN go get -d -v ./...
RUN go install 

EXPOSE ${DOCKER_DEV_PORT}
CMD ${APP_NAME}

