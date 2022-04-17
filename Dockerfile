# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# App code
COPY . ./
RUN go build -o ./arisa3

ENTRYPOINT ["./arisa3"]
