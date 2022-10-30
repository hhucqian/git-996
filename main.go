package main

import (
	"fmt"
	"git-analyse/get_repository"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage : ", os.Args[0], "<repository path>")
		return
	}

	repository_path := os.Args[1]
	fmt.Println("repository path:", repository_path)

	repository := get_repository.New(repository_path)

	repository.Anaylse()
}
