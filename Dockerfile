FROM golang:1.24-alpine AS build

RUN apk add --no-cache git ca-certificates tzdata postgresql-client

RUN adduser -D -s /bin/sh -u 1001 appuser

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY migrations/ ./migrations/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -o /bookstore-api ./cmd/main.go

RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.16.0/migrate.linux-amd64.tar.gz \
    && tar -xzf migrate.tar.gz migrate \
    && chmod +x migrate

FROM alpine

RUN apk add --no-cache ca-certificates tzdata postgresql-client

WORKDIR /app

COPY --from=build /bookstore-api /bookstore-api
COPY --from=build /app/migrations/ /migrations/
COPY --from=build /app/migrate /usr/local/bin/migrate
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

USER 1001

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/bookstore-api"]
