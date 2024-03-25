package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Menu:")
		fmt.Println("1. View all tasks")
		fmt.Println("2. Update task")
		fmt.Println("3. Delete task")
		fmt.Println("4. Add task")
		fmt.Println("5. Exit")

		fmt.Print("Enter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = choice[:len(choice)-1] // remove newline character

		switch choice {
		case "1":
			getAllTasks()
		case "2":
			updateTask(reader)
		case "3":
			deleteTask(reader)
		case "4":
			addTask(reader)
		case "5":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 5.")
		}
	}
}

func getAllTasks() {
	response, err := http.Get("http://localhost:8080/all_tasks")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	var tasks []Task
	if err := json.NewDecoder(response.Body).Decode(&tasks); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println("Tasks:")
	for _, task := range tasks {
		fmt.Printf("\t----------------------\n\tTask ID: %d\n\tCompleted: %t\n\tDescription: %s\n \t----------------------\n", task.ID, task.Completed, task.Description)

	}
}

func updateTask(reader *bufio.Reader) {
	fmt.Print("Enter task ID to update: ")
	idStr, _ := reader.ReadString('\n')
	id, _ := strconv.Atoi(idStr[:len(idStr)-1]) // remove newline character

	fmt.Print("Enter new task description: ")
	description, _ := reader.ReadString('\n')
	description = description[:len(description)-1] // remove newline character

	fmt.Print("Enter new completion status (true/false): ")
	completedStr, _ := reader.ReadString('\n')
	completedStr = completedStr[:len(completedStr)-1] // remove newline character
	completed, _ := strconv.ParseBool(completedStr)

	// Create request body
	body, _ := json.Marshal(map[string]interface{}{
		"id":          id,
		"description": description,
		"completed":   completed,
	})

	response, err := http.Post("http://localhost:8080/update_task", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Failed to update task. Status:", response.Status)
		return
	}

	fmt.Println("Task updated successfully.")
}

func deleteTask(reader *bufio.Reader) {
	fmt.Print("Enter task ID to delete: ")
	idStr, _ := reader.ReadString('\n')
	id, _ := strconv.Atoi(idStr[:len(idStr)-1]) // remove newline character

	// Create request body
	body, _ := json.Marshal(map[string]interface{}{
		"id": id,
	})

	req, _ := http.NewRequest("DELETE", "http://localhost:8080/delete_task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Failed to delete task. Status:", response.Status)
		return
	}

	fmt.Println("Task deleted successfully.")
}

func addTask(reader *bufio.Reader) {
	fmt.Print("Enter task content: ")
	task, _ := reader.ReadString('\n')
	task = task[:len(task)-1] // remove newline character

	// Create request body
	body, _ := json.Marshal(map[string]interface{}{
		"description": task,
		"completed":   false,
	})

	response, err := http.Post("http://localhost:8080/add_task", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		fmt.Println("Failed to add task. Status:", response.Status)
		return
	}

	fmt.Println("Task added successfully.")
}

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
