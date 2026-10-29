FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY . .
RUN go mod download
RUN go build -o things-mcp .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

# Hugging Face Spaces runs the container as uid 1000, so /data must be
# writable by that user. Root-run hosts (Render, Fly.io) are unaffected.
RUN mkdir -p /data && chown 1000:1000 /data

COPY --from=builder /app/things-mcp /usr/local/bin/things-mcp

EXPOSE 8080

ENV PORT=8080
ENV DATA_DIR=/data

VOLUME ["/data"]

CMD ["things-mcp"]
