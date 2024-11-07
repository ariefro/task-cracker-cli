package main

import (
	"fmt"
	"os"
)

var path = "./notes.json"

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func createFile() {
	// detect if file already exists
	_, err := os.Stat(path)

	// create a new file if it doesn't exist
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if isError(err) {
			return
		}
		defer file.Close()
	}
}

func main() {
	createFile()
}
