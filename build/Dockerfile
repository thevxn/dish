#
# dish / Dockerfile
#

#
# stage 0 --- build
#

# https://hub.docker.com/_/golang

ARG ALPINE_VERSION 3.20
ARG GOLANG_VERSION 1.23
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS dish-build

LABEL org.opencontainers.image.authors="krusty@vxn.dev, tack@vxn.dev, krixlion@vxn.dev"

ARG APP_NAME dish

WORKDIR /go/src/${APP_NAME}
COPY . /go/src/${APP_NAME}

# build and install the binary
RUN go install ./cmd/dish/

#
# stage 1 --- release
#

FROM alpine:${ALPINE_VERSION} AS dish-release

WORKDIR /opt

COPY configs/demo_sockets.json /opt/
COPY --from=dish-build /go/bin/dish /opt/dish
RUN ln -s /opt/dish /usr/local/bin

ENTRYPOINT [ "dish" ]
