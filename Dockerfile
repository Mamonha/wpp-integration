# BUILD STAGE
FROM golang:1.23-alpine AS builder

ARG BUILD_REF
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY .env ./
RUN go mod download

COPY . ./

# Build do binário
RUN env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.build=${BUILD_REF}" -o main ./main.go

# FINAL STAGE
FROM alpine:latest

WORKDIR /root/

# Copia o binário e o .env
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expondo a porta 8000
EXPOSE 8000

# Define o comando para iniciar o contêiner
CMD ["./main"]
