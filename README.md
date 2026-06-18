# myproject — REST API на Go + Gin + PostgreSQL

Учебный проект. Полноценный REST API с авторизацией, пагинацией, поиском, комментариями, лайками и загрузкой файлов. Плюс React фронтенд.

## Стек

**Backend**
- **Go** + **Gin** — HTTP фреймворк
- **PostgreSQL** + **sqlx** — база данных
- **JWT** — авторизация (access 15мин + refresh 7 дней)
- **zerolog** — структурированные логи
- **Prometheus** — метрики
- **Swagger** — документация (`/swagger/index.html`)

**Frontend**
- **React 19** + **Vite** — UI
- **react-router-dom v7** — роутинг

**Инфраструктура**
- **Docker Compose** — запуск одной командой
- **GitHub Actions** — CI (unit тесты на каждый пуш)

## Быстрый старт

```bash
docker compose up --build
```

- API: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger/index.html`
- Frontend: `http://localhost:5173` (отдельно через `npm run dev` в папке `frontend/`)

## Переменные окружения

| Переменная   | Описание                  | Пример                                                    |
|--------------|---------------------------|-----------------------------------------------------------|
| DATABASE_URL | Строка подключения к БД   | postgres://postgres:password@db:5432/mydb?sslmode=disable |
| JWT_SECRET   | Секрет для JWT токенов    | supersecret                                               |
| PORT         | Порт сервера              | 8080                                                      |

## Эндпоинты

### Auth
| Метод | URL           | Описание                     |
|-------|---------------|------------------------------|
| POST  | /auth/login   | Логин, возвращает оба токена |
| POST  | /auth/refresh | Обновить access token        |
| POST  | /auth/logout  | Отозвать refresh token       |

### Users
| Метод  | URL                      | Auth | Описание                             |
|--------|--------------------------|------|--------------------------------------|
| GET    | /api/v1/users            | —    | Список юзеров (пагинация)            |
| POST   | /api/v1/users            | —    | Создать юзера                        |
| GET    | /api/v1/users/:id        | —    | Получить юзера                       |
| PUT    | /api/v1/users/:id        | ✓    | Обновить себя                        |
| DELETE | /api/v1/users/:id        | ✓    | Удалить себя                         |
| POST   | /api/v1/users/:id/avatar | ✓    | Загрузить аватар (jpg/png, макс 2MB) |

### Posts
| Метод  | URL               | Auth | Описание                                |
|--------|-------------------|------|-----------------------------------------|
| GET    | /api/v1/posts     | —    | Список постов (`?page=&limit=&search=`) |
| POST   | /api/v1/posts     | ✓    | Создать пост                            |
| GET    | /api/v1/posts/:id | —    | Получить пост                           |
| PUT    | /api/v1/posts/:id | ✓    | Обновить свой пост                      |
| DELETE | /api/v1/posts/:id | ✓    | Удалить свой пост                       |

### Comments
| Метод  | URL                        | Auth | Описание                 |
|--------|----------------------------|------|--------------------------|
| GET    | /api/v1/posts/:id/comments | —    | Комментарии к посту      |
| POST   | /api/v1/posts/:id/comments | ✓    | Добавить комментарий     |
| DELETE | /api/v1/comments/:id       | ✓    | Удалить свой комментарий |

### Likes
| Метод  | URL                     | Auth | Описание      |
|--------|-------------------------|------|---------------|
| GET    | /api/v1/posts/:id/likes | —    | Кол-во лайков |
| POST   | /api/v1/posts/:id/like  | ✓    | Лайкнуть      |
| DELETE | /api/v1/posts/:id/like  | ✓    | Убрать лайк   |

### Системные
| Метод | URL           | Описание               |
|-------|---------------|------------------------|
| GET   | /health       | Проверка работоспособности |
| GET   | /metrics      | Prometheus метрики     |
| GET   | /swagger/*any | Swagger UI             |

## Тесты

```bash
go test ./...
```

Покрыто: сервисный слой (PostService, UserService) и HTTP-хендлеры (посты, авторизация) через `httptest`.

## Миграции

```bash
for f in migrations/*.sql; do psql $DATABASE_URL -f $f; done
```

Или вручную по порядку: `001` → `009`.

## Архитектура

```
cmd/myapp/       — точка входа, сборка зависимостей, graceful shutdown
internal/
  config/        — конфиг из env
  handler/       — HTTP хендлеры + middleware (JWT, rate limit, метрики)
  service/       — бизнес-логика, sentinel errors
  repository/    — работа с БД (sqlx)
  models/        — структуры данных
  logger/        — zerolog
migrations/      — SQL миграции (001–009)
docs/            — сгенерированный Swagger
frontend/        — React приложение (Vite)
```

## Оптимизации

- **Индексы** на `posts.user_id`, `comments.post_id`, `likes.post_id`
- **pg_trgm** + GIN индексы для быстрого ILIKE-поиска
- **Window function** `COUNT(*) OVER()` — пагинация за один запрос вместо двух
- **Connection pool**: `SetMaxOpenConns(25)`, `SetMaxIdleConns(10)`
