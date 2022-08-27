FROM golang

ARG APP_NAME

ENV APP_NAME ${APP_NAME}

WORKDIR /go/src/${APP_NAME}
COPY . /go/src/${APP_NAME}
COPY .docker/resolv.conf /etc/resolv.conf

RUN go build main.go

CMD [ ./main ]