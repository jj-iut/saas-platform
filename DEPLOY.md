# Деплой на Vercel + Render

Руководство по деплою SaaS Platform на Vercel (frontend) и Render (backend).

## Архитектура деплоя

- **Frontend (Next.js)**: Vercel
- **Backend (Go)**: Render
- **Database (PostgreSQL)**: Render PostgreSQL

## Деплой Backend на Render

### 1. Подготовка

1. Убедитесь, что в репозитории есть файл `render.yaml`
2. Проверьте, что все зависимости указаны в `go.mod`

### 2. Создание сервиса на Render

1. Зайдите на [Render Dashboard](https://dashboard.render.com/)
2. Нажмите "New +" → "Blueprint"
3. Подключите ваш GitHub репозиторий
4. Render автоматически обнаружит `render.yaml` и создаст сервисы

Или создайте вручную:

1. **Создайте PostgreSQL Database:**
   - Нажмите "New +" → "PostgreSQL"
   - Выберите план (Free для начала)
   - Назовите: `saas-platform-db`
   - Запишите credentials

2. **Создайте Web Service для Backend:**
   - Нажмите "New +" → "Web Service"
   - Подключите репозиторий
   - Настройки:
     - **Name**: `saas-platform-backend`
     - **Environment**: `Go`
     - **Region**: выберите ближайший
     - **Branch**: `main`
     - **Root Directory**: `.` (root)
     - **Build Command**: `go mod download && go build -o main .`
     - **Start Command**: `./main`

3. **Настройте Environment Variables:**
   - `SERVER_PORT`: `8080`
   - `ENVIRONMENT`: `production`
   - `DB_HOST`: из PostgreSQL service (Internal Database URL)
   - `DB_PORT`: из PostgreSQL service
   - `DB_USER`: из PostgreSQL service
   - `DB_PASSWORD`: из PostgreSQL service
   - `DB_NAME`: из PostgreSQL service
   - `DB_SSLMODE`: `require`
   - `JWT_ACCESS_SECRET`: сгенерируйте случайную строку (минимум 32 символа)
   - `JWT_REFRESH_SECRET`: сгенерируйте случайную строку (минимум 32 символа)
   - `JWT_ACCESS_TTL`: `15m`
   - `JWT_REFRESH_TTL`: `168h`

4. **Используйте Internal Database URL:**
   Render предоставляет Internal Database URL, который должен использоваться для подключения внутри сети Render.

### 3. Получение URL Backend

После деплоя Render предоставит URL вида: `https://saas-platform-backend.onrender.com`

Запишите этот URL - он понадобится для настройки frontend.

## Деплой Frontend на Vercel

### 1. Подготовка

1. Убедитесь, что в репозитории есть файл `vercel.json`
2. Проверьте структуру frontend директории

### 2. Создание проекта на Vercel

1. Зайдите на [Vercel Dashboard](https://vercel.com/dashboard)
2. Нажмите "Add New..." → "Project"
3. Подключите ваш GitHub репозиторий
4. Выберите репозиторий `saas-platform`

### 3. Настройка проекта

1. **Framework Preset**: Vercel должен автоматически определить Next.js
2. **Root Directory**: выберите `frontend`
3. **Build Command**: `npm run build` (уже в package.json)
4. **Output Directory**: `.next` (по умолчанию для Next.js)

### 4. Environment Variables

Добавьте переменные окружения:

- `NEXT_PUBLIC_API_URL`: URL вашего backend на Render
  - Например: `https://saas-platform-backend.onrender.com/api/v1`

**Важно**: Переменные с префиксом `NEXT_PUBLIC_` доступны в браузере.

### 5. Деплой

1. Нажмите "Deploy"
2. Vercel автоматически соберет и задеплоит приложение
3. После успешного деплоя вы получите URL вида: `https://saas-platform.vercel.app`

## Проверка деплоя

### Backend Health Check

```bash
curl https://saas-platform-backend.onrender.com/health
```

Должен вернуть:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "checks": {
    "database": "healthy"
  }
}
```

### Frontend Health Check

Откройте в браузере: `https://saas-platform.vercel.app/api/health`

Должен вернуть данные от backend.

## Обновление переменных окружения

### Backend (Render)

1. Зайдите в Dashboard → ваш Web Service
2. Settings → Environment
3. Добавьте/измените переменные
4. Нажмите "Save Changes"
5. Сервис автоматически перезапустится

### Frontend (Vercel)

1. Зайдите в Dashboard → ваш Project
2. Settings → Environment Variables
3. Добавьте/измените переменные
4. Нажмите "Save"
5. Запустите новый деплой (или он произойдет автоматически при следующем push)

## Полезные команды

### Проверка логов

**Render:**
- Dashboard → ваш сервис → Logs

**Vercel:**
- Dashboard → ваш проект → Deployments → выберите deployment → Logs

### Переменные окружения для локальной разработки

Создайте `.env` в корне проекта (для backend) и `.env.local` в `frontend/` (для frontend).

Смотрите `.env.example` для примера.

## Troubleshooting

### Backend не может подключиться к БД

1. Проверьте, что используете Internal Database URL от Render
2. Убедитесь, что `DB_SSLMODE=require`
3. Проверьте логи в Render Dashboard

### Frontend не может подключиться к Backend

1. Проверьте `NEXT_PUBLIC_API_URL` в Vercel
2. Убедитесь, что backend доступен (проверьте health endpoint)
3. Проверьте CORS настройки в backend (должны разрешать запросы с Vercel домена)

### CORS ошибки

Backend должен разрешать запросы с вашего Vercel домена. Обновите CORS middleware в `internal/router/router.go` если нужно.

## Мониторинг

### Render

- Автоматические health checks
- Логи доступны в реальном времени
- Metrics в Dashboard

### Vercel

- Analytics доступны в Dashboard
- Логи каждого deployment
- Performance metrics

## Стоимость

- **Render Free tier**: 
  - Web Service: спит после 15 минут неактивности
  - PostgreSQL: до 90 дней, ограниченное хранилище
- **Vercel Free tier**: 
  - Достаточно для большинства проектов
  - Ограничения на bandwidth и build time

Для production рекомендуется использовать paid планы.
