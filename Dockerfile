# Build stage
FROM golang:1.19.5-alpine3.17 AS builder
WORKDIR /app
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
COPY . .
RUN go build -o main main.go


# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY *.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration


EXPOSE 8080
EXPOSE 8090

CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
