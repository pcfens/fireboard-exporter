FROM golang:1.17-alpine AS build

WORKDIR /go/src/github.com/pcfens/firebase-exporter
COPY . .

RUN apk add --no-cache ca-certificates \
    && CGO_ENABLED=0 GOOS=linux go build -a -o firebase-exporter

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build --chmod=755 /go/src/github.com/pcfens/firebase-exporter/firebase-exporter firebase-exporter

ENTRYPOINT ["/firebase-exporter"]
