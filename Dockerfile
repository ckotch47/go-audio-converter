# =============================================
# STAGE 1: Сборка (Builder)
# =============================================
FROM golang:1.25-alpine AS builder

# Устанавливаем зависимости для сборки Go и FFmpeg
RUN apk add --no-cache \
    gcc \
    musl-dev \
    ffmpeg 


# Копируем модули и скачиваем зависимости
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник с CGO_ENABLED=1 (для FFmpeg)
RUN CGO_ENABLED=1  go build -a -installsuffix cgo -o server .

# =============================================
# STAGE 2: Финальный образ (Runtime)
# =============================================
FROM alpine:latest

# Устанавливаем FFmpeg и ca-certificates (для HTTPS)
RUN apk add --no-cache \
    ffmpeg \
    tzdata

# Создаём рабочую директорию и временные папки
RUN mkdir -p /app/uploads /tmp

# Копируем бинарник
COPY --from=builder /app/server /app/server
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/public /app/public

# Устанавливаем права
RUN chmod +x /app/server


# Работаем в /app
WORKDIR /app

# Экспорт порта
EXPOSE 3000

# Запускаем сервер
CMD ["./server"]