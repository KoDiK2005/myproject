# myproject

Учебный REST API на Go с PostgreSQL.

## Стек

- Go + Gin
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

## Авторизация

PUT и DELETE защищены JWT. Сначала получи токен:

```
POST /auth/login → { "token": "..." }
```

Затем передавай в заголовке:

```
Authorization: Bearer <token>
```

## Эндпоинты

| Метод  | URL                | Авторизация | Описание              |
|--------|--------------------|-------------|-----------------------|
| POST   | /auth/login        | —           | Получить JWT токен    |
| POST   | /api/v1/users      | —           | Создать пользователя  |
| GET    | /api/v1/users      | —           | Список пользователей  |
| GET    | /api/v1/users/{id} | —           | Получить по ID        |
| PUT    | /api/v1/users/{id} | JWT         | Обновить пользователя |
| DELETE | /api/v1/users/{id} | JWT         | Удалить пользователя  |
| GET    | /api/v1/posts      | —           | Список постов         |
| GET    | /api/v1/posts/{id} | —           | Получить пост по ID   |
| POST   | /api/v1/posts      | JWT         | Создать пост          |
| DELETE | /api/v1/posts/{id} | JWT         | Удалить пост          |

### POST /auth/login

```json
// запрос
{ "email": "mark@example.com", "password": "secret123" }

// ответ 200
{ "token": "eyJhbGci..." }
```

### POST /api/v1/users

```json
// запрос
{ "name": "Mark", "email": "mark@example.com", "password": "secret123" }

// ответ 201
{ "id": 1, "name": "Mark", "email": "mark@example.com" }
```

### PUT /api/v1/users/{id}

```json
// запрос
{ "name": "Mark Updated", "email": "mark2@example.com" }

// ответ 200
{ "id": 1, "name": "Mark Updated", "email": "mark2@example.com" }
```
