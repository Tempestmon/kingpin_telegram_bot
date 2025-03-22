# Этап 1: Сборка приложения
FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /bot .

# Этап 2: Создание Distroless-образа
FROM gcr.io/distroless/static-debian12

# Копируем бинарник из этапа сборки
COPY --from=builder /bot /bot

# Указываем точку входа
ENTRYPOINT ["/bot"]
