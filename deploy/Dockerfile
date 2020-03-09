FROM docker.io/library/golang:1.14-alpine as builder

MAINTAINER Jack Murdock <jack_murdock@comcast.com>

WORKDIR /go/src/github.com/xmidt-org/petasos

ARG VERSION=unknown
ARG GITCOMMIT=unknown
ARG BUILDTIME=unknown

ADD https://github.com/geofffranks/spruce/releases/download/v1.25.2/spruce-linux-amd64 /usr/local/bin/spruce
RUN chmod +x /usr/local/bin/spruce

RUN apk add --no-cache --no-progress \
    ca-certificates \
    make \
    git \
    openssh \
    gcc \
    libc-dev \
    upx

COPY . .
RUN make build

FROM alpine:latest

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/xmidt-org/petasos/petasos.yaml /petasos.yaml
COPY --from=builder /go/src/github.com/xmidt-org/petasos/petasos /petasos
COPY --from=builder /go/src/github.com/xmidt-org/petasos/deploy/Dockerfile /go/src/github.com/xmidt-org/petasos/NOTICE /go/src/github.com/xmidt-org/petasos/LICENSE /go/src/github.com/xmidt-org/petasos/CHANGELOG.md /
COPY --from=builder /go/src/github.com/xmidt-org/petasos/deploy/packaging/entrypoint.sh /entrypoint.sh
COPY --from=builder /go/src/github.com/xmidt-org/petasos/deploy/packaging/petasos_spruce.yaml /tmp/petasos_spruce.yaml
COPY --from=builder /usr/local/bin/spruce /spruce

RUN mkdir /etc/petasos/ && touch /etc/petasos/petasos.yaml && chmod 666 /etc/petasos/petasos.yaml

USER nobody

ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 6400
EXPOSE 6401
EXPOSE 6402
EXPOSE 6403

CMD ["/petasos"]