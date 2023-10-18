package tools

import (
	"fmt"
	"os"
	"strings"
)

func GenerateAst(outputDir string, expr string, types []string) {
	path := outputDir + "/" + expr + ".go"
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file")
		os.Exit(74)
	}
	defer file.Close()
	file.Write([]byte("package Lox\n\n"))
	file.Write([]byte(fmt.Sprintf("type %v struct {\n}", expr)))

	for _, t := range types {
		splitType := strings.Split(t, ":")
		className := strings.Trim(splitType[0], " ")
		fields := strings.Trim(splitType[1], " ")

		defineType(file, expr, className, fields)
	}
}

func defineType(file *os.File, baseName string, className string, fieldList string) {
	file.Write([]byte(fmt.Sprintf("\ntype %v struct {\n", className)))
	fields := strings.Split(fieldList, ", ")
	for _, field := range fields {
		file.Write([]byte(fmt.Sprintf("\t%v\n", field)))
	}
	file.Write([]byte("}\n"))
}

func Generate(outputDir string) {
	GenerateAst(outputDir, "Expr", []string{
		"Binary   : left Expr, operator Token , right Expr ",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Unary    : operator Token , right Expr",
	})
}
