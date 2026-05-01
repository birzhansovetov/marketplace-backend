# ── Stage 1: build ──────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Зависимости отдельным слоем — кэшируется если go.mod не менялся
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o marketplace-api ./cmd/...

# ── Stage 2: run ─────────────────────────────────────────────
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/marketplace-api .

EXPOSE 8080

CMD ["./marketplace-api"]