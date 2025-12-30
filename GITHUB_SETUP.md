# Инструкция по добавлению проекта в GitHub

## Шаг 1: Установите Git (если не установлен)

### Windows:
1. Скачайте Git с https://git-scm.com/download/win
2. Установите с настройками по умолчанию
3. Перезапустите терминал

### Проверка установки:
```bash
git --version
```

## Шаг 2: Настройте Git (первый раз)

```bash
git config --global user.name "Ваше Имя"
git config --global user.email "your.email@example.com"
```

## Шаг 3: Создайте репозиторий на GitHub

1. Перейдите на https://github.com
2. Войдите в свой аккаунт
3. Нажмите кнопку "+" в правом верхнем углу
4. Выберите "New repository"
5. Заполните:
   - **Repository name**: `saas-platform` (или другое имя)
   - **Description**: "SaaS Platform - Restaurant Management System"
   - **Visibility**: Public или Private (на ваш выбор)
   - **НЕ** создавайте README, .gitignore или лицензию (мы уже создали их)
6. Нажмите "Create repository"

## Шаг 4: Инициализируйте Git репозиторий локально

Откройте терминал в корневой директории проекта и выполните:

```bash
# Инициализация репозитория
git init

# Добавление всех файлов
git add .

# Создание первого коммита
git commit -m "Initial commit: SaaS Platform with backend, frontend and Docker setup"

# Добавление remote репозитория (замените YOUR_USERNAME на ваш GitHub username)
git remote add origin https://github.com/YOUR_USERNAME/saas-platform.git

# Переименование основной ветки в main (если нужно)
git branch -M main

# Отправка кода на GitHub
git push -u origin main
```

## Шаг 5: Альтернативный способ через GitHub CLI

Если у вас установлен GitHub CLI:

```bash
# Авторизация
gh auth login

# Создание репозитория и push
gh repo create saas-platform --public --source=. --remote=origin --push
```

## Шаг 6: Проверка

Перейдите на https://github.com/YOUR_USERNAME/saas-platform и убедитесь, что все файлы загружены.

## Дополнительные команды Git

### Просмотр статуса:
```bash
git status
```

### Просмотр изменений:
```bash
git diff
```

### Добавление изменений:
```bash
git add .
# или конкретный файл
git add filename
```

### Создание коммита:
```bash
git commit -m "Описание изменений"
```

### Отправка изменений:
```bash
git push
```

### Получение изменений:
```bash
git pull
```

### Просмотр истории:
```bash
git log
```

## Структура .gitignore

Проект уже содержит `.gitignore` файл, который исключает:
- Скомпилированные файлы
- Зависимости (node_modules, vendor)
- Файлы окружения (.env)
- IDE файлы
- Логи и временные файлы

## Важные файлы для GitHub

Убедитесь, что следующие файлы добавлены:
- ✅ README.md
- ✅ .gitignore
- ✅ docker-compose.yml
- ✅ Dockerfile.backend
- ✅ Все исходные файлы проекта

## Troubleshooting

### Ошибка: "remote origin already exists"
```bash
git remote remove origin
git remote add origin https://github.com/YOUR_USERNAME/saas-platform.git
```

### Ошибка: "failed to push some refs"
```bash
git pull origin main --allow-unrelated-histories
git push -u origin main
```

### Ошибка аутентификации
Используйте Personal Access Token вместо пароля:
1. GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate new token
3. Выберите scope: `repo`
4. Используйте токен как пароль при push

## Следующие шаги

После загрузки проекта:
1. Добавьте описание репозитория
2. Настройте темы (topics): `go`, `react`, `postgresql`, `docker`, `saas`
3. Добавьте лицензию (если нужно)
4. Настройте GitHub Actions для CI/CD (опционально)
5. Добавьте badges в README (опционально)

