package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
}

type Task struct {
	gorm.Model
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

var task string
var nextID int = 1

type requestBody struct {
	Task string `json:"task"`
	ID   int    `json:"id"`
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	DB.Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	DB.Create(&task)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /api/tasks/{id} вызван") // Лог для проверки

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("Ошибка ID:", err)
		http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
		return
	}

	fmt.Println("Обновление задачи с ID:", id) // Проверяем, какой ID передаётся

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	var task Task
	if err := DB.First(&task, id).Error; err != nil {
		fmt.Println("Задача не найдена:", err)
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Обновляем поля
	task.Task = updatedTask.Task
	task.IsDone = updatedTask.IsDone

	// Сохраняем изменения в базе
	if err := DB.Save(&task).Error; err != nil {
		fmt.Println("Ошибка обновления в БД:", err)
		http.Error(w, "Ошибка обновления в БД", http.StatusInternalServerError)
		return
	}

	fmt.Println("Задача успешно обновлена:", task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
		return
	}

	var task Task
	if err := DB.First(&task, id).Error; err != nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	DB.Delete(&task)

	w.WriteHeader(http.StatusNoContent)
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Запрос:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func main() {
	InitDB()
	DB.AutoMigrate(&Task{})
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", GetTask).Methods("GET")
	router.HandleFunc("/api/tasks/{id}", UpdateTask).Methods("PUT")
	router.HandleFunc("/api/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/api/tasks/{id}", DeleteTask).Methods("DELETE")
	fmt.Println("Registered routes:")
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		fmt.Println(methods, path)
		return nil
	})

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", router)
}
