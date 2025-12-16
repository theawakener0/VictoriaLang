package main

import (
	"fmt"
	"os"
	"path/filepath"

	"victoria/errors"
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
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s%serror%s: could not read file '%s'\n", errors.Bold, errors.BrightRed, errors.Reset, filename)
		fmt.Printf("  %s|%s %v\n", errors.Cyan, errors.Reset, err)
		return
	}

	source := string(data)
	absPath, _ := filepath.Abs(filename)

	env := object.NewEnvironment()
	evaluator.RegisterBuiltinModules()
	evaluator.SetEvalContext(source, absPath)
	defer evaluator.ClearEvalContext()

	l := lexer.New(source)
	p := parser.New(l)
	p.SetSource(source, absPath)

	program := p.ParseProgram()
	if p.HasErrors() {
		richErrors := p.RichErrors()
		if len(richErrors) > 0 {
			for _, err := range richErrors {
				fmt.Print(err.Format())
				fmt.Println()
			}
			fmt.Printf("\n%s%serror%s: could not compile due to %d previous error(s)\n",
				errors.Bold, errors.BrightRed, errors.Reset, len(richErrors))
		} else {
			// Fallback to simple errors
			for _, msg := range p.Errors() {
				fmt.Println(msg)
			}
		}
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		errObj := evaluated.(*object.Error)
		fmt.Print(evaluator.FormatRichError(errObj))
	}
}
