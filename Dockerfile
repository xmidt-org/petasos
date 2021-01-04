FROM docker.io/library/golang:1.15-alpine as builder

MAINTAINER Jack Murdock <jack_murdock@comcast.com>

WORKDIR /src

ARG VERSION
ARG GITCOMMIT
ARG BUILDTIME


RUN apk add --no-cache --no-progress \
    ca-certificates \
    make \
    git \
    openssh \
    gcc \
    libc-dev \
    upx

RUN go get github.com/geofffranks/spruce/cmd/spruce && chmod +x /go/bin/spruce
COPY . .
RUN make test release

FROM alpine:3.12.1

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/petasos /src/petasos.yaml /src/deploy/packaging/entrypoint.sh /go/bin/spruce /src/Dockerfile /src/NOTICE /src/LICENSE /src/CHANGELOG.md /
COPY --from=builder /src/deploy/packaging/petasos_spruce.yaml /tmp/petasos_spruce.yaml

RUN mkdir /etc/petasos/ && touch /etc/petasos/petasos.yaml && chmod 666 /etc/petasos/petasos.yaml

USER nobody

ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 6400
EXPOSE 6401
EXPOSE 6402
EXPOSE 6403

CMD ["/petasos"]
