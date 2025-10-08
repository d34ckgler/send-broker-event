# Build stage
FROM golang:1.24.1-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum
COPY go.mod go.sum* ./

# Descargar dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o send-broker-event .

# Runtime stage
FROM alpine:latest

# Instalar ca-certificates para conexiones HTTPS
RUN apk --no-cache add ca-certificates

# Crear directorio de trabajo
WORKDIR /app

# Copiar el binario compilado desde el builder
COPY --from=builder /app/send-broker-event /usr/local/bin/send-broker-event

# Copiar el archivo .env (opcional, mejor usar variables de entorno)
COPY .env .env

# Dar permisos de ejecución
RUN chmod +x /usr/local/bin/send-broker-event

# El contenedor quedará en ejecución esperando comandos
# Usamos tail -f /dev/null para mantener el contenedor vivo sin consumir recursos
CMD ["tail", "-f", "/dev/null"]
