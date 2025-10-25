# Etapa de construcción
FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copiar archivos de módulo Go
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consumer_rabbit .

# Imagen final mínima
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/consumer_rabbit .

CMD ["./consumer_rabbit"]
