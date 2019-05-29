# Building stage
FROM golang:1.11-alpine3.7

WORKDIR /banner-manager

RUN apk add --no-cache protobuf ca-certificates make git

# Source code, building tools and dependences
COPY src /banner-manager/src
COPY Makefile /banner-manager
# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod /banner-manager


ENV CGO_ENABLED 0
ENV GOOS linux
ENV TIMEZONE "Asia/Tokyo"
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime && \
    echo $TIMEZONE > /etc/timezone && \
    apk del tzdata
# Force the go compiler to use modules
ENV GO111MODULE=on

RUN go mod download
RUN make build
RUN mv banner-manager-client /go/bin
RUN mv banner-manager-server /go/bin

# Production stage
FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /banner-manager

# copy the go binaries from the building stage
COPY --from=0 /go/bin /banner-manager

COPY config /banner-manager/config
COPY data /banner-manager/data
