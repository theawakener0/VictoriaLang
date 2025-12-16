package evaluator

import (
	"fmt"
	"strings"
	"victoria/ast"
	"victoria/errors"
	"victoria/object"
)

// EvalContext holds context for evaluation including source code for error messages
type EvalContext struct {
	SourceCode string
	Filename   string
}

// Global context for error reporting - set by the runner
var currentContext *EvalContext

// SetEvalContext sets the global evaluation context for error reporting
func SetEvalContext(source string, filename string) {
	currentContext = &EvalContext{
		SourceCode: source,
		Filename:   filename,
	}
}

// ClearEvalContext clears the evaluation context
func ClearEvalContext() {
	currentContext = nil
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// newErrorWithLocation creates an error with location information
func newErrorWithLocation(format string, line, col, endCol int, a ...interface{}) *object.Error {
	return &object.Error{
		Message:   fmt.Sprintf(format, a...),
		Line:      line,
		Column:    col,
		EndColumn: endCol,
	}
}

// newTypeMismatchError creates a rich type mismatch error
func newTypeMismatchError(left, operator, right string, tok *ast.InfixExpression) *object.Error {
	line := tok.Token.Line
	col := tok.Token.Column
	endCol := tok.Token.EndColumn
	msg := fmt.Sprintf("type mismatch: %s %s %s", left, operator, right)
	return &object.Error{
		Message:   msg,
		Line:      line,
		Column:    col,
		EndColumn: endCol,
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// FormatRichError formats an object.Error into a rich error display
func FormatRichError(err *object.Error) string {
	if currentContext == nil || err.Line == 0 {
		return err.Inspect()
	}

	loc := errors.SourceLocation{
		Line:      err.Line,
		Column:    err.Column,
		EndColumn: err.EndColumn,
		Filename:  currentContext.Filename,
	}

	richErr := errors.NewRuntimeError(err.Message, loc, currentContext.SourceCode)

	// Add context-specific help and notes based on error message
	msg := err.Message

	// TYPE MISMATCH ERRORS
	if strings.Contains(msg, "type mismatch") {
		richErr.WithCode("E0001")
		richErr.WithNote("Victoria is dynamically typed, but operators require compatible types")

		if strings.Contains(msg, "STRING") && strings.Contains(msg, "INTEGER") {
			richErr.WithNote("strings and integers cannot be combined directly with arithmetic operators")
			richErr.WithHelp("use string() to convert integers to strings: \"text\" + string(42)")
		}
		if strings.Contains(msg, "STRING") && strings.Contains(msg, "FLOAT") {
			richErr.WithNote("strings and floats cannot be combined directly")
			richErr.WithHelp("use string() to convert floats to strings: \"value: \" + string(3.14)")
		}
		if strings.Contains(msg, "BOOLEAN") {
			richErr.WithNote("booleans cannot be used in arithmetic operations")
			if strings.Contains(msg, "INTEGER") {
				richErr.WithHelp("use int() to convert boolean to integer: int(true) returns 1")
			} else {
				richErr.WithHelp("use string() to convert boolean to string: string(true) returns \"true\"")
			}
		}
		if strings.Contains(msg, "ARRAY") {
			richErr.WithNote("arrays can only be concatenated with other arrays using '+'")
			richErr.WithHelp("use push() to add elements: push(arr, element)")
		}
		if strings.Contains(msg, "HASH") {
			richErr.WithNote("hashes do not support arithmetic operators")
			richErr.WithHelp("access hash values with hash[\"key\"] or hash.key syntax")
		}

	} else if strings.Contains(msg, "identifier not found") {
		name := strings.TrimPrefix(msg, "identifier not found: ")
		richErr.WithCode("E0002")
		richErr.WithNote("variables must be declared before use with 'let' or 'const'")

		commonBuiltins := map[string]string{
			"println":   "use 'print' instead - Victoria uses print() for output",
			"printf":    "use 'format' instead - Victoria uses format() for formatted strings",
			"console":   "use 'print' instead - Victoria uses print() for output",
			"log":       "use 'print' instead - Victoria uses print() for output",
			"echo":      "use 'print' instead - Victoria uses print() for output",
			"puts":      "use 'print' instead - Victoria uses print() for output",
			"str":       "use 'string' instead - Victoria uses string() for type conversion",
			"toString":  "use 'string' instead - Victoria uses string() for type conversion",
			"toInt":     "use 'int' instead - Victoria uses int() for type conversion",
			"parseInt":  "use 'int' instead - Victoria uses int() for type conversion",
			"size":      "use 'len' instead - Victoria uses len() for collection length",
			"length":    "use 'len' instead - Victoria uses len() for collection length",
			"count":     "use 'len' instead - Victoria uses len() for collection length",
			"append":    "use 'push' instead - Victoria uses push(array, element)",
			"add":       "use 'push' instead - Victoria uses push(array, element)",
			"remove":    "use 'pop' instead - Victoria uses pop(array) to remove last element",
			"delete":    "use filter() to create a new array without elements",
			"substr":    "use string slicing instead: str[start:end]",
			"substring": "use string slicing instead: str[start:end]",
			"forEach":   "use a for-in loop: for item in array { ... }",
			"foreach":   "use a for-in loop: for item in array { ... }",
			"nil":       "use 'null' instead - Victoria uses null for no value",
			"none":      "use 'null' instead - Victoria uses null for no value",
			"None":      "use 'null' instead - Victoria uses null for no value",
			"undefined": "use 'null' instead - Victoria uses null for no value",
			"fn":        "use 'define' instead - Victoria uses 'define' to create functions",
			"func":      "use 'define' instead - Victoria uses 'define' to create functions",
			"function":  "use 'define' instead - Victoria uses 'define' to create functions",
			"lambda":    "use 'define' instead - Victoria uses 'define' to create functions",
			"def":       "use 'define' instead - Victoria uses 'define' to create functions",
			"var":       "use 'let' instead - Victoria uses 'let' for variable declaration",
			"elif":      "use 'else if' - Victoria uses 'else if' not 'elif'",
			"elsif":     "use 'else if' - Victoria uses 'else if' not 'elsif'",
			"elseif":    "use 'else if' (with space) - Victoria uses 'else if'",
			"and":       "'and' is an operator, not a function - use it between expressions: a and b",
			"or":        "'or' is an operator, not a function - use it between expressions: a or b",
			"not":       "use '!' for negation - Victoria uses !value not not(value)",
			"self":      "Victoria uses the struct instance name directly in methods",
			"this":      "Victoria uses the struct instance name directly in methods",
			"require":   "use 'include' instead - Victoria uses include \"filename\"",
			"import":    "use 'include' instead - Victoria uses include \"filename\"",
		}

		if suggestion, ok := commonBuiltins[name]; ok {
			richErr.WithHelp(suggestion)
		} else if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
			richErr.WithHelp(fmt.Sprintf("'%s' looks like a type name; did you mean to create an instance?", name))
			richErr.WithNote("use struct instantiation: new StructName { field: value }")
		} else {
			richErr.WithHelp(fmt.Sprintf("declare with: let %s = <value>", name))
			richErr.WithNote("for constants, use: const " + name + " = <value>")
		}

	} else if strings.Contains(msg, "index operator not supported") {
		richErr.WithCode("E0006")
		richErr.WithNote("indexing is only supported for arrays, strings, and hashes")

		if strings.Contains(msg, "INTEGER") {
			richErr.WithHelp("integers cannot be indexed; did you mean to use an array?")
		} else if strings.Contains(msg, "NULL") {
			richErr.WithHelp("cannot index null; ensure the value is initialized")
		} else {
			richErr.WithHelp("only arrays, strings, and hashes support [] indexing")
		}

	} else if strings.Contains(msg, "not a function") {
		richErr.WithCode("E0005")

		if strings.Contains(msg, "INTEGER") {
			richErr.WithNote("integers cannot be called as functions")
			richErr.WithHelp("remove the parentheses, or did you mean to use a function?")
		} else if strings.Contains(msg, "ARRAY") {
			richErr.WithNote("arrays cannot be called as functions")
			richErr.WithHelp("use array[index] to access elements, not array(index)")
		} else if strings.Contains(msg, "HASH") {
			richErr.WithNote("hashes cannot be called as functions")
			richErr.WithHelp("use hash[\"key\"] or hash.key to access values")
		} else if strings.Contains(msg, "nil") || strings.Contains(msg, "NULL") {
			richErr.WithNote("attempted to call null as a function")
			richErr.WithHelp("ensure the variable is assigned a function before calling")
		} else {
			richErr.WithNote("only functions and builtin functions can be called")
		}

	} else if strings.Contains(msg, "unknown operator") {
		richErr.WithCode("E0003")

		if strings.Contains(msg, "STRING") && strings.Contains(msg, "-") {
			richErr.WithNote("strings only support the '+' operator for concatenation")
		} else if strings.Contains(msg, "BOOLEAN") {
			richErr.WithNote("booleans only support comparison operators (==, !=)")
			richErr.WithHelp("use 'and', 'or', '!' for boolean logic")
		} else {
			richErr.WithNote("this operator is not supported for the given types")
		}

	} else if strings.Contains(msg, "variable not defined") {
		name := strings.TrimPrefix(msg, "variable not defined: ")
		richErr.WithCode("E0002")
		richErr.WithNote("cannot modify a variable that hasn't been declared")
		richErr.WithHelp(fmt.Sprintf("declare first: let %s = <initial_value>", name))

	} else if strings.Contains(msg, "struct not found") {
		name := strings.TrimPrefix(msg, "struct not found: ")
		richErr.WithCode("E0009")
		richErr.WithNote("structs must be defined before instantiation")
		richErr.WithHelp(fmt.Sprintf("define the struct first: struct %s { field1, field2 }", name))

	} else if strings.Contains(msg, "wrong number of arguments") {
		richErr.WithCode("E0010")
		richErr.WithNote("function called with incorrect number of arguments")

		if strings.Contains(msg, "len") {
			richErr.WithHelp("len() takes exactly 1 argument: len(array) or len(string)")
		} else if strings.Contains(msg, "push") {
			richErr.WithHelp("push() takes exactly 2 arguments: push(array, element)")
		} else if strings.Contains(msg, "map") {
			richErr.WithHelp("map() takes exactly 2 arguments: map(array, fn)")
		} else if strings.Contains(msg, "filter") {
			richErr.WithHelp("filter() takes exactly 2 arguments: filter(array, predicate)")
		} else if strings.Contains(msg, "reduce") {
			richErr.WithHelp("reduce() takes 2 or 3 arguments: reduce(array, fn) or reduce(array, fn, initial)")
		}

	} else if strings.Contains(msg, "spread operator") {
		richErr.WithCode("E0012")
		if strings.Contains(msg, "requires an array") {
			richErr.WithNote("the spread operator (...) can only expand arrays")
			richErr.WithHelp("example: [...arr1, ...arr2] combines two arrays")
		} else {
			richErr.WithNote("spread operator must be used inside array literals")
		}

	} else if strings.Contains(msg, "unusable as hash key") {
		richErr.WithCode("E0013")
		richErr.WithNote("hash keys must be hashable types: strings, integers, or booleans")
		richErr.WithHelp("use a string, integer, or boolean as the hash key")

	} else if strings.Contains(msg, "reduce of empty array with no initial value") {
		richErr.WithCode("E0021")
		richErr.WithNote("reduce() needs an initial value when the array is empty")
		richErr.WithHelp("provide an initial value: reduce([], fn, 0)")

	} else if strings.Contains(msg, "division by zero") {
		richErr.WithCode("E0023")
		richErr.WithNote("cannot divide by zero")
		richErr.WithHelp("check the divisor is not zero before dividing")
	}

	return richErr.Format()
}
