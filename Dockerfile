# syntax=docker/dockerfile:1

FROM golang:1.18

WORKDIR /app

COPY . ./
RUN go build -o ./arisa3

ENTRYPOINT ["./arisa3"]
