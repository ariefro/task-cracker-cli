package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var dataFile = "tasks.json"

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func createFile(err error) {
	// create a new file if it doesn't exist
	if os.IsNotExist(err) {
		file, err := os.Create(dataFile)
		if isError(err) {
			return
		}
		defer file.Close()
	}
}

func loadTasks() ([]Task, error) {
	file, err := os.ReadFile(dataFile)
	if isError(err) {
		createFile(err)
	}

	if len(file) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	err = json.Unmarshal(file, &tasks)

	return tasks, err
}

func saveTasks(tasks []Task) error {
	dataByte, err := json.MarshalIndent(tasks, "", "  ")
	if isError(err) {
		return err
	}

	err = os.WriteFile(dataFile, dataByte, 0644)

	return err
}

func addTask(description string) error {
	tasks, err := loadTasks()
	if isError(err) {
		return err
	}

	id := 1
	if len(tasks) > 0 {
		id = tasks[len(tasks)-1].Id + 1
	}

	task := &Task{
		Id:          id,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks = append(tasks, *task)

	return saveTasks(tasks)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-tracker <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println(" add <description> Add a new task")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-tracker add <description>")
			return
		}

		description := os.Args[2]
		err := addTask(description)
		if isError(err) {
			fmt.Printf("Error adding task: %v\n", err)
		} else {
			fmt.Println("Task added successfully!")
		}

	default:
		fmt.Println("Uknown command:", command)
		fmt.Println("Available commands: add, list, view, update, delete")
	}
}
