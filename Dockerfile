FROM golang:1.13-alpine as builder

WORKDIR /app

RUN apk update && apk add curl
RUN curl -L 'https://github.com/golang-migrate/migrate/releases/download/v4.5.0/migrate.linux-amd64.tar.gz' | tar xvz
RUN mv migrate.linux-amd64 /usr/bin/go-migrate && chmod +x /usr/bin/go-migrate

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./bin/bids

FROM alpine:latest

RUN apk update && apk add bash tzdata postgresql-client && rm -rf /var/cache/apk/*

COPY --from=builder /usr/bin/go-migrate /app/
COPY docker-entrypoint.sh .
RUN chmod a+x /docker-entrypoint.sh

COPY /migrations /app/migrations
COPY --from=builder /app/bin/bids /app/

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["./app/bids"]
