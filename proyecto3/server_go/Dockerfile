# =========================
# Dockerfile para server_go
# =========================

# 1️⃣ Imagen base de Go
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server_go main.go

# =========================
# Imagen final ligera
# =========================
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/server_go .

EXPOSE 50051
CMD ["./server_go"]
