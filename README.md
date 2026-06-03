# myproject

Учебный REST API на Go с PostgreSQL.

## Стек

- Go + net/http
- PostgreSQL + sqlx
- Layered architecture: handler -> service -> repository
- Docker + docker-compose

## Структура

```
cmd/myapp/        — точка входа
internal/
  config/         — конфиг из env
  handler/        — HTTP-хендлеры + middleware
  service/        — бизнес-логика
  repository/     — работа с БД
  models/         — структуры данных
migrations/       — SQL-миграции
```

## Запуск через Docker

```bash
docker compose up --build
```

Миграции применяются автоматически при первом запуске.

## Запуск без Docker

Скопируй `.env.example` в `.env` и заполни своими данными:

```bash
cp .env.example .env
go run ./cmd/myapp
```

## Эндпоинты

| Метод | URL         | Описание             |
|-------|-------------|----------------------|
| POST  | /users      | Создать пользователя |
| GET   | /users      | Список пользователей |
| GET   | /users/{id} | Получить по ID       |

### POST /users

```json
// запрос
{ "name": "Mark", "email": "mark@example.com" }

// ответ 201
{ "id": 1, "name": "Mark", "email": "mark@example.com" }
```
