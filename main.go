package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

// Response represents a JSON response with a message and status.
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Handler for each route
func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request, allowedMethod string) bool {
	if r.Method != allowedMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func sendErrorResponse(w http.ResponseWriter, status int, errorMessage string) {
	http.Error(w, errorMessage, status)
}

func sendJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

func sendSuccessResponse(w http.ResponseWriter, status int, message string) {
	response := Response{
		Status:  "ok",
		Message: message,
	}
	sendJsonResponse(w, status, response)
}

func main() {
	loadTasksFromJSONFile("tasks.json")
	fmt.Println(reflect.TypeOf(tasks))
	http.HandleFunc("/get_all_tasks", getAllTasksHandler)
	http.HandleFunc("/mark_task_complete", markTaskCompleteHandler)
	http.HandleFunc("/update_task", updateTaskHandler)
	http.HandleFunc("/delete_task", deleteTaskHandler)
	http.HandleFunc("/add_task", addTaskHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var tasksSlice []Task
	for _, task := range tasks {
		tasksSlice = append(tasksSlice, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasksSlice)
}

func markTaskCompleteHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	task, exists := tasks[req.ID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	task.Completed = true
	tasks[req.ID] = task

	w.WriteHeader(http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, exists := tasks[req.ID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	tasks[req.ID] = req

	w.WriteHeader(http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, exists := tasks[req.ID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	delete(tasks, req.ID)

	w.WriteHeader(http.StatusOK)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addTaskHandler")
	if !methodNotAllowedHandler(w, r, http.MethodPost) {
		return
	}

	var req Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	mu.Lock()
	defer mu.Unlock()

	log.Printf("JSON body: %v\n", req)
	req.ID = nextID
	tasks[nextID] = req
	nextID++

	// Write tasks to JSON file
	if err := writeTasksToJSONFile(tasks, "tasks.json"); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "Task added successfully")
}
