## Эндпоинты


- **GET /api/v1/songs -** Получение данных библиотеки с фильтрацией по всем полям и пагинацией
- **GET /api/v1/song/:id/lyrics -** Получение текста песни с пагинацией по куплетам
- **DELETE /api/v1/song/:id -** Удаление песни
- **PUT /api/v1/song -** Изменение данных песни
- **POST /api/v1/song -** Добавление новой песни в формате JSON

## Используемые технологии

- Go
- PostgreSQL
- Docker
- Gin
- Uber-fx
- Swagger


## Настройка окружения

1. Скопируйте файл `.env.example` в `.env` и заполните необходимые переменные окружения:
   ```bash
   cp .env.example .env
   ```

2. Убедитесь, что у вас установлен Docker и Docker Compose.

## Запуск приложения из Docker

1. Соберите и запустите контейнеры с помощью Docker Compose:
   ```bash
   docker-compose up --build
   ```

Приложение будет доступно по адресу `http://localhost:{APP_PORT}/swagger/index.html#/`

## Запуск приложения локально

1. Установка зависимостей
   ```bash
   go mod download
   ```

2. Запустите мигратор:
   ```bash
   go run cmd/migrator/main.go -up
   ```

3. Запустите приложение:
   ```bash
   go run cmd/app/main.go
   ```
Приложение будет доступно по адресу `http://localhost:{APP_PORT}/swagger/index.html#/`

## Логирование

Логи приложения можно найти в папке `logs`. Они помогут вам отслеживать работу сервера и выявлять возможные ошибки.