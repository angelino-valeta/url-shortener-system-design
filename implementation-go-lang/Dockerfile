# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copia go mod e sum primeiro (melhor cache)
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN go build -o url-shortener cmd/main.go

# Imagem final leve
FROM alpine:latest

WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/url-shortener .

# Expõe a porta
EXPOSE 8080

# Comando para rodar
CMD ["./url-shortener"]