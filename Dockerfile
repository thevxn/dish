#
# savla-dish / Dockerfile
#

# https://hub.docker.com/_/golang

ARG GOLANG_VERSION_MINOR=1.17
FROM golang:${GOLANG_VERSION_MINOR}-alpine

MAINTAINER krusty@savla.dev
MAINTAINER tack@savla.dev

ARG APP_NAME

ENV GOLANG_VERSION_MINOR ${GOLANG_VERSION_MINOR}
ENV APP_NAME ${APP_NAME}
ENV APP_VERSION ${APP_NAME}_${GOLANG_VERSION_MINOR}

WORKDIR /go/src/${APP_NAME}
COPY . /go/src/${APP_NAME}
COPY .docker/resolv.conf /etc/resolv.conf

# run "build job"
RUN go mod init 
RUN go get -d -v ./...
RUN go version; go env
RUN go install ${APP_NAME} && \
	ln -s ${GOPATH}/bin/${APP_NAME} ${GOPATH}/bin/${APP_VERSION}

CMD ${APP_VERSION}

