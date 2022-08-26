#
# savla-dish / Dockerfile
#

# https://hub.docker.com/_/golang

ARG GOLANG_VERSION
FROM golang

ARG APP_NAME

ENV APP_NAME ${APP_NAME}

WORKDIR /go/src/${APP_NAME}
COPY . /go/src/${APP_NAME}
COPY .docker/resolv.conf /etc/resolv.conf

RUN go version; go env
RUN go build main.go

CMD [ ./main ]

