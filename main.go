package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Task represents a task with an ID, Description, and Completion status.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var (
	// Initialize the tasks map with a test task.
	tasks = map[int]Task{
		1: {ID: 1, Description: "Test task", Completed: false},
	}
	nextID = 2 // Start nextID at 2 since we already have a task with ID 1.
	mu     sync.Mutex
)

func main() {
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

	req.ID = nextID
	tasks[nextID] = req
	nextID++

	w.WriteHeader(http.StatusCreated)
}
