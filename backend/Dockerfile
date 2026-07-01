# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache ca-certificates git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
COPY db/migrations /db/migrations

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget

WORKDIR /app

COPY --from=builder /api /app/api
COPY --from=builder /db/migrations /app/db/migrations

ENV MIGRATIONS_PATH=/app/db/migrations
ENV PORT=8080

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["/app/api"]
