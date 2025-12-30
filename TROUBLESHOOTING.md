# Troubleshooting Guide

## Проблема: 404 ошибка при попытке входа

Если вы получаете 404 ошибку при попытке войти, проверьте следующее:

### 1. Проверьте, что Backend запущен и доступен

Backend должен быть доступен по URL, указанному в `NEXT_PUBLIC_API_URL`.

**Для локальной разработки:**
- Backend должен быть доступен на `http://localhost:8080`
- Проверьте: откройте в браузере `http://localhost:8080/health`
- Должен вернуться JSON: `{"status":"healthy",...}`

**Для production (Vercel + Render):**
- Backend должен быть деплоен на Render и доступен по URL вида: `https://saas-platform-backend.onrender.com`
- Проверьте: откройте в браузере `https://your-backend-url.onrender.com/health`
- Должен вернуться JSON: `{"status":"healthy",...}`

### 2. Проверьте переменную окружения NEXT_PUBLIC_API_URL

**На Vercel:**

1. Зайдите в Dashboard вашего проекта
2. Settings → Environment Variables
3. Убедитесь, что есть переменная `NEXT_PUBLIC_API_URL`
4. Значение должно быть: `https://your-backend-url.onrender.com/api/v1`
   - **ВАЖНО:** НЕ добавляйте слэш в конце!
   - Правильно: `https://saas-platform-backend.onrender.com/api/v1`
   - Неправильно: `https://saas-platform-backend.onrender.com/api/v1/`

5. После изменения переменной окружения:
   - Перейдите в Deployments
   - Нажмите на последний deployment
   - Нажмите "Redeploy" для применения изменений

### 3. Проверьте логи Backend на Render

1. Зайдите в Render Dashboard
2. Выберите ваш backend service
3. Перейдите на вкладку "Logs"
4. Проверьте, что сервер запустился успешно
5. Ищите ошибки подключения к базе данных

### 4. Проверьте CORS настройки

Backend уже настроен с CORS, который разрешает запросы от любого origin. Если у вас все еще проблемы:

1. Убедитесь, что backend действительно запущен
2. Проверьте, что frontend отправляет запросы на правильный URL

### 5. Проверьте Network tab в браузере

1. Откройте Developer Tools (F12)
2. Перейдите на вкладку Network
3. Попробуйте войти снова
4. Посмотрите на запрос к `/auth/login`:
   - Проверьте URL запроса (должен быть полный URL с вашим backend)
   - Проверьте статус код (404 = endpoint не найден, 500 = ошибка сервера)
   - Посмотрите на Response (если есть) для деталей ошибки

### Типичные проблемы и решения

**Проблема:** 404 Not Found
- **Причина:** Backend не запущен или неправильный URL
- **Решение:** Проверьте, что backend доступен и `NEXT_PUBLIC_API_URL` указывает на правильный URL

**Проблема:** CORS ошибка
- **Причина:** Backend не разрешает запросы с вашего домена
- **Решение:** Проверьте, что backend запущен и CORS настроен (уже настроен в коде)

**Проблема:** Connection refused / Network error
- **Причина:** Backend недоступен по указанному URL
- **Решение:** Проверьте, что backend деплоится на Render и доступен

**Проблема:** 500 Internal Server Error
- **Причина:** Ошибка на стороне backend (обычно проблема с БД)
- **Решение:** Проверьте логи backend на Render, убедитесь, что БД настроена правильно

### Пример правильной настройки

**На Vercel (Environment Variables):**
```
NEXT_PUBLIC_API_URL=https://saas-platform-backend.onrender.com/api/v1
```

**На Render:**
- Backend service должен быть доступен по URL: `https://saas-platform-backend.onrender.com`
- Health check должен работать: `https://saas-platform-backend.onrender.com/health`
- База данных должна быть подключена и доступна

### Тестирование

1. **Тест backend:**
   ```bash
   curl https://your-backend-url.onrender.com/health
   ```
   Должен вернуть: `{"status":"healthy",...}`

2. **Тест login endpoint:**
   ```bash
   curl -X POST https://your-backend-url.onrender.com/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@saas-platform.com","password":"Admin123!"}'
   ```
   Должен вернуть JSON с токенами и данными пользователя

Если эти команды работают, но frontend все еще получает 404, проблема в `NEXT_PUBLIC_API_URL` на Vercel.

