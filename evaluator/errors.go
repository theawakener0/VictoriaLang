package evaluator

import (
	"fmt"
	"strings"
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
		_ = richErr.WithCode("E0001")

		// Handle typed variable declaration errors
		if strings.Contains(msg, "cannot assign") && strings.Contains(msg, "to variable of type") {
			_ = richErr.WithNote("Victoria supports optional static typing for type safety")
			parts := strings.Split(msg, "cannot assign ")
			if len(parts) > 1 {
				typeParts := strings.Split(parts[1], " to variable of type ")
				if len(typeParts) == 2 {
					actualType := strings.TrimSpace(typeParts[0])
					expectedType := strings.TrimSpace(typeParts[1])
					_ = richErr.WithHelp(fmt.Sprintf("either change the value to a %s, or change the type annotation to :%s", expectedType, actualType))
				}
			}
		} else if strings.Contains(msg, "cannot assign") && strings.Contains(msg, "to constant of type") {
			_ = richErr.WithNote("constants with type annotations must match the declared type")
			_ = richErr.WithHelp("ensure the value matches the declared type, or remove the type annotation")
		} else if strings.Contains(msg, "for parameter") {
			// Function parameter type mismatch
			_ = richErr.WithNote("function parameters with type annotations enforce type checking at runtime")
			if strings.Contains(msg, "expected int") {
				_ = richErr.WithHelp("pass an integer value, or use int() to convert: int(value)")
			} else if strings.Contains(msg, "expected string") {
				_ = richErr.WithHelp("pass a string value, or use string() to convert: string(value)")
			} else if strings.Contains(msg, "expected float") {
				_ = richErr.WithHelp("pass a float value - integers are automatically converted to float")
			} else if strings.Contains(msg, "expected bool") {
				_ = richErr.WithHelp("pass a boolean value (true or false)")
			} else if strings.Contains(msg, "expected array") {
				_ = richErr.WithHelp("pass an array value: [element1, element2, ...]")
			} else {
				_ = richErr.WithHelp("ensure the argument matches the expected parameter type")
			}
		} else if strings.Contains(msg, "return type mismatch") {
			// Return type mismatch
			_ = richErr.WithNote("functions with return type annotations must return matching types")
			if strings.Contains(msg, "expected int") {
				_ = richErr.WithHelp("return an integer value, or change the return type annotation")
			} else if strings.Contains(msg, "expected string") {
				_ = richErr.WithHelp("return a string value, or change the return type annotation")
			} else if strings.Contains(msg, "expected bool") {
				_ = richErr.WithHelp("return true or false, or change the return type annotation")
			} else if strings.Contains(msg, "expected void") {
				_ = richErr.WithHelp("remove the return statement, or change the return type from 'void'")
			} else {
				_ = richErr.WithHelp("ensure the return value matches the declared return type")
			}
		} else {
			_ = richErr.WithNote("Victoria is dynamically typed, but operators require compatible types")
		}

		if strings.Contains(msg, "STRING") && strings.Contains(msg, "INTEGER") {
			_ = richErr.WithNote("strings and integers cannot be combined directly with arithmetic operators")
			_ = richErr.WithHelp("use string() to convert integers to strings: \"text\" + string(42)")
		}
		if strings.Contains(msg, "STRING") && strings.Contains(msg, "FLOAT") {
			_ = richErr.WithNote("strings and floats cannot be combined directly")
			_ = richErr.WithHelp("use string() to convert floats to strings: \"value: \" + string(3.14)")
		}
		if strings.Contains(msg, "BOOLEAN") {
			_ = richErr.WithNote("booleans cannot be used in arithmetic operations")
			if strings.Contains(msg, "INTEGER") {
				_ = richErr.WithHelp("use int() to convert boolean to integer: int(true) returns 1")
			} else {
				_ = richErr.WithHelp("use string() to convert boolean to string: string(true) returns \"true\"")
			}
		}
		if strings.Contains(msg, "ARRAY") {
			_ = richErr.WithNote("arrays can only be concatenated with other arrays using '+'")
			_ = richErr.WithHelp("use push() to add elements: push(arr, element)")
		}
		if strings.Contains(msg, "HASH") {
			_ = richErr.WithNote("hashes do not support arithmetic operators")
			_ = richErr.WithHelp("access hash values with hash[\"key\"] or hash.key syntax")
		}

	} else if strings.Contains(msg, "identifier not found") {
		name := strings.TrimPrefix(msg, "identifier not found: ")
		_ = richErr.WithCode("E0002")
		_ = richErr.WithNote("variables must be declared before use with 'let' or 'const'")

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
			_ = richErr.WithHelp(suggestion)
		} else if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
			_ = richErr.WithHelp(fmt.Sprintf("'%s' looks like a type name; did you mean to create an instance?", name))
			_ = richErr.WithNote("use struct instantiation: new StructName { field: value }")
		} else {
			_ = richErr.WithHelp(fmt.Sprintf("declare with: let %s = <value>", name))
			_ = richErr.WithNote("for constants, use: const " + name + " = <value>")
		}

	} else if strings.Contains(msg, "index operator not supported") {
		_ = richErr.WithCode("E0006")
		_ = richErr.WithNote("indexing is only supported for arrays, strings, and hashes")

		if strings.Contains(msg, "INTEGER") {
			_ = richErr.WithHelp("integers cannot be indexed; did you mean to use an array?")
		} else if strings.Contains(msg, "NULL") {
			_ = richErr.WithHelp("cannot index null; ensure the value is initialized")
		} else {
			_ = richErr.WithHelp("only arrays, strings, and hashes support [] indexing")
		}

	} else if strings.Contains(msg, "not a function") {
		_ = richErr.WithCode("E0005")

		if strings.Contains(msg, "INTEGER") {
			_ = richErr.WithNote("integers cannot be called as functions")
			_ = richErr.WithHelp("remove the parentheses, or did you mean to use a function?")
		} else if strings.Contains(msg, "ARRAY") {
			_ = richErr.WithNote("arrays cannot be called as functions")
			_ = richErr.WithHelp("use array[index] to access elements, not array(index)")
		} else if strings.Contains(msg, "HASH") {
			_ = richErr.WithNote("hashes cannot be called as functions")
			_ = richErr.WithHelp("use hash[\"key\"] or hash.key to access values")
		} else if strings.Contains(msg, "nil") || strings.Contains(msg, "NULL") {
			_ = richErr.WithNote("attempted to call null as a function")
			_ = richErr.WithHelp("ensure the variable is assigned a function before calling")
		} else {
			_ = richErr.WithNote("only functions and builtin functions can be called")
		}

	} else if strings.Contains(msg, "unknown operator") {
		_ = richErr.WithCode("E0003")

		if strings.Contains(msg, "STRING") && strings.Contains(msg, "-") {
			_ = richErr.WithNote("strings only support the '+' operator for concatenation")
		} else if strings.Contains(msg, "BOOLEAN") {
			_ = richErr.WithNote("booleans only support comparison operators (==, !=)")
			_ = richErr.WithHelp("use 'and', 'or', '!' for boolean logic")
		} else {
			_ = richErr.WithNote("this operator is not supported for the given types")
		}

	} else if strings.Contains(msg, "variable not defined") {
		name := strings.TrimPrefix(msg, "variable not defined: ")
		_ = richErr.WithCode("E0002")
		_ = richErr.WithNote("cannot modify a variable that hasn't been declared")
		_ = richErr.WithHelp(fmt.Sprintf("declare first: let %s = <initial_value>", name))

	} else if strings.Contains(msg, "struct not found") {
		name := strings.TrimPrefix(msg, "struct not found: ")
		_ = richErr.WithCode("E0009")
		_ = richErr.WithNote("structs must be defined before instantiation")
		_ = richErr.WithHelp(fmt.Sprintf("define the struct first: struct %s { field1, field2 }", name))

	} else if strings.Contains(msg, "wrong number of arguments") {
		_ = richErr.WithCode("E0010")

		// Check if this is from a typed function
		if strings.Contains(msg, "expected") && strings.Contains(msg, "got") {
			_ = richErr.WithNote("typed functions enforce parameter count at runtime")
			_ = richErr.WithHelp("check the function definition for the correct number of parameters")
		} else {
			_ = richErr.WithNote("function called with incorrect number of arguments")
		}

		if strings.Contains(msg, "len") {
			_ = richErr.WithHelp("len() takes exactly 1 argument: len(array) or len(string)")
		} else if strings.Contains(msg, "push") {
			_ = richErr.WithHelp("push() takes exactly 2 arguments: push(array, element)")
		} else if strings.Contains(msg, "map") {
			_ = richErr.WithHelp("map() takes exactly 2 arguments: map(array, fn)")
		} else if strings.Contains(msg, "filter") {
			_ = richErr.WithHelp("filter() takes exactly 2 arguments: filter(array, predicate)")
		} else if strings.Contains(msg, "reduce") {
			_ = richErr.WithHelp("reduce() takes 2 or 3 arguments: reduce(array, fn) or reduce(array, fn, initial)")
		}

	} else if strings.Contains(msg, "spread operator") {
		_ = richErr.WithCode("E0012")
		if strings.Contains(msg, "requires an array") {
			_ = richErr.WithNote("the spread operator (...) can only expand arrays")
			_ = richErr.WithHelp("example: [...arr1, ...arr2] combines two arrays")
		} else {
			_ = richErr.WithNote("spread operator must be used inside array literals")
		}

	} else if strings.Contains(msg, "unusable as hash key") {
		_ = richErr.WithCode("E0013")
		_ = richErr.WithNote("hash keys must be hashable types: strings, integers, or booleans")
		_ = richErr.WithHelp("use a string, integer, or boolean as the hash key")

	} else if strings.Contains(msg, "reduce of empty array with no initial value") {
		_ = richErr.WithCode("E0021")
		_ = richErr.WithNote("reduce() needs an initial value when the array is empty")
		_ = richErr.WithHelp("provide an initial value: reduce([], fn, 0)")

	} else if strings.Contains(msg, "division by zero") {
		_ = richErr.WithCode("E0023")
		_ = richErr.WithNote("cannot divide by zero")
		_ = richErr.WithHelp("check the divisor is not zero before dividing")

	} else if strings.Contains(msg, "expected type") || strings.Contains(msg, "invalid type") {
		// Type annotation parsing errors
		_ = richErr.WithCode("E0030")
		_ = richErr.WithNote("type annotations specify the expected type of a value")
		_ = richErr.WithHelp("valid types: int, float, string, bool, char, array, map, any, void")

	} else if strings.Contains(msg, "cannot use type") {
		// Type usage errors
		_ = richErr.WithCode("E0031")
		_ = richErr.WithNote("some types have restrictions on where they can be used")
		if strings.Contains(msg, "void") {
			_ = richErr.WithHelp("'void' can only be used as a function return type")
		}
	}

	return richErr.Format()
}
