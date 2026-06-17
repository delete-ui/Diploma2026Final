FROM golang:1.25.4-alpine AS builder

RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-w -s" \
    -o /app/bin/api \
    ./cmd/api

FROM alpine:3.19

LABEL version="1.0.0"
LABEL description="Battery Shop API"

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    docker-cli \
    sudo \
    && update-ca-certificates

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

# Даём appuser права на выполнение docker без пароля
RUN echo "appuser ALL=(root) NOPASSWD: /usr/bin/docker" >> /etc/sudoers

WORKDIR /app

COPY --from=builder /app/bin/api .

COPY --from=builder /app/internal/migrate/migrations ./internal/migrate/migrations

RUN echo "ENV=production" > .env && \
    echo "DB_HOST=postgres" >> .env

RUN chown -R appuser:appgroup /app

USER appuser

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/batteries?limit=1 || exit 1

EXPOSE 8080

ENV GO_ENV=production \
    LOG_FORMAT=json \
    HTTP_ADDRESS=:8080

ENTRYPOINT ["./api"]