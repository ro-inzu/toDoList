package main

import "sync"

// Task represents a task with an ID, Description, and Completion status.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var (
	// Initialize the tasks map with a test task.
	tasks  = map[int]Task{}
	nextID = 2 // Start nextID at 2 since we already have a task with ID 1.
	mu     sync.Mutex
)
