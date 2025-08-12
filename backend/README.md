# SpeakApper Backend

Бэкенд для приложения SpeakApper, написанный на Go.

## 🚀 Запуск

### 1. MongoDB Atlas

Проект использует MongoDB Atlas - облачную базу данных MongoDB.
Подключение настроено автоматически.

### 2. Запуск Go сервера

```bash
# Установка зависимостей
go mod tidy

# Запуск сервера
go run main.go
```

Сервер запустится на порту 8080.

## 📡 API Endpoints

### Регистрация
```
POST /api/signup
Content-Type: application/json

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
Content-Type: application/json

{
  "email": "ivan@example.com",
  "password": "password123"
}
```

### Регистрация через Google
```
POST /api/google-signup
Content-Type: application/json

{
  "token": "google-jwt-token"
}
```

### Проверка здоровья сервера
```
GET /api/health
```

## 🔧 Технологии

- **Go** - основной язык
- **MongoDB** - база данных
- **Gorilla Mux** - роутинг
- **Gorilla Handlers** - CORS middleware
- **bcrypt** - хеширование паролей
- **JWT** - аутентификация
- **Docker** - контейнеризация

## 📝 Особенности

- ✅ CORS настроен для фронтенда
- ✅ Валидация данных
- ✅ Хеширование паролей
- ✅ JWT токены
- ✅ MongoDB для хранения данных
- ✅ Docker контейнеризация
- ✅ Mongo Express для управления БД
- ✅ Обработка ошибок
- ✅ Индексы для оптимизации запросов

## 🔒 Безопасность

- Пароли хешируются с помощью bcrypt
- JWT токены с истечением через 7 дней
- Валидация всех входных данных
- CORS настроен только для разрешенных доменов 