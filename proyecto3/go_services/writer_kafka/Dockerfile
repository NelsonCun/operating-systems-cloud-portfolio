# Etapa de compilación
FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar binario
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kafka_writer .

# Imagen final
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/kafka_writer .

EXPOSE 6001

CMD ["./kafka_writer"]
