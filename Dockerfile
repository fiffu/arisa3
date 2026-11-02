# syntax=docker/dockerfile:1

FROM golang:1.24

WORKDIR /app

COPY . ./
RUN go build -o ./arisa3

# Secrets and DATABASE_URL come from env vars, which override the sample's placeholders.
ENTRYPOINT ["./arisa3", "-config-file", "config.sample.yml"]
