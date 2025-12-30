# Быстрое исправление 404 ошибки для /api/v1/auth/login

## Проблема
Health endpoint работает (`/health` возвращает 200), но `/api/v1/auth/login` возвращает 404.

## Причина
Backend на Render не был перезапущен после добавления новых модулей аутентификации, или код не был правильно задеплоен.

## Решение

### Вариант 1: Перезапуск через Render Dashboard (рекомендуется)

1. Зайдите в Render Dashboard: https://dashboard.render.com/
2. Найдите ваш backend service (`saas-platform-backend`)
3. Нажмите на три точки (...) → **"Manual Deploy"** → **"Deploy latest commit"**
4. Или просто нажмите кнопку **"Manual Deploy"** если она есть
5. Дождитесь завершения деплоя (обычно 2-3 минуты)
6. Проверьте логи, убедитесь что нет ошибок компиляции

### Вариант 2: Проверка через Render Logs

1. В Render Dashboard → ваш backend service
2. Перейдите на вкладку **"Logs"**
3. Проверьте последние логи:
   - Должны быть сообщения типа "Server starting on port 8080"
   - Не должно быть ошибок типа "undefined: authModule" или подобных
   - Если есть ошибки компиляции, backend не запустится

### Вариант 3: Принудительный редеплой через Git

Если Render использует автоматический деплой из GitHub:

1. Убедитесь, что все изменения закоммичены и запушены в GitHub
2. Создайте пустой коммит для принудительного деплоя:
   ```bash
   git commit --allow-empty -m "Trigger redeploy"
   git push origin main
   ```
3. Render автоматически задеплоит новую версию

### Проверка после редеплоя

После редеплоя проверьте:

1. Health endpoint: `https://ваш-backend-url.onrender.com/health` - должен работать
2. Login endpoint через curl или браузер:
   ```bash
   curl -X POST https://ваш-backend-url.onrender.com/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@saas-platform.com","password":"Admin123!"}'
   ```

Должен вернуться JSON с токенами, а не 404.

### Типичные ошибки в логах

Если видите ошибки типа:
- `undefined: authModule` - код не компилируется, проверьте что все файлы на месте
- `failed to connect to database` - проблема с подключением к БД
- `no such table: users` - миграции не применились

### Если ничего не помогает

1. Проверьте, что на Render действительно используется последний коммит из GitHub
2. Убедитесь, что все файлы из `internal/modules/auth/` и `internal/modules/restaurants/` закоммичены
3. Попробуйте пересоздать service на Render с нуля (если это возможно)

