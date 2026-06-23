# myproject — социальная сеть на Go + React

Учебный проект. Полноценная социальная сеть: посты с приватностью, система дружбы, реалтайм-мессенджер через WebSocket, комментарии, лайки, поиск.

## Стек

**Backend**
- **Go** + **Gin** — HTTP фреймворк
- **PostgreSQL** + **sqlx** — база данных
- **JWT** — авторизация (access 15мин + refresh 7 дней с ротацией)
- **gorilla/websocket** — реалтайм чат
- **zerolog** — структурированные логи
- **Prometheus** — метрики
- **Swagger** — документация (`/swagger/index.html`)

**Frontend**
- **React 19** + **Vite**
- **react-router-dom v7** — роутинг
- WebSocket API — реалтайм сообщения

**Инфраструктура**
- **Docker Compose** — запуск одной командой
- **GitHub Actions** — CI (unit тесты на каждый пуш)

## Быстрый старт

```bash
docker compose up --build
```

Применить миграции:
```bash
for f in migrations/*.sql; do psql $DATABASE_URL -f $f; done
```

- API: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger/index.html`
- Frontend: `http://localhost:5173` (отдельно: `npm run dev` в `frontend/`)

## Переменные окружения

| Переменная      | Описание                              | Пример                                                    |
|-----------------|----------------------------------------|-----------------------------------------------------------|
| DATABASE_URL    | Строка подключения к БД               | postgres://postgres:password@db:5432/mydb?sslmode=disable |
| JWT_SECRET      | Секрет для JWT токенов                | supersecret                                               |
| PORT            | Порт сервера                          | 8080                                                      |
| ALLOWED_ORIGINS | CORS + WebSocket origin (через запятую) | http://localhost:3000,http://localhost:5173              |

## Эндпоинты

### Auth
| Метод | URL           | Описание                      |
|-------|---------------|-------------------------------|
| POST  | /auth/login   | Логин, возвращает оба токена  |
| POST  | /auth/refresh | Обновить токены (ротация)     |
| POST  | /auth/logout  | Отозвать refresh token        |

### WebSocket
| URL                  | Auth | Описание                              |
|----------------------|------|---------------------------------------|
| GET /ws?token=...    | ✓    | WS соединение для входящих сообщений  |

### Users
| Метод  | URL                      | Auth | Описание                             |
|--------|--------------------------|------|--------------------------------------|
| GET    | /api/v1/users            | —    | Список / поиск (`?search=имя`)       |
| POST   | /api/v1/users            | —    | Регистрация                          |
| GET    | /api/v1/users/:id        | —    | Профиль пользователя                 |
| PUT    | /api/v1/users/:id        | ✓    | Обновить себя                        |
| DELETE | /api/v1/users/:id        | ✓    | Удалить себя                         |
| POST   | /api/v1/users/:id/avatar | ✓    | Загрузить аватар (jpg/png, макс 2MB) |
| GET    | /api/v1/users/:id/posts  | —    | Посты пользователя (с учётом дружбы) |

### Posts
| Метод  | URL               | Auth    | Описание                                               |
|--------|-------------------|---------|--------------------------------------------------------|
| GET    | /api/v1/posts     | optional| Лента: гость — публичные, авторизованный — персональная|
| POST   | /api/v1/posts     | ✓       | Создать пост (`visibility: public\|friends`)           |
| GET    | /api/v1/posts/:id | —       | Получить пост                                          |
| PUT    | /api/v1/posts/:id | ✓       | Обновить свой пост                                     |
| DELETE | /api/v1/posts/:id | ✓       | Удалить свой пост                                      |

### Дружба (только взаимная, без подписок)
| Метод  | URL                               | Auth | Описание               |
|--------|-----------------------------------|------|------------------------|
| POST   | /api/v1/friends/request/:id       | ✓    | Отправить заявку       |
| POST   | /api/v1/friends/accept/:id        | ✓    | Принять заявку         |
| POST   | /api/v1/friends/reject/:id        | ✓    | Отклонить заявку       |
| DELETE | /api/v1/friends/:id               | ✓    | Удалить из друзей      |
| GET    | /api/v1/friends                   | ✓    | Мой список друзей      |
| GET    | /api/v1/friends/requests/incoming | ✓    | Входящие заявки        |
| GET    | /api/v1/friends/requests/outgoing | ✓    | Исходящие заявки       |
| GET    | /api/v1/friends/status/:id        | ✓    | Статус с конкретным    |

### Сообщения
| Метод | URL                  | Auth | Описание                          |
|-------|----------------------|------|-----------------------------------|
| GET   | /api/v1/messages     | ✓    | Список переписок                  |
| GET   | /api/v1/messages/:id | ✓    | История с пользователем           |
| POST  | /api/v1/messages/:id | ✓    | Отправить сообщение (только другу)|

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
| Метод | URL           | Описание                   |
|-------|---------------|----------------------------|
| GET   | /health       | Проверка (503 если БД упала)|
| GET   | /metrics      | Prometheus метрики         |
| GET   | /swagger/*any | Swagger UI                 |

## Фронт — страницы

| Путь           | Доступ | Описание                                    |
|----------------|--------|---------------------------------------------|
| /              | ✓      | Лента постов, создание поста                |
| /people        | ✓      | Поиск людей по имени                        |
| /friends       | ✓      | Друзья, входящие/исходящие заявки           |
| /messages      | ✓      | Список переписок                            |
| /messages/:id  | ✓      | Чат (реалтайм WS)                           |
| /profile       | ✓      | Свой профиль, редактирование, аватар        |
| /users/:id     | —      | Профиль другого юзера + кнопка дружбы      |
| /posts/:id     | ✓      | Пост с комментариями и лайками              |

## Тесты

```bash
go test ./...
```

Покрыто: PostService, UserService, CommentService, LikeService + HTTP-хендлеры (посты, авторизация).

## Миграции

```
001 — users
002 — posts
003 — comments
004 — likes
005 — refresh_tokens
006 — avatar поле
007 — rate_limit (не используется)
008 — индексы на posts/comments/likes
009 — pg_trgm + GIN индексы для поиска
010 — (reserved)
011 — friendships
012 — posts.visibility (public|friends)
013 — messages + индексы
```

## Архитектура

```
cmd/myapp/       — точка входа, DI, graceful shutdown
internal/
  config/        — конфиг из env
  handler/       — HTTP хендлеры + middleware (JWT, rate limit, CORS, метрики)
  service/       — бизнес-логика, sentinel errors
  repository/    — работа с БД (sqlx, context timeouts)
  models/        — структуры данных
  ws/            — WebSocket hub (горутина-диспетчер)
  logger/        — zerolog
migrations/      — SQL миграции (001–013)
docs/            — сгенерированный Swagger
frontend/        — React приложение (Vite)
  src/
    api/         — fetch-обёртки (auth, posts, friends, messages...)
    pages/       — страницы
    components/  — PostCard, Toast
    hooks/       — useFriendBadge
```

## Оптимизации

- **Connection pool**: `SetMaxOpenConns(25)`, `SetMaxIdleConns(10)`, `SetConnMaxLifetime(5min)`
- **Context timeouts** (5s) во всех DB-запросах
- **pg_trgm** + GIN индексы для ILIKE-поиска
- **Window function** `COUNT(*) OVER()` — пагинация за один запрос
- **Персональная лента** — один SQL с подзапросом по дружбам вместо кода
- **WebSocket hub** — горутина на весь сервер, не поток на соединение
- **apiFetch** — авто-рефреш токена на фронте, юзер не видит истечения сессии
- **WS reconnect** — чат переподключается при разрыве соединения с экспоненциальной задержкой
