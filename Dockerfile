FROM golang:1.15.8-alpine AS golang

RUN apk add --no-cache \
  git

WORKDIR /go/src
COPY . .

ENV CGO_ENABLED=0
RUN go get -d -v ./... && \
  go install -v ./...

FROM litusproject/php_codesniffer:latest AS php_codesniffer

COPY --from=golang /go/bin/php-codesniffer-action /usr/bin/

WORKDIR /
COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
