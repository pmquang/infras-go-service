FROM golang:1.13.7-alpine3.11 AS build-base

ENV GO111MODULE auto
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0

WORKDIR /root/project
COPY go.mod go.sum ./

RUN apk add --no-cache openssh-client git \
 && go mod download

FROM build-base AS builder
COPY . .
ARG VERSION
RUN go build -o /root/app --ldflags "-w -s -X main.version=${VERSION:-0.unknown}" .

FROM alpine:3.11
COPY --from=builder /root/app /root/app
ENTRYPOINT [ "/root/app" ]
