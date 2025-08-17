package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// OpenAI API ключ теперь берём из переменной окружения
var openaiAPIKey string

// OpenAI API структуры
type WhisperRequest struct {
	Model    string `json:"model"`
	File     string `json:"file"`
	Language string `json:"language,omitempty"`
}

type WhisperResponse struct {
	Text string `json:"text"`
}

// Note структура для MongoDB
type Note struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title      string             `bson:"title" json:"title"`
	Content    string             `bson:"content" json:"content"`
	Type       string             `bson:"type" json:"type"`
	Tab        string             `bson:"tab" json:"tab"`
	LastOpened string             `bson:"last_opened" json:"last_opened"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// GPT generation types
type GenerateRequest struct {
	Transcript string `json:"transcript"`
	Language   string `json:"language,omitempty"`
}

type Flashcard struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
	Example    string `json:"example,omitempty"`
}

type QuizQuestion struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}

type GeneratePayload struct {
	Flashcards   []Flashcard    `json:"flashcards"`
	Quiz         []QuizQuestion `json:"quiz"`
	LanguageCode string         `json:"languageCode,omitempty"`
}

// Учебные материалы (материализованные карточки/квиз)
type Material struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Transcript string             `bson:"transcript" json:"transcript"`
	Flashcards []Flashcard        `bson:"flashcards" json:"flashcards"`
	Quiz       []QuizQuestion     `bson:"quiz" json:"quiz"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

func handleMaterials(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Handle GET request - fetch user materials
	if r.Method == "GET" {
		collection := client.Database("speakapper").Collection("materials")

		// Find all materials for this user
		cursor, err := collection.Find(context.Background(), bson.M{"user_id": userID})
		if err != nil {
			log.Printf("Error fetching materials: %v", err)
			http.Error(w, "Failed to fetch materials", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		var materials []Material
		if err = cursor.All(context.Background(), &materials); err != nil {
			log.Printf("Error decoding materials: %v", err)
			http.Error(w, "Failed to decode materials", http.StatusInternalServerError)
			return
		}

		// Convert ObjectIDs to strings for JSON response
		var responseMaterials []map[string]interface{}
		for _, mat := range materials {
			responseMaterials = append(responseMaterials, map[string]interface{}{
				"id":         mat.ID.Hex(),
				"transcript": mat.Transcript,
				"flashcards": mat.Flashcards,
				"quiz":       mat.Quiz,
				"created_at": mat.CreatedAt,
				"updated_at": mat.UpdatedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"materials": responseMaterials,
		})
		return
	}

	// Handle POST request - create new material
	var materialData struct {
		Transcript string         `json:"transcript"`
		Flashcards []Flashcard    `json:"flashcards"`
		Quiz       []QuizQuestion `json:"quiz"`
	}

	if err := json.NewDecoder(r.Body).Decode(&materialData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create material
	material := Material{
		UserID:     userID,
		Transcript: materialData.Transcript,
		Flashcards: materialData.Flashcards,
		Quiz:       materialData.Quiz,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to MongoDB
	collection := client.Database("speakapper").Collection("materials")
	result, err := collection.InsertOne(context.Background(), material)
	if err != nil {
		log.Printf("Error saving material: %v", err)
		http.Error(w, "Failed to save material", http.StatusInternalServerError)
		return
	}

	// Set the ID from the inserted document
	material.ID = result.InsertedID.(primitive.ObjectID)

	// Return the created material
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"material": material,
	})
}

func handleNotes(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Handle GET request - fetch user notes
	if r.Method == "GET" {
		collection := client.Database("speakapper").Collection("notes")

		// Find all notes for this user
		cursor, err := collection.Find(context.Background(), bson.M{"user_id": userID})
		if err != nil {
			log.Printf("Error fetching notes: %v", err)
			http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		var notes []Note
		if err = cursor.All(context.Background(), &notes); err != nil {
			log.Printf("Error decoding notes: %v", err)
			http.Error(w, "Failed to decode notes", http.StatusInternalServerError)
			return
		}

		// Convert ObjectIDs to strings for JSON response
		var responseNotes []map[string]interface{}
		for _, note := range notes {
			responseNotes = append(responseNotes, map[string]interface{}{
				"id":          note.ID.Hex(),
				"title":       note.Title,
				"content":     note.Content,
				"type":        note.Type,
				"tab":         note.Tab,
				"last_opened": note.LastOpened,
				"created_at":  note.CreatedAt,
				"updated_at":  note.UpdatedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"notes":   responseNotes,
		})
		return
	}

	// Handle POST request - create new note
	// Parse request body
	var noteData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Type    string `json:"type"`
		Tab     string `json:"tab"`
	}

	if err := json.NewDecoder(r.Body).Decode(&noteData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create note
	note := Note{
		UserID:     userID,
		Title:      noteData.Title,
		Content:    noteData.Content,
		Type:       noteData.Type,
		Tab:        noteData.Tab,
		LastOpened: "Just now",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to MongoDB
	collection := client.Database("speakapper").Collection("notes")
	result, err := collection.InsertOne(context.Background(), note)
	if err != nil {
		log.Printf("Error saving note: %v", err)
		http.Error(w, "Failed to save note", http.StatusInternalServerError)
		return
	}

	// Set the ID from the inserted document
	note.ID = result.InsertedID.(primitive.ObjectID)

	// Return the created note
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"note":    note,
	})
}

func handleTranscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (allow large)
	if err := r.ParseMultipartForm(1024 << 20); err != nil { // 1GB streamed to temp
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "No audio file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Received audio file: %s, size: %d bytes", header.Filename, header.Size)

	// Save upload to temp file on disk to avoid memory blowups
	tmpDir := os.TempDir()
	tmpIn := filepath.Join(tmpDir, fmt.Sprintf("upload_%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename)))
	out, err := os.Create(tmpIn)
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		os.Remove(tmpIn)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	out.Close()
	defer os.Remove(tmpIn)

	// Threshold for long audio (e.g., > 20MB)
	const longThreshold = 20 * 1024 * 1024
	if header.Size > longThreshold {
		text, err := transcribeLongAudio(tmpIn)
		if err != nil {
			log.Printf("Long transcription error: %v", err)
			http.Error(w, "Transcription failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":       true,
			"transcription": text,
			"filename":      header.Filename,
			"size":          header.Size,
			"mode":          "segmented",
		})
		return
	}

	// Small file: read and send directly
	fileBytes, err := os.ReadFile(tmpIn)
	if err != nil {
		http.Error(w, "Failed to read temp file", http.StatusInternalServerError)
		return
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		http.Error(w, "Failed to create form file", http.StatusInternalServerError)
		return
	}
	part.Write(fileBytes)
	writer.WriteField("model", "whisper-1")
	writer.Close()

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &requestBody)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		http.Error(w, "Failed to transcribe audio", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("OpenAI API error: %s - %s", resp.Status, string(respBody))
		http.Error(w, "Transcription failed", http.StatusInternalServerError)
		return
	}

	var whisperResp WhisperResponse
	if err := json.Unmarshal(respBody, &whisperResp); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"transcription": whisperResp.Text,
		"filename":      header.Filename,
		"size":          header.Size,
		"mode":          "single",
	})
}

func transcribeLongAudio(inputPath string) (string, error) {
	// Check ffmpeg availability
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("ffmpeg not found: %w", err)
	}

	workDir, err := os.MkdirTemp(os.TempDir(), "chunks_*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(workDir)

	// Segment into ~10-minute chunks, mono 16kHz low bitrate to reduce size
	chunkPattern := filepath.Join(workDir, "chunk_%03d.mp3")
	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, "-ac", "1", "-ar", "16000", "-b:a", "64k", "-f", "segment", "-segment_time", "600", chunkPattern)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg segment failed: %w", err)
	}

	// Collect chunks
	entries, err := os.ReadDir(workDir)
	if err != nil {
		return "", err
	}
	var chunkFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "chunk_") {
			chunkFiles = append(chunkFiles, filepath.Join(workDir, e.Name()))
		}
	}
	if len(chunkFiles) == 0 {
		return "", fmt.Errorf("no chunks produced")
	}
	// Sort by name to keep order
	slices.Sort(chunkFiles)

	var fullText strings.Builder
	httpClient := &http.Client{Timeout: 180 * time.Second}
	for _, path := range chunkFiles {
		// Build multipart per chunk
		var body bytes.Buffer
		mw := multipart.NewWriter(&body)
		chunk, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		p, err := mw.CreateFormFile("file", filepath.Base(path))
		if err != nil {
			return "", err
		}
		if _, err := p.Write(chunk); err != nil {
			return "", err
		}
		mw.WriteField("model", "whisper-1")
		mw.Close()

		req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &body)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+openaiAPIKey)
		req.Header.Set("Content-Type", mw.FormDataContentType())

		resp, err := httpClient.Do(req)
		if err != nil {
			return "", err
		}
		respBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("chunk transcribe error: %s - %s", resp.Status, string(respBytes))
		}
		var wr WhisperResponse
		if err := json.Unmarshal(respBytes, &wr); err != nil {
			return "", err
		}
		fullText.WriteString(strings.TrimSpace(wr.Text))
		fullText.WriteString("\n")
		// throttle a bit to be safe
		time.Sleep(500 * time.Millisecond)
	}

	return fullText.String(), nil
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Transcript == "" {
		http.Error(w, "Transcript is required", http.StatusBadRequest)
		return
	}

	// Build chat completion request
	chatReq := map[string]interface{}{
		"model":       "gpt-4o-mini",
		"temperature": 0.3,
		"messages": []map[string]string{
			{"role": "system", "content": "You convert transcripts into study materials. IMPORTANT RULES: 1) USE ONLY words and facts from the transcript; DO NOT invent. 2) Flashcards: term must be an exact word/phrase from transcript; definition should be a short sentence fragment from transcript (or closest sentence). 3) Quiz: each question must be based on transcript; the correct answer and ALL options must be text spans that appear in the transcript. 4) If there is not enough material, return fewer items. Respond strictly as compact JSON with keys: flashcards (array of {term, definition, example?}) and quiz (array of {id?, topicId?, question, options[4], answer, hint?}). No extra commentary, no markdown fences."},
			{"role": "user", "content": fmt.Sprintf("Language hint: %s\nTranscript:\n%s\n\nReturn JSON only.", reqBody.Language, reqBody.Transcript)},
		},
	}

	buf, _ := json.Marshal(chatReq)
	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(buf))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	httpReq.Header.Set("Authorization", "Bearer "+openaiAPIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 45 * time.Second}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Printf("OpenAI chat API error: %v", err)
		http.Error(w, "Failed to generate materials", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("OpenAI chat API error: %s - %s", resp.Status, string(respBytes))
		http.Error(w, "Generation failed", http.StatusInternalServerError)
		return
	}

	// Parse choices[0].message.content
	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBytes, &openaiResp); err != nil {
		http.Error(w, "Failed to parse OpenAI response", http.StatusInternalServerError)
		return
	}
	if len(openaiResp.Choices) == 0 {
		http.Error(w, "Empty OpenAI response", http.StatusInternalServerError)
		return
	}

	var payload GeneratePayload
	if err := json.Unmarshal([]byte(openaiResp.Choices[0].Message.Content), &payload); err != nil {
		// If model returned fenced code block, try strip
		clean := openaiResp.Choices[0].Message.Content
		// naive cleanup of Markdown fences
		clean = strings.TrimPrefix(clean, "```json")
		clean = strings.TrimPrefix(clean, "```")
		clean = strings.TrimSuffix(clean, "```")
		if err2 := json.Unmarshal([]byte(clean), &payload); err2 != nil {
			log.Printf("Failed to unmarshal AI JSON: %v; content: %s", err, openaiResp.Choices[0].Message.Content)
			http.Error(w, "Invalid JSON from model", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"flashcards": payload.Flashcards,
		"quiz":       payload.Quiz,
	})
}

func main() {
	// Читаем OpenAI API ключ из переменной окружения
	openaiAPIKey = os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		log.Fatal("❌ OPENAI_API_KEY не задан в переменных окружения!")
	}

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
	r.HandleFunc("/api/transcribe", handleTranscribe).Methods("POST")
	r.HandleFunc("/api/notes", handleNotes).Methods("POST")
	r.HandleFunc("/api/notes", handleNotes).Methods("GET")
	r.HandleFunc("/api/generate", handleGenerate).Methods("POST")
	r.HandleFunc("/api/materials", handleMaterials).Methods("POST", "GET")

	// Применяем CORS middleware
	handler := corsMiddleware(r)

	fmt.Println("🚀 SpeakApper Backend запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// Обработчик регистрации
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли пользователь
	exists, err := UserExists(req.Email)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(SignupResponse{
			Success: false,
			Message: "User already exists",
		})
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Создаем пользователя
	user := &User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
	}

	// Сохраняем в базу данных
	if err := CreateUser(user); err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Генерируем JWT токен
	token, err := generateJWT(user)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SignupResponse{
		Success: true,
		Message: "User registered successfully",
		User:    user,
		Token:   token,
	})
}

// Обработчик входа
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем пользователя из базы данных
	user, err := GetUserByEmail(req.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid credentials",
		})
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid credentials",
		})
		return
	}

	// Генерируем JWT токен
	token, err := generateJWT(user)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

// Обработчик регистрации через Google
func googleSignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GoogleSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидируем Google токен
	googleUser, err := ValidateGoogleToken(req.Token)
	if err != nil {
		log.Printf("Error validating Google token: %v", err)
		http.Error(w, "Invalid Google token", http.StatusUnauthorized)
		return
	}

	// Проверяем, существует ли пользователь
	exists, err := UserExists(googleUser.Email)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var user *User
	if exists {
		// Пользователь существует, получаем его данные
		user, err = GetUserByEmail(googleUser.Email)
		if err != nil {
			log.Printf("Error getting existing user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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

		if err := CreateUser(user); err != nil {
			log.Printf("Error creating Google user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Генерируем JWT токен
	token, err := generateJWT(user)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Google authentication successful",
		"user":    user,
		"token":   token,
	})
}

// Обработчик проверки здоровья сервера
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"time":   time.Now(),
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
