# SpeakApper - AI Note Taking App

Современное веб-приложение для создания заметок с помощью ИИ, написанное на Vue.js и Go.

## 🚀 Быстрый старт

### 1. Настройка Google OAuth

1. Следуйте инструкциям в файле `GOOGLE_OAUTH_SETUP.md`
2. Скопируйте `env.example` в `.env` и заполните `VITE_GOOGLE_CLIENT_ID`

### 2. Запуск фронтенда (Vue.js)

```bash
# Установка зависимостей
npm install

# Запуск в режиме разработки
npm run dev
```

Фронтенд будет доступен по адресу: http://localhost:3001

### 3. Запуск бэкенда (Go)

```bash
# Переход в папку бэкенда
cd backend

# Установка зависимостей
go mod tidy

# Запуск сервера
go run .
```

Бэкенд будет доступен по адресу: http://localhost:8080

## 📁 Структура проекта

```
turboaii/
├── src/                    # Фронтенд (Vue.js)
│   ├── App.vue            # Корневой компонент
│   ├── Home.vue           # Главная страница
│   ├── Signup.vue         # Страница регистрации
│   ├── main.js            # Точка входа
│   ├── router.js          # Роутинг
│   ├── style.css          # Стили
│   └── animations.js      # Анимации
├── backend/               # Бэкенд (Go)
│   ├── main.go            # Основной сервер
│   ├── README.md          # Документация бэкенда
│   └── test_api.http      # Тесты API
├── package.json           # Зависимости фронтенда
└── README.md             # Этот файл
```

## 🎨 Особенности фронтенда

- **Современный дизайн** с темной темой и градиентами
- **Анимированные звезды** на фоне
- **Адаптивная верстка** для всех устройств
- **Плавные переходы** между страницами
- **Интерактивные элементы** с hover эффектами
- **Google OAuth интеграция** для быстрой регистрации

## 🔧 API Endpoints

### Регистрация
```
POST /api/signup
{
  "firstName": "Иван",
  "lastName": "Иванов",
  "email": "ivan@example.com",
  "password": "password123"
}
```

### Вход
```
POST /api/login
{
  "email": "ivan@example.com",
  "password": "password123"
}
```

### Регистрация через Google
```
POST /api/google-signup
{
  "token": "google-access-token",
  "idToken": "google-id-token"
}
```

### Проверка здоровья
```
GET /api/health
```

### Получение всех пользователей
```
GET /api/users
```

## 🛠 Технологии

### Фронтенд
- **Vue.js 3** - фреймворк
- **Vue Router** - роутинг
- **Vite** - сборщик
- **CSS3** - стили и анимации

### Бэкенд
- **Go** - основной язык
- **MongoDB Atlas** - облачная база данных
- **Gorilla Mux** - роутинг
- **JWT** - аутентификация
- **bcrypt** - хеширование паролей

## 🔒 Безопасность

- ✅ Хеширование паролей (bcrypt)
- ✅ JWT токены с истечением
- ✅ Валидация данных
- ✅ CORS настройки
- ✅ Защита от CSRF
- ✅ Google OAuth валидация токенов
- ✅ Автоматическая проверка Google API

## 📱 Функциональность

- ✅ Регистрация пользователей
- ✅ Вход в систему
- ✅ Регистрация через Google
- ✅ JWT аутентификация
- ✅ Сохранение токенов в localStorage
- ✅ Красивый UI/UX

## 🚀 Развертывание

### Фронтенд
```bash
npm run build
```

### Бэкенд
```bash
cd backend
go build -o speakapper-backend
./speapper-backend
```

## 📝 Лицензия

MIT License 