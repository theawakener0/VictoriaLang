package repl

import (
	"bufio"
	"fmt"
	"io"

	"victoria/errors"
	"victoria/evaluator"
	"victoria/lexer"
	"victoria/object"
	"victoria/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	evaluator.RegisterBuiltinModules()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		evaluator.SetEvalContext(line, "<repl>")

		l := lexer.New(line)
		p := parser.New(l)
		p.SetSource(line, "<repl>")

		program := p.ParseProgram()
		if p.HasErrors() {
			richErrors := p.RichErrors()
			if len(richErrors) > 0 {
				for _, err := range richErrors {
					_, _ = io.WriteString(out, err.Format())
					_, _ = io.WriteString(out, "\n")
				}
			} else {
				printParserErrors(out, p.Errors())
			}
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			if evaluated.Type() == object.ERROR_OBJ {
				errObj := evaluated.(*object.Error)
				_, _ = io.WriteString(out, evaluator.FormatRichError(errObj))
				_, _ = io.WriteString(out, "\n")
			} else if evaluated.Type() != object.NULL_OBJ {
				_, _ = io.WriteString(out, evaluated.Inspect())
				_, _ = io.WriteString(out, "\n")
			}
		}
	}
}

func printParserErrors(out io.Writer, errs []string) {
	_, _ = io.WriteString(out, fmt.Sprintf("%s%serror%s: parser errors found\n", errors.Bold, errors.BrightRed, errors.Reset))
	for _, msg := range errs {
		_, _ = io.WriteString(out, fmt.Sprintf("  %s|%s %s\n", errors.Cyan, errors.Reset, msg))
	}
}
