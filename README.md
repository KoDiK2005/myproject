# myproject — REST API на Go + Gin + PostgreSQL

Учебный проект. Полноценный REST API с авторизацией, пагинацией, поиском, комментариями, лайками и загрузкой файлов.

## Стек

- **Go 1.22** + **Gin** — HTTP фреймворк
- **PostgreSQL** + **sqlx** — база данных
- **JWT** — авторизация (access 15мин + refresh 7 дней)
- **zerolog** — структурированные логи
- **Prometheus** — метрики
- **Swagger** — документация (`/swagger/index.html`)
- **Docker Compose** — запуск одной командой
- **GitHub Actions** — CI (unit + integration тесты)

## Запуск

```bash
docker compose up --build
```

API будет доступен на `http://localhost:8080`

## Переменные окружения

| Переменная    | Описание              | Пример                                                |
|---------------|-----------------------|-------------------------------------------------------|
| DATABASE_URL  | Строка подключения к БД | postgres://postgres:password@db:5432/mydb?sslmode=disable |
| JWT_SECRET    | Секрет для JWT токенов | supersecret                                          |
| PORT          | Порт сервера          | 8080                                                  |

## Эндпоинты

### Auth
| Метод | URL              | Описание                    |
|-------|------------------|-----------------------------|
| POST  | /auth/login      | Логин, возвращает оба токена |
| POST  | /auth/refresh    | Обновить access token       |
| POST  | /auth/logout     | Отозвать refresh token      |

### Users
| Метод  | URL                    | Auth | Описание              |
|--------|------------------------|------|-----------------------|
| GET    | /api/v1/users          | —    | Список юзеров (пагинация) |
| POST   | /api/v1/users          | —    | Создать юзера         |
| GET    | /api/v1/users/:id      | —    | Получить юзера        |
| PUT    | /api/v1/users/:id      | ✓    | Обновить себя         |
| DELETE | /api/v1/users/:id      | ✓    | Удалить себя          |
| POST   | /api/v1/users/:id/avatar | ✓  | Загрузить аватар (jpg/png, макс 2MB) |

### Posts
| Метод  | URL                | Auth | Описание                        |
|--------|--------------------|------|---------------------------------|
| GET    | /api/v1/posts      | —    | Список постов (`?page=&limit=&search=`) |
| POST   | /api/v1/posts      | ✓    | Создать пост                    |
| GET    | /api/v1/posts/:id  | —    | Получить пост                   |
| PUT    | /api/v1/posts/:id  | ✓    | Обновить свой пост              |
| DELETE | /api/v1/posts/:id  | ✓    | Удалить свой пост               |

### Comments
| Метод  | URL                         | Auth | Описание              |
|--------|-----------------------------|------|-----------------------|
| GET    | /api/v1/posts/:id/comments  | —    | Комментарии к посту   |
| POST   | /api/v1/posts/:id/comments  | ✓    | Добавить комментарий  |
| DELETE | /api/v1/comments/:id        | ✓    | Удалить свой комментарий |

### Likes
| Метод  | URL                      | Auth | Описание        |
|--------|--------------------------|------|-----------------|
| GET    | /api/v1/posts/:id/likes  | —    | Кол-во лайков   |
| POST   | /api/v1/posts/:id/like   | ✓    | Лайкнуть        |
| DELETE | /api/v1/posts/:id/like   | ✓    | Убрать лайк     |

### Системные
| Метод | URL              | Описание               |
|-------|------------------|------------------------|
| GET   | /health          | Проверка работоспособности |
| GET   | /metrics         | Prometheus метрики     |
| GET   | /swagger/*any    | Swagger UI             |

## Запуск тестов

```bash
# юнит-тесты
go test ./...

# интеграционные (нужен PostgreSQL)
DATABASE_URL=postgres://... go test -tags integration ./internal/repository/... -v
```

## Миграции

Применяются вручную:
```bash
psql $DATABASE_URL -f migrations/001_create_users_table.sql
psql $DATABASE_URL -f migrations/002_add_password.sql
psql $DATABASE_URL -f migrations/003_create_posts.sql
psql $DATABASE_URL -f migrations/004_create_refresh_tokens.sql
psql $DATABASE_URL -f migrations/005_create_comments.sql
psql $DATABASE_URL -f migrations/006_create_likes.sql
psql $DATABASE_URL -f migrations/007_add_avatar_to_users.sql
```

## Архитектура

```
cmd/myapp/         — точка входа
internal/
  config/          — конфиг из env
  handler/         — HTTP хендлеры + middleware
  service/         — бизнес-логика
  repository/      — работа с БД
  models/          — структуры данных
  logger/          — zerolog
migrations/        — SQL миграции
docs/              — сгенерированный Swagger
```
