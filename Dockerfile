# Используем официальный образ Go
FROM golang:1.23.4 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Устанавливаем зависимости
RUN go mod download

# Собираем приложение и мигратор
RUN go build -o app ./cmd/app/main.go
RUN go build -o migrator ./cmd/migrator/main.go

# Команда по умолчанию для запуска приложения (будет переопределена в docker-compose.yml для мигратора)
CMD ["./app"]
