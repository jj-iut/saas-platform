# Структура проекта

## Обзор

Проект состоит из двух основных частей:
- **Backend** (Go) - REST API сервер
- **Frontend** (Next.js) - веб-приложение

## Структура директорий

```
saas-platform/
├── frontend/                 # Next.js frontend приложение
│   ├── app/                  # Next.js App Router
│   │   ├── layout.tsx        # Root layout
│   │   ├── page.tsx          # Главная страница
│   │   ├── globals.css       # Глобальные стили
│   │   └── api/              # API routes (proxy к backend)
│   ├── public/               # Статические файлы
│   ├── package.json          # Node.js зависимости
│   ├── tsconfig.json         # TypeScript конфигурация
│   ├── next.config.js        # Next.js конфигурация
│   ├── tailwind.config.ts    # Tailwind CSS конфигурация
│   └── .env.example          # Пример переменных окружения
│
├── internal/                 # Backend (Go) - внутренние пакеты
│   ├── config/               # Конфигурация приложения
│   │   └── config.go         # Загрузка конфигурации из env
│   ├── database/             # Работа с базой данных
│   │   ├── connection.go     # Подключение к PostgreSQL
│   │   └── migrations.go     # Миграции базы данных
│   ├── handlers/             # HTTP обработчики
│   │   └── health.go         # Health check endpoint
│   ├── router/               # Маршрутизация
│   │   └── router.go         # Настройка роутера и middleware
│   └── modules/              # Бизнес-модули (для будущего расширения)
│       └── ...               # Модули будут добавлены по необходимости
│
├── main.go                   # Точка входа backend приложения
├── go.mod                    # Go модули и зависимости
├── go.sum                    # Checksums зависимостей
│
├── docker-compose.yml        # Docker Compose для production
├── docker-compose.dev.yml    # Docker Compose для development
├── Dockerfile.backend        # Dockerfile для backend
├── Dockerfile.frontend       # Dockerfile для frontend (production)
├── Dockerfile.frontend.dev   # Dockerfile для frontend (development)
│
├── vercel.json               # Конфигурация Vercel для frontend
├── render.yaml               # Конфигурация Render для backend
│
├── README.md                 # Основная документация
├── DEPLOY.md                 # Руководство по деплою
├── README.DOCKER.md          # Docker документация
└── PROJECT_STRUCTURE.md      # Этот файл
```

## Backend (Go)

### Архитектура

Backend следует принципам **Modular Monolith** архитектуры:

- **Handler** → обработка HTTP запросов/ответов
- **Service** → бизнес-логика и валидация
- **Repository** → работа с базой данных

### Основные компоненты

1. **config/** - Загрузка конфигурации из переменных окружения
2. **database/** - Подключение к PostgreSQL и миграции
3. **handlers/** - HTTP обработчики для endpoints
4. **router/** - Настройка маршрутов и middleware (CORS, Request ID)

### API Endpoints

- `GET /health` - Health check (проверка состояния сервера и БД)
- `GET /api/v1/*` - API endpoints (будут добавлены в модулях)

## Frontend (Next.js)

### Технологии

- **Next.js 14** с App Router
- **TypeScript** для типизации
- **Tailwind CSS** для стилизации

### Структура

- **app/** - Next.js App Router директория
  - `layout.tsx` - Root layout компонент
  - `page.tsx` - Главная страница
  - `api/` - API routes (используются как proxy к backend)

### Переменные окружения

- `NEXT_PUBLIC_API_URL` - URL backend API (например: `http://localhost:8080/api/v1`)

## Деплой

### Backend → Render

1. Создайте PostgreSQL database на Render
2. Создайте Web Service для Go приложения
3. Настройте environment variables
4. Render автоматически соберет и задеплоит backend

Подробности в [DEPLOY.md](./DEPLOY.md)

### Frontend → Vercel

1. Подключите GitHub репозиторий к Vercel
2. Укажите root directory: `frontend`
3. Настройте `NEXT_PUBLIC_API_URL`
4. Vercel автоматически соберет и задеплоит frontend

Подробности в [DEPLOY.md](./DEPLOY.md)

## Локальная разработка

### С Docker Compose

```bash
# Development режим
docker-compose -f docker-compose.dev.yml up

# Production режим
docker-compose up
```

### Без Docker

**Backend:**
```bash
# Установите PostgreSQL локально
# Создайте базу данных
createdb saas_platform

# Создайте .env файл
cp .env.example .env

# Запустите
go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## Переменные окружения

См. `.env.example` в корне проекта для backend и `frontend/.env.example` для frontend.

## Дальнейшее развитие

### Добавление нового модуля

1. Создайте директорию в `internal/modules/your-module/`
2. Создайте файлы:
   - `model.go` - модели данных
   - `repository.go` - работа с БД
   - `service.go` - бизнес-логика
   - `handler.go` - HTTP handlers
3. Зарегистрируйте маршруты в `internal/router/router.go`

### Добавление нового API endpoint

1. Создайте handler в соответствующем модуле
2. Добавьте маршрут в `internal/router/router.go`
3. Используйте паттерн Handler → Service → Repository

