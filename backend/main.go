package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// User представляет структуру пользователя
type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Не отправляем пароль в JSON
	CreatedAt time.Time `json:"createdAt"`
}

// SignupRequest представляет запрос на регистрацию
type SignupRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// SignupResponse представляет ответ на регистрацию
type SignupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
	Token   string `json:"token,omitempty"`
}

// LoginRequest представляет запрос на вход
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GoogleSignupRequest представляет запрос на регистрацию через Google
type GoogleSignupRequest struct {
	Token   string `json:"token"`
	IDToken string `json:"idToken"`
}

// JWT секретный ключ (в продакшене используйте переменную окружения)
var jwtSecret = []byte("your-secret-key")

func main() {
	// Подключаемся к MongoDB
	if err := ConnectDB(); err != nil {
		log.Fatal("❌ Ошибка подключения к MongoDB:", err)
	}
	defer DisconnectDB()
	r := mux.NewRouter()

	// Настройка CORS
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3001", "http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	// Роуты
	r.HandleFunc("/api/signup", signupHandler).Methods("POST")
	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/google-signup", googleSignupHandler).Methods("POST")
	r.HandleFunc("/api/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/users", getAllUsersHandler).Methods("GET")

	// Применяем CORS middleware
	handler := corsMiddleware(r)

	fmt.Println("🚀 SpeakApper Backend запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// Обработчик регистрации
func signupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success": false, "message": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	// Валидация
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		http.Error(w, `{"success": false, "message": "Все поля обязательны для заполнения"}`, http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, `{"success": false, "message": "Пароль должен содержать минимум 6 символов"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли пользователь
	exists, err := UserExists(req.Email)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Ошибка при проверке пользователя"}`, http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, `{"success": false, "message": "Пользователь с таким email уже существует"}`, http.StatusConflict)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Ошибка при создании пользователя"}`, http.StatusInternalServerError)
		return
	}

	// Создаем пользователя
	user := &User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
	}

	// Сохраняем пользователя в MongoDB
	if err := CreateUser(user); err != nil {
		http.Error(w, `{"success": false, "message": "Ошибка при сохранении пользователя"}`, http.StatusInternalServerError)
		return
	}

	// Генерируем JWT токен
	token, err := generateJWT(user)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Ошибка при создании токена"}`, http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	response := SignupResponse{
		Success: true,
		Message: "Пользователь успешно зарегистрирован",
		User:    user,
		Token:   token,
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик входа
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success": false, "message": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	// Получаем пользователя из базы данных
	user, err := GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Неверный email или пароль"}`, http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, `{"success": false, "message": "Неверный email или пароль"}`, http.StatusUnauthorized)
		return
	}

	// Генерируем JWT токен
	token, err := generateJWT(user)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Ошибка при создании токена"}`, http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	response := SignupResponse{
		Success: true,
		Message: "Успешный вход",
		User:    user,
		Token:   token,
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик регистрации через Google
func googleSignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req GoogleSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success": false, "message": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Token == "" {
		http.Error(w, `{"success": false, "message": "Токен Google обязателен"}`, http.StatusBadRequest)
		return
	}

	// Валидируем Google токен и получаем информацию о пользователе
	googleUser, err := ValidateGoogleToken(req.Token)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Неверный Google токен"}`, http.StatusUnauthorized)
		return
	}

	// Проверяем, существует ли пользователь с таким email
	exists, err := UserExists(googleUser.Email)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Ошибка при проверке пользователя"}`, http.StatusInternalServerError)
		return
	}

	var user *User
	if exists {
		// Пользователь уже существует, получаем его данные
		user, err = GetUserByEmail(googleUser.Email)
		if err != nil {
			http.Error(w, `{"success": false, "message": "Ошибка при получении пользователя"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// Создаем нового пользователя
		user = &User{
			FirstName: googleUser.GivenName,
			LastName:  googleUser.FamilyName,
			Email:     googleUser.Email,
			Password:  "", // Google пользователи не имеют пароля
		}

		// Сохраняем пользователя в MongoDB
		if err := CreateUser(user); err != nil {
			http.Error(w, `{"success": false, "message": "Ошибка при сохранении пользователя"}`, http.StatusInternalServerError)
			return
		}
	}

	// Генерируем JWT токен
	token, err := generateJWT(user)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Ошибка при создании токена"}`, http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	response := SignupResponse{
		Success: true,
		Message: "Успешная авторизация через Google",
		User:    user,
		Token:   token,
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик проверки здоровья сервера
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "SpeakApper Backend работает",
		"time":    time.Now(),
	})
}

// Обработчик получения всех пользователей
func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Ошибка при получении пользователей",
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"users":   users,
		"count":   len(users),
	})
}

// Генерация JWT токена
func generateJWT(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
	})

	return token.SignedString(jwtSecret)
}
