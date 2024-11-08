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

	var tasks []Task
	err = json.Unmarshal(file, &tasks)
	return tasks, err
}

func saveNotes(tasks []Task) error {
	dataByte, err := json.MarshalIndent(tasks, "", "  ")
	if isError(err) {
		return err
	}

	err = os.WriteFile(dataFile, dataByte, 0644)

	return err
}

func main() {
	tasks, err := loadTasks()
	if isError(err) {
		fmt.Println(err)
	}

	fmt.Println("====", tasks)
}
