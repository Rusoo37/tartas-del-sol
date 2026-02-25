# ==========================================
# Etapa 1: Compilación (Builder)
# ==========================================
# Usamos la imagen oficial de Go en su última versión liviana
FROM golang:alpine AS builder

# Creamos y nos paramos en la carpeta /app adentro del contenedor
WORKDIR /app

# Copiamos el archivo go.mod (y todo el resto del código)
COPY go.mod ./
COPY . .

# Compilamos la aplicación para que sea un ejecutable independiente
RUN CGO_ENABLED=0 GOOS=linux go build -o servidor-tartas main.go

# ==========================================
# Etapa 2: Imagen Final de Producción
# ==========================================
# Usamos Alpine, una versión de Linux que pesa menos de 5MB
FROM alpine:latest

# Nos paramos en la carpeta principal
WORKDIR /root/

# Copiamos SOLO el binario compilado desde la etapa anterior
COPY --from=builder /app/servidor-tartas .

# Copiamos los archivos sueltos que el servidor necesita para andar
COPY --from=builder /app/tartas.json .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Le avisamos al contenedor que nuestra app usa el puerto 8080
EXPOSE 8080

# Comando final que arranca tu página web
CMD ["./servidor-tartas"]