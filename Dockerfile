FROM golang

ARG APP_NAME

ENV APP_NAME ${APP_NAME}

WORKDIR /go/src/${APP_NAME}
ADD . /go/src/${APP_NAME}
ADD .docker/resolv.conf /etc/resolv.conf

RUN go build cmd/main.go

CMD [ ./main ]