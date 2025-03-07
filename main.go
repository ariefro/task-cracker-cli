package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

func listTasks(status string) error {
	tasks, err := loadTasks()
	if isError(err) {
		return err
	}

	filteredTasks := tasks
	if status != "" {
		if status != "done" && status != "in-progress" && status != "todo" {
			fmt.Println("Invalid parameter. Available options: todo, in-progress, done")
			return nil
		}

		filteredTasks = []Task{}
		for _, task := range tasks {
			if task.Status == status {
				filteredTasks = append(filteredTasks, task)
			}
		}
	}

	if len(filteredTasks) == 0 {
		fmt.Printf("No tasks found with status '%s'.\n", status)
		return nil
	}

	for _, task := range filteredTasks {
		fmt.Printf("ID: %d\n", task.Id)
		fmt.Printf("Description: %s\n", task.Description)
		fmt.Printf("Status: %s\n", task.Status)
		fmt.Printf("Created At: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Update At: %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("-------------------------------------------------")
	}

	return nil
}

func updateTaskById(id int, description string) error {
	tasks, err := loadTasks()
	if isError(err) {
		return err
	}

	for i, task := range tasks {
		if task.Id == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()

			err := saveTasks(tasks)
			if isError(err) {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("Task with ID %d not found", id)
}

func markDone(id int) error {
	tasks, err := loadTasks()
	if isError(err) {
		return err
	}

	for i, task := range tasks {
		if task.Id == id {
			tasks[i].Status = "done"
			tasks[i].UpdatedAt = time.Now()

			err := saveTasks(tasks)
			if isError(err) {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("Task with ID %d not found", id)
}

func markInProgress(id int) error {
	tasks, err := loadTasks()
	if isError(err) {
		return err
	}

	for i, task := range tasks {
		if task.Id == id {
			tasks[i].Status = "in-progress"
			tasks[i].UpdatedAt = time.Now()

			err := saveTasks(tasks)
			if isError(err) {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("Task with ID %d not found", id)
}

func deleteTaskById(id int) error {
	tasks, err := loadTasks()
	if isError(err) {
		return err
	}

	for i, task := range tasks {
		if task.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return saveTasks(tasks)
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-tracker <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println(" add		<description>		Add a new task")
		fmt.Println(" list [status]				List all tasks, optionally filtered by status")
		fmt.Println(" update		<id> <description>	Update a task description by ID")
		fmt.Println(" mark-done	<id>			Mark a task status as done")
		fmt.Println(" mark-in-progress	<id>			Mark a task status as in progress")
		fmt.Println(" delete		<id>			Delete a task by ID")
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

	case "list":
		status := ""
		if len(os.Args) > 2 {
			status = os.Args[2]
		}

		err := listTasks(status)
		if isError(err) {
			fmt.Printf("Error listing tasks: %v\n", err)
			return
		}

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: task-tracker update <id> <description>")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if isError(err) {
			fmt.Printf("Invalid ID: %v\n", err)
			return
		}

		description := os.Args[3]
		err = updateTaskById(id, description)
		if isError(err) {
			fmt.Printf("Error updating task: %v\n", err)
		} else {
			fmt.Println("Task updated successfully!")
		}

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-tracker mark-done <id>")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if isError(err) {
			fmt.Errorf("Invalid ID: %v\n", err)
			return
		}

		err = markDone(id)
		if isError(err) {
			fmt.Printf("Error updating status task: %v\n", err)
		} else {
			fmt.Println("Task updated successfully!")
		}

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-tracker mark-in-progress <id>")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if isError(err) {
			fmt.Errorf("Invalid ID: %v\n", err)
			return
		}

		err = markInProgress(id)
		if isError(err) {
			fmt.Printf("Error updating status task: %v\n", err)
		} else {
			fmt.Println("Task updated successfully!")
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-tracker delete <id>")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if isError(err) {
			fmt.Printf("Invalid ID: %v\n", err)
			return
		}

		err = deleteTaskById(id)
		if isError(err) {
			fmt.Printf("Error deleting task: %v\n", err)
		} else {
			fmt.Println("Task deleted successfully!")
		}

	default:
		fmt.Println("Uknown command:", command)
		fmt.Println("Available commands: add, list, update, mark-done, mark-in-progress, delete")
	}
}
