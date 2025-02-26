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

func GETHandler(w http.ResponseWriter, r *http.Request) {
	if task == "" {
		fmt.Fprintln(w, "Hello!")
	} else {
		fmt.Fprintln(w, "Hello!", task)
	}
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody requestBody
	json.NewDecoder(r.Body).Decode(&requestBody)

	task = requestBody.Task
	fmt.Fprintf(w, "Task updated:%s", task)
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api", GETHandler).Methods("GET")
	router.HandleFunc("/api", POSTHandler).Methods("POST")
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", router)
}
