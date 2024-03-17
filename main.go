package main

import (
	"encoding/json"
	"flag"
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

func sendJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	http.HandleFunc("/all_tasks", getAllTasksHandler)
	http.HandleFunc("/update_task", updateTaskHandler)
	http.HandleFunc("/delete_task", deleteTaskHandler)
	http.HandleFunc("/add_task", addTaskHandler)

	port := flag.Int("port", 8080, "HTTP server port")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)

	fmt.Printf("Server is running on http://localhost:%d\n", *port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	if !methodNotAllowedHandler(w, r, http.MethodGet) {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var all_tasks []Task
	for _, task := range tasks {
		all_tasks = append(all_tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all_tasks)
}

/*
Updates the task based on its TASK ID.
Can update the task description and use it to update the Completed status (bool)
*/
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !methodNotAllowedHandler(w, r, http.MethodPost) {
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

	// Write tasks to JSON file
	if err := writeTasksToJSONFile(tasks, "tasks.json"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Updated Task successfully")
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addTaskHandler")
	if !methodNotAllowedHandler(w, r, http.MethodPost) {
		return
	}

	var req Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "Task added successfully")
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !methodNotAllowedHandler(w, r, http.MethodDelete) {
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

	// Write tasks to JSON file
	if err := writeTasksToJSONFile(tasks, "tasks.json"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "Task deleted successfully")
}
