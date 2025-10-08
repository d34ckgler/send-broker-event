# Docker Setup para send-broker-event

Este proyecto está configurado para ejecutarse en Docker sin iniciar automáticamente la aplicación. El contenedor permanece en ejecución y puedes ejecutar comandos cuando lo necesites.

## Construcción de la imagen

### Opción 1: Usando Docker directamente
```bash
docker build -t send-broker-event .
```

### Opción 2: Usando Docker Compose
```bash
docker-compose build
```

## Iniciar el contenedor

### Opción 1: Usando Docker directamente
```bash
docker run -d --name send-broker-event send-broker-event
```

### Opción 2: Usando Docker Compose
```bash
docker-compose up -d
```

## Ejecutar el programa

Una vez que el contenedor está en ejecución, puedes ejecutar el comando con:

```bash
# Sintaxis general
docker exec send-broker-event send-broker-event <evento> [argumentos...]

# Ejemplo
docker exec send-broker-event send-broker-event product.updated
```

### Con Docker Compose:
```bash
docker-compose exec send-broker-event send-broker-event product.updated
```

## Verificar logs

```bash
# Ver logs del contenedor
docker logs send-broker-event

# Seguir logs en tiempo real
docker logs -f send-broker-event
```

## Detener y eliminar el contenedor

### Con Docker:
```bash
docker stop send-broker-event
docker rm send-broker-event
```

### Con Docker Compose:
```bash
docker-compose down
```

## Variables de entorno

Puedes configurar las variables de entorno de tres formas:

1. **Archivo .env**: Se copia al contenedor durante la construcción
2. **docker-compose.yml**: Define las variables en la sección `environment`
3. **Línea de comandos**: Usa `-e` al ejecutar `docker run`

```bash
docker run -d --name send-broker-event \
  -e CENTRAL_RABBITHOST=192.168.30.55 \
  -e CENTRAL_RABBITPORT=5673 \
  send-broker-event
```

## Notas importantes

- El contenedor usa `tail -f /dev/null` para mantenerse en ejecución sin consumir recursos
- El binario se instala en `/usr/local/bin/send-broker-event` y está disponible globalmente
- La imagen usa Alpine Linux para mantener el tamaño pequeño
- Se utiliza compilación multi-stage para optimizar el tamaño de la imagen final
