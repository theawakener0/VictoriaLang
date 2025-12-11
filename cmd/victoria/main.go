package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"victoria/evaluator"
	"victoria/lexer"
	"victoria/object"
	"victoria/parser"
	"victoria/repl"
)

func main() {
	if len(os.Args) > 1 {
		filename := os.Args[1]
		runFile(filename)
	} else {
		fmt.Printf("Victoria Programming Language\n")
		fmt.Printf("Type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}

func runFile(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	env := object.NewEnvironment()
	evaluator.RegisterBuiltinModules()
	l := lexer.New(string(data))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
	}
}
