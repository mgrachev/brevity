# Image for building an application with dependencies
FROM golang:1.9-alpine3.7 as builder

RUN apk upgrade --update \
 && apk --no-cache add postgresql-client openssl ca-certificates alpine-sdk

# Main application
RUN go get github.com/kardianos/govendor
COPY . /go/src/github.com/mgrachev/brevity
WORKDIR /go/src/github.com/mgrachev/brevity
RUN govendor sync

RUN GOOS=linux go build -a -installsuffix cgo cmd/brevity-http-server/brevity-http-server.go

# Postgres database manager
RUN go get github.com/lib/pq && go get github.com/urfave/cli
RUN mkdir -p /go/src/github.com/rnubel/pgmgr
RUN git clone https://github.com/rnubel/pgmgr /go/src/github.com/rnubel/pgmgr
RUN GOOS=linux go build -a -installsuffix cgo /go/src/github.com/rnubel/pgmgr/main.go

# Image for running the application
FROM alpine:3.7

RUN apk upgrade --update \
 && apk --no-cache add postgresql-client openssl ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /go/src/github.com/mgrachev/brevity/brevity-http-server app
COPY --from=builder /go/src/github.com/mgrachev/brevity/main pgmgr

COPY db /app/db/
ENV APP_ENV=production

CMD ["/app/app"]
