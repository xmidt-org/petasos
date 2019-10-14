FROM golang:alpine as builder
MAINTAINER Jack Murdock <jack_murdock@comcast.com>

WORKDIR /go/src/github.com/xmidt-org/petasos
ARG VERSION=undefined
ARG GITCOMMIT=undefined
ARG BUILDTIME=undefined

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

COPY . .
RUN GO111MODULE=on go build -ldflags "-X 'main.BuildTime=${BUILDTIME}' -X main.GitCommit=${GITCOMMIT} -X main.Version=${VERSION}" -o petasos_linux_amd64

FROM alpine

RUN apk --no-cache add ca-certificates
RUN mkdir -p /etc/petasos
VOLUME /etc/petasos

EXPOSE 8080

COPY --from=builder /go/src/github.com/xmidt-org/petasos/petasos_linux_amd64 /
COPY petasos.yaml /
COPY deploy/Dockerfile NOTICE LICENSE CHANGELOG.md /
ENTRYPOINT ["/petasos_linux_amd64"]
