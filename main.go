package main

import (
	"fmt"
	"go-compiler/main/tools"
	"os"

	"go-compiler/main/lox"
)

func main() {
	args := os.Args[1:]

	lox := lox.NewLox()
	if len(args) > 1 {
		fmt.Println("%s,%s", len(args), args)
		if (args[0] == "-g") || (args[0] == "--generate") && len(args) == 2 {
			tools.Generate(args[1])
			return
		}
		fmt.Println("Usage: golox [script]")
		fmt.Println("-g: golox -g|--generate [output directory]: Generates AST files")
		os.Exit(64)
		return
	} else if len(args) == 1 {
		lox.RunFile(args[0])
	} else {
		lox.RunPrompt()
	}

}
