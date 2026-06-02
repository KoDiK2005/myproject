# myproject

Учебный REST API на Go с PostgreSQL.

## Стек

- Go + net/http
- PostgreSQL + sqlx
- Layered architecture: handler → service → repository

## Структура

```
cmd/myapp/        — точка входа
internal/
  config/         — конфиг из env
  handler/        — HTTP-хэндлеры
  service/        — бизнес-логика
  repository/     — работа с БД
  models/         — структуры данных
migrations/       — SQL-миграции
```

## Запуск

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/myproject?sslmode=disable"
export PORT=8080
go run ./cmd/myapp
```

## Эндпоинты

| Метод | URL     | Описание             |
|-------|---------|----------------------|
| POST  | /users  | Создать пользователя |
| GET   | /users  | Список пользователей |

### POST /users

```json
// запрос
{ "name": "Mark", "email": "mark@example.com" }

// ответ 201
{ "id": 1, "name": "Mark", "email": "mark@example.com" }
```

## Миграции

```bash
psql $DATABASE_URL -f migrations/001_create_users_table.sql
```
