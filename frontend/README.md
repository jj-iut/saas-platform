# SaaS Platform Frontend

Next.js frontend для SaaS-платформы.

## Технологии

- **Next.js 14** - React фреймворк
- **TypeScript** - типизация
- **Tailwind CSS** - стилизация

## Разработка

```bash
# Установка зависимостей
npm install

# Запуск dev сервера
npm run dev
```

Откройте [http://localhost:3000](http://localhost:3000) в браузере.

## Переменные окружения

Создайте файл `.env.local`:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Деплой на Vercel

1. Подключите репозиторий к Vercel
2. Настройте переменные окружения:
   - `NEXT_PUBLIC_API_URL` - URL вашего backend API
3. Деплой произойдет автоматически при push в main ветку

