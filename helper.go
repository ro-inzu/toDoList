package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

/*
Loads json file containing all the entire task list (Completed(true) * Completed(false))
*/
func loadTasksFromJSONFile(filename string) {
	fmt.Println("loadTasksFromJSONFile")
	// Open the JSON file for reading
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode the JSON data from the file into the tasks map
	tasks = make(map[int]Task) // Initialize tasks to an empty map

	// If the file is empty, set nextID to 1 and return
	fileStat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fileStat.Size() == 0 {
		nextID = 1
		return
	}

	// Decode the JSON data from the file into the tasks map
	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		log.Fatal(err)
	}

	// Determine the next available ID
	nextID = 1
	for id := range tasks {
		if id >= nextID {
			nextID = id + 1
		}
	}
}

/*
writeTasksToJSONFile writes tasks to a JSON file.
*/
func writeTasksToJSONFile(tasks map[int]Task, filename string) error {
	fmt.Println("writeTasksToJSONFile")

	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the tasks into JSON and write them to the file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(tasks); err != nil {
		return err
	}

	return nil
}
