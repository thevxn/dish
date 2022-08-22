#
# savla-dish / Dockerfile
#

# https://hub.docker.com/_/golang

ARG GOLANG_VERSION
FROM golang:${GOLANG_VERSION}-alpine as dish-build

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
RUN go install -tags dev ${APP_NAME}

#
# "flatten" the base image
#

FROM alpine:3.16 as dish-binary

COPY --from=dish-build /go/bin/${APP_NAME} /usr/local/bin/${APP_NAME}

ENTRYPOINT [ ${APP_NAME} ]

