# Docker Setup Guide

Руководство по запуску проекта с помощью Docker и Docker Compose.

## Требования

- Docker 20.10+
- Docker Compose 2.0+

## Быстрый старт

### 1. Настройка переменных окружения

Создайте файл `.env` на основе `.env.example`:

```bash
cp .env.example .env
```

Отредактируйте `.env` и укажите свои значения, особенно:
- `JWT_ACCESS_SECRET` - секретный ключ для JWT (минимум 32 символа)
- `JWT_REFRESH_SECRET` - секретный ключ для refresh token (минимум 32 символа)
- `DB_PASSWORD` - пароль для PostgreSQL

### 2. Запуск в production режиме

```bash
# Сборка и запуск всех сервисов
docker-compose up -d

# Или используя Makefile
make build
make up
```

### 3. Запуск в development режиме

```bash
# Используя docker-compose.dev.yml
docker-compose -f docker-compose.dev.yml up -d

# Или используя Makefile
make dev
```

## Сервисы

### PostgreSQL
- **Порт**: 5432 (по умолчанию)
- **База данных**: saas_platform
- **Пользователь**: postgres
- **Пароль**: из переменной `DB_PASSWORD`
- **Volume**: `postgres_data` - данные сохраняются между перезапусками

### Backend (Go)
- **Порт**: 8080 (по умолчанию)
- **Health check**: `http://localhost:8080/health`
- **API**: `http://localhost:8080/api/v1`

### Frontend (Vite/React)
- **Порт**: 3000 (по умолчанию)
- **URL**: `http://localhost:3000`
- **API URL**: настраивается через `VITE_API_URL`

## Полезные команды

### Просмотр логов

```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

### Остановка сервисов

```bash
docker-compose down

# С удалением volumes (удалит данные БД!)
docker-compose down -v
```

### Перезапуск сервисов

```bash
docker-compose restart

# Конкретный сервис
docker-compose restart backend
```

### Пересборка образов

```bash
# Обычная пересборка
docker-compose build

# Пересборка без кеша
docker-compose build --no-cache
```

### Просмотр статуса

```bash
docker-compose ps
```

## Volumes

### PostgreSQL Data
- **Volume**: `postgres_data`
- **Путь в контейнере**: `/var/lib/postgresql/data`
- **Назначение**: Хранение данных PostgreSQL

### Development Volumes
В development режиме используются bind mounts для hot reload:
- Backend: `./:/app` - весь код проекта
- Frontend: `./admin-panel:/app` - код frontend

## Сети

Все сервисы подключены к сети `saas-network` (или `saas-network-dev` в dev режиме), что позволяет им общаться по именам сервисов:
- Backend обращается к БД через `postgres:5432`
- Frontend может обращаться к backend через `backend:8080` (если настроен proxy)

## Troubleshooting

### Проблема: Backend не может подключиться к БД

**Решение**: Убедитесь, что:
1. PostgreSQL запущен и здоров (`docker-compose ps`)
2. В `.env` указан правильный `DB_HOST=postgres` (не localhost!)
3. Пароли совпадают в настройках postgres и backend

### Проблема: Frontend не может подключиться к API

**Решение**: 
1. Проверьте `VITE_API_URL` в `.env`
2. В production используйте `http://backend:8080/api/v1` (внутри Docker сети)
3. Или `http://localhost:8080/api/v1` (снаружи Docker)

### Проблема: Порт уже занят

**Решение**: Измените порты в `.env`:
- `SERVER_PORT=8081` для backend
- `FRONTEND_PORT=3001` для frontend
- `DB_PORT=5433` для postgres

### Очистка всего

```bash
# Остановить и удалить все контейнеры, volumes, сети
docker-compose down -v --remove-orphans
```

