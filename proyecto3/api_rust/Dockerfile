# Etapa de compilación
FROM rust:alpine AS builder

# Crear directorio de trabajo
WORKDIR /app

# Instalar dependencias necesarias para compilar
RUN apk add --no-cache protobuf gcc musl-dev

# Copiar todo el código fuente
COPY . .

# Compilar en modo release
RUN cargo build --release

# Etapa final (imagen ligera)
FROM alpine

# Crear usuario sin privilegios
RUN adduser -D app_user

# Crear directorio de trabajo
WORKDIR /app

# Copiar el binario compilado desde la etapa anterior
COPY --from=builder /app/target/release/api_rust ./api_rust

# Otorgar permisos correctos
RUN chown app_user:app_user ./api_rust

# Cambiar a usuario sin privilegios
USER app_user

# Exponer el puerto HTTP de la API
EXPOSE 8080

# Variable de entorno para el servidor gRPC
ENV GRPC_SERVER_URL=http://server-go-service:50051

# Comando por defecto
CMD ["./api_rust"]
