package lox

import (
	"bufio"
	"fmt"
	"go-compiler/main/errors"
	"go-compiler/main/evaluator"
	"go-compiler/main/object"
	"go-compiler/main/parser"
	"go-compiler/main/scanner"
	"os"
	"os/user"
)

type Lox struct {
}

func NewLox() *Lox {
	return &Lox{}
}

func (l *Lox) RunFile(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file")
		os.Exit(74)
	}
	env := object.NewEnvironment(nil)
	l.Run(string(b), env)
	if errors.HadError {
		os.Exit(65)
	}
}

func (l *Lox) RunPrompt() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Welcome to %s! Let's get down to monkey business!\n", u.Username)
	fmt.Println("Explore mokey by writting some code:")
	env := object.NewEnvironment(nil)
	for {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		// Potetially use ReadBytes('\n') instead
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("Error reading input")
			os.Exit(74)
		}
		if line == nil {
			break
		}
		if string(line) == "exit" {
			return
		}
		l.Run(string(line), env)
		errors.HadError = false
	}
}

func (l *Lox) Run(source string, env *object.Environment) {
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()
	p := parser.NewParser(tokens)
	program := p.Parse()
	if len(p.Errors()) != 0 {
		return
	}

	eval := evaluator.Eval(program, env)
	if eval != nil {
		fmt.Printf("%s\n", eval.Inspect())
	}

}
