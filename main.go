package main

import (
	"os"
)

func main() {
	crib := NewCrib(os.Args[1:])
	repl := NewREPL(crib)

	for repl.update() {
		repl.render()
	}
}

