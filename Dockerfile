#
# savla-dish / Dockerfile
#

# https://hub.docker.com/_/golang

ARG GOLANG_VERSION=1.18
FROM golang:${GOLANG_VERSION}-alpine

MAINTAINER krusty@savla.dev
MAINTAINER tack@savla.dev

ARG APP_NAME

ENV GOLANG_VERSION ${GOLANG_VERSION}
ENV APP_NAME ${APP_NAME}
ENV APP_VERSION ${APP_NAME}_${GOLANG_VERSION}

WORKDIR /go/src/${APP_NAME}
COPY . /go/src/${APP_NAME}
COPY .docker/resolv.conf /etc/resolv.conf

# run "build job"
RUN go mod init 
RUN go get -d -v ./...
RUN go version; go env
RUN go install -tags dev ${APP_NAME} && \
	ln -s ${GOPATH}/bin/${APP_NAME} ${GOPATH}/bin/${APP_VERSION}

CMD ${APP_VERSION}

