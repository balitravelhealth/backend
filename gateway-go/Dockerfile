# ── Stage 1: build ────────────────────────────────────────────────────────
FROM golang:1.25.10-alpine3.23 AS builder

WORKDIR /app

# Layer cache: unduh dependency dulu, baru salin source
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o gateway ./cmd/server

# ── Stage 2: runtime (image minimal) ──────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/gateway .

EXPOSE 8080
CMD ["./gateway"]
