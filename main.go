package main

import (
	"fmt"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/pkg/plugin"
)

func main() {
	config.Setup("./config.lua")

	// d, err := linking.ReadNotesDir("/home/patrick/notes/")
	// // d, err := linking.ReadDocument("/home/patrick/notes/todo.md")
	// if err != nil {
	// 	return
	// }
	// _ = d

	x := config.Config()
	x.NGramRange = append(x.NGramRange, 1)
	fmt.Println(config.Config().NGramRange)

	plugin.LuaPlugin()
}
