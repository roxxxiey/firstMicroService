# Stage 1: Build stage
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код приложения
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Run stage
FROM alpine:latest

# Добавляем CA сертификаты для HTTPS
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированный бинарник из build-стадии
COPY --from=builder /app/main .

# Запускаем приложение
CMD ["./main"]
