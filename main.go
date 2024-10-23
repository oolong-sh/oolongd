package main

import (
	"fmt"

	"github.com/oolong-sh/oolong/internal/linking"
)

func main() {
	fmt.Println("todo")

	d, err := linking.ReadNotesDir("/home/patrick/notes/")
	// d, err := linking.ReadDocument("/home/patrick/notes/todo.md")
	if err != nil {
		return
	}
	_ = d
}
