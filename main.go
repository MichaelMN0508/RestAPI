package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var task string

type requestBody struct {
	Task string `json:"task"`
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
func main() {
	InitDB()
	DB.AutoMigrate(&Task{})
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", GetTask).Methods("GET")
	router.HandleFunc("/api/tasks", CreateTask).Methods("POST")
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", router)
}
