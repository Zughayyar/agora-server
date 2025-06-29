FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY ./internal ./internal
COPY ./cmd ./cmd
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN mkdir -p bin
RUN go build -o bin/server ./cmd/server
RUN go build -o bin/migration ./cmd/migration

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/bin/server /app/bin/server
COPY --from=builder /app/bin/migration /app/bin/migration

CMD ["./bin/server"]