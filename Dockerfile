# syntax=docker/dockerfile:1

FROM golang:1.26

WORKDIR /app

COPY . ./
RUN go build -o ./arisa3

# Send OpenTelemetry metrics (Go runtime, host, and app metrics) to Honeycomb.
# HONEYCOMB_API_KEY is a secret, so it comes from the vault instead.
ENV OTEL_METRICS_ENABLED=true
ENV HONEYCOMB_METRICS_DATASET=arisa3

# Secrets and DATABASE_URL come from env vars, which override the sample's placeholders.
ENTRYPOINT ["./arisa3", "-config-file", "config.sample.yml"]
