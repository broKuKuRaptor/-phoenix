# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-alpine AS builder

ENV GOTOOLCHAIN=auto

RUN apk add --no-cache ca-certificates git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG CMD_PATH=./cmd/accounts
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/phoenix ${CMD_PATH}

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /out/phoenix /app/phoenix

ENTRYPOINT ["/app/phoenix"]
