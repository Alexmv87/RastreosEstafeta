
FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

# Se compila la aplicación Go en un binario estático
RUN CGO_ENABLED=0 go build -o ./main main.go

FROM alpine

WORKDIR /app

# Se instalan certificados CA para conexiones HTTPS
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

# Se configura DNS para usar servidores de Google y Cloudflare   
RUN echo 'nameserver 8.8.8.8' > /etc/resolv.conf

RUN echo 'nameserver 1.1.1.1' >> /etc/resolv.conf

EXPOSE 9000

CMD ["./main"]