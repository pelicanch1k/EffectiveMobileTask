# Используем официальный образ Go
FROM golang:1.23.4 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

RUN go version
ENV GOPATH=/

# Копируем файлы проекта
COPY . .

# Собираем приложение
RUN go mod download
RUN go build -o main.go ./cmd/app/main.go

CMD ["./main.go"]