package evaluator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"victoria/ast"
	"victoria/errors"
	"victoria/lexer"
	"victoria/object"
	"victoria/parser"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func init() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
}

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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.IncludeStatement:
		return evalIncludeStatement(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.SetConst(node.Name.Value, val)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return evalStringLiteral(node.Value, env)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		// Handle prefix ++ and -- (they need to modify the variable)
		if node.Operator == "++" || node.Operator == "--" {
			return evalPrefixIncDec(node, env)
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		result := evalPrefixExpression(node.Operator, right)
		// Add location info to errors
		if errObj, ok := result.(*object.Error); ok && errObj.Line == 0 {
			errObj.Line = node.Token.Line
			errObj.Column = node.Token.Column
			errObj.EndColumn = node.Token.EndColumn
		}
		return result

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.InfixExpression:
		// Special handling for assignment if we treat it as infix
		// But we don't have assignment expression in AST yet, only LetStatement.
		// Wait, I said I would handle reassignment `x = 5`.
		// In parser, I didn't implement `parseAssignmentExpression` specifically, but `parseInfixExpression` handles `=`.
		// So `x = 5` becomes `InfixExpression(x, =, 5)`.
		if node.Operator == "=" || node.Operator == "+=" || node.Operator == "-=" || node.Operator == "*=" || node.Operator == "/=" || node.Operator == "%=" {
			return evalAssignmentExpression(node, env)
		}

		// Special handling for dot operator (field access or method call preparation)
		if node.Operator == "." {
			return evalDotExpression(node, env)
		}

		// Short-circuit evaluation for && and ||
		if node.Operator == "&&" || node.Operator == "and" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if !isTruthy(left) {
				return FALSE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}

		if node.Operator == "||" || node.Operator == "or" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if isTruthy(left) {
				return TRUE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		result := evalInfixExpression(node.Operator, left, right)
		// Add location info to errors
		if errObj, ok := result.(*object.Error); ok && errObj.Line == 0 {
			errObj.Line = node.Token.Line
			errObj.Column = node.Token.Column
			errObj.EndColumn = node.Token.EndColumn
		}
		return result

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.TryStatement:
		return evalTryStatement(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.ArrowFunction:
		params := node.Parameters
		body := node.Body
		return &object.ArrowFunction{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		result := applyFunction(function, args)
		// Add location info to errors
		if errObj, ok := result.(*object.Error); ok && errObj.Line == 0 {
			errObj.Line = node.Token.Line
			errObj.Column = node.Token.Column
			errObj.EndColumn = node.Token.EndColumn
		}
		return result

	case *ast.ArrayLiteral:
		elements := evalArrayElements(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		result := evalIndexExpression(left, index)
		// Add location info to errors
		if errObj, ok := result.(*object.Error); ok && errObj.Line == 0 {
			errObj.Line = node.Token.Line
			errObj.Column = node.Token.Column
			errObj.EndColumn = node.Token.EndColumn
		}
		return result

	case *ast.SliceExpression:
		return evalSliceExpression(node, env)

	case *ast.SpreadExpression:
		// Spread expressions are handled in array literals
		return newError("spread operator can only be used in array literals")

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.StructLiteral:
		// Register struct definition in env
		s := &object.Struct{Name: node.Name.Value}
		for _, f := range node.Fields {
			s.Fields = append(s.Fields, f.Value)
		}
		env.Set(node.Name.Value, s)
		return NULL

	case *ast.StructInstantiation:
		return evalStructInstantiation(node, env)

	case *ast.MethodDefinition:
		// Register method in env as "StructName.MethodName"
		// We store it as a Function object
		fn := &object.Function{Parameters: node.Parameters, Env: env, Body: node.Body}
		key := node.StructName.Value + "." + node.MethodName.Value
		env.Set(key, fn)
		return NULL

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.ForInIndexExpression:
		return evalForInIndexExpression(node, env)

	case *ast.CForExpression:
		return evalCForExpression(node, env)

	case *ast.BreakStatement:
		return &object.Break{}

	case *ast.ContinueStatement:
		return &object.Continue{}

	case *ast.SwitchExpression:
		return evalSwitchExpression(node, env)

	case *ast.TernaryExpression:
		return evalTernaryExpression(node, env)

	case *ast.RangeExpression:
		return evalRangeExpression(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	// Create a new scope for the block
	blockEnv := object.NewEnclosedEnvironment(env)

	for _, statement := range block.Statements {
		result = Eval(statement, blockEnv)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ || rt == object.BREAK_OBJ || rt == object.CONTINUE_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}
	return newError("unknown operator: -%s", right.Type())
}

func evalStringLiteral(s string, env *object.Environment) object.Object {
	// First process escape sequences
	s = processEscapeSequences(s)

	// Handle string interpolation: ${expr}
	result := ""
	i := 0
	for i < len(s) {
		if i+1 < len(s) && s[i] == '$' && s[i+1] == '{' {
			// Find matching }
			j := i + 2
			depth := 1
			for j < len(s) && depth > 0 {
				if s[j] == '{' {
					depth++
				} else if s[j] == '}' {
					depth--
				}
				j++
			}
			if depth == 0 {
				// Extract expression
				exprStr := s[i+2 : j-1]
				// Parse and evaluate the expression
				l := lexer.New(exprStr)
				p := parser.New(l)
				program := p.ParseProgram()
				if len(p.Errors()) > 0 {
					return newError("string interpolation parse error: %s", p.Errors()[0])
				}
				if len(program.Statements) == 0 {
					result += ""
				} else {
					// Evaluate first statement/expression
					val := Eval(program.Statements[0], env)
					if isError(val) {
						return val
					}
					if val != nil {
						result += val.Inspect()
					}
				}
				i = j
			} else {
				result += string(s[i])
				i++
			}
		} else {
			result += string(s[i])
			i++
		}
	}
	return &object.String{Value: result}
}

func processEscapeSequences(s string) string {
	result := ""
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result += "\n"
				i += 2
			case 't':
				result += "\t"
				i += 2
			case 'r':
				result += "\r"
				i += 2
			case '\\':
				result += "\\"
				i += 2
			case '"':
				result += "\""
				i += 2
			case '$':
				result += "$"
				i += 2
			default:
				result += string(s[i])
				i++
			}
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, &object.Float{Value: float64(left.(*object.Integer).Value)}, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalFloatInfixExpression(operator, left, &object.Float{Value: float64(right.(*object.Integer).Value)})
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalTryStatement(node *ast.TryStatement, env *object.Environment) object.Object {
	result := Eval(node.Block, env)

	if isError(result) {
		if node.CatchBlock != nil {
			catchEnv := object.NewEnclosedEnvironment(env)
			if node.CatchVar != nil {
				msg := result.(*object.Error).Message
				catchEnv.Set(node.CatchVar.Value, &object.String{Value: msg})
			}
			return Eval(node.CatchBlock, catchEnv)
		}
		return NULL
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newErrorWithLocation("identifier not found: "+node.Value, node.Token.Line, node.Token.Column, node.Token.EndColumn)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
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

	// ═══════════════════════════════════════════════════════════════════════════
	// TYPE MISMATCH ERRORS
	// ═══════════════════════════════════════════════════════════════════════════
	if strings.Contains(msg, "type mismatch") {
		richErr.WithCode("E0001")
		richErr.WithNote("Victoria is dynamically typed, but operators require compatible types")

		// String + Integer
		if strings.Contains(msg, "STRING") && strings.Contains(msg, "INTEGER") {
			richErr.WithNote("strings and integers cannot be combined directly with arithmetic operators")
			richErr.WithHelp("use string() to convert integers to strings: \"text\" + string(42)")
		}
		// String + Float
		if strings.Contains(msg, "STRING") && strings.Contains(msg, "FLOAT") {
			richErr.WithNote("strings and floats cannot be combined directly")
			richErr.WithHelp("use string() to convert floats to strings: \"value: \" + string(3.14)")
		}
		// Boolean + Integer/String
		if strings.Contains(msg, "BOOLEAN") {
			richErr.WithNote("booleans cannot be used in arithmetic operations")
			if strings.Contains(msg, "INTEGER") {
				richErr.WithHelp("use int() to convert boolean to integer: int(true) returns 1")
			} else {
				richErr.WithHelp("use string() to convert boolean to string: string(true) returns \"true\"")
			}
		}
		// Array operations
		if strings.Contains(msg, "ARRAY") {
			richErr.WithNote("arrays can only be concatenated with other arrays using '+'")
			richErr.WithHelp("use push() to add elements: push(arr, element)")
		}
		// Hash operations
		if strings.Contains(msg, "HASH") {
			richErr.WithNote("hashes do not support arithmetic operators")
			richErr.WithHelp("access hash values with hash[\"key\"] or hash.key syntax")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// IDENTIFIER NOT FOUND ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "identifier not found") {
		name := strings.TrimPrefix(msg, "identifier not found: ")
		richErr.WithCode("E0002")
		richErr.WithNote("variables must be declared before use with 'let' or 'const'")

		// Common typos and suggestions - values are complete help messages
		commonBuiltins := map[string]string{
			// Output functions
			"println": "use 'print' instead - Victoria uses print() for output",
			"printf":  "use 'format' instead - Victoria uses format() for formatted strings",
			"console": "use 'print' instead - Victoria uses print() for output",
			"log":     "use 'print' instead - Victoria uses print() for output",
			"echo":    "use 'print' instead - Victoria uses print() for output",
			"puts":    "use 'print' instead - Victoria uses print() for output",
			// Type conversions
			"str":      "use 'string' instead - Victoria uses string() for type conversion",
			"toString": "use 'string' instead - Victoria uses string() for type conversion",
			"toInt":    "use 'int' instead - Victoria uses int() for type conversion",
			"parseInt": "use 'int' instead - Victoria uses int() for type conversion",
			// Length/size functions
			"size":   "use 'len' instead - Victoria uses len() for collection length",
			"length": "use 'len' instead - Victoria uses len() for collection length",
			"count":  "use 'len' instead - Victoria uses len() for collection length",
			// Array operations
			"append": "use 'push' instead - Victoria uses push(array, element)",
			"add":    "use 'push' instead - Victoria uses push(array, element)",
			"remove": "use 'pop' instead - Victoria uses pop(array) to remove last element",
			"delete": "use filter() to create a new array without elements",
			// String operations
			"substr":    "use string slicing instead: str[start:end]",
			"substring": "use string slicing instead: str[start:end]",
			// Iteration
			"forEach": "use a for-in loop: for item in array { ... }",
			"foreach": "use a for-in loop: for item in array { ... }",
			// Function-like values
			"map_func":  "use 'map' - Victoria has map(array, fn)",
			"filter_fn": "use 'filter' - Victoria has filter(array, fn)",
			"reduce_fn": "use 'reduce' - Victoria has reduce(array, fn, initial)",
			// Keywords from other languages
			"nil":       "use 'null' instead - Victoria uses null for no value",
			"none":      "use 'null' instead - Victoria uses null for no value",
			"None":      "use 'null' instead - Victoria uses null for no value",
			"undefined": "use 'null' instead - Victoria uses null for no value",
			"true_val":  "use 'true' - it's a built-in boolean literal",
			"false_val": "use 'false' - it's a built-in boolean literal",
			// Function definition alternatives
			"fn":       "use 'define' instead - Victoria uses 'define' to create functions: let f = define(x) { x * 2 }",
			"func":     "use 'define' instead - Victoria uses 'define' to create functions: let f = define(x) { x * 2 }",
			"function": "use 'define' instead - Victoria uses 'define' to create functions: let f = define(x) { x * 2 }",
			"lambda":   "use 'define' instead - Victoria uses 'define' to create functions: let f = define(x) { x * 2 }",
			"def":      "use 'define' instead - Victoria uses 'define' to create functions: let f = define(x) { x * 2 }",
			// Variable declaration
			"var": "use 'let' instead - Victoria uses 'let' for variable declaration",
			// Control flow
			"elif":   "use 'else if' - Victoria uses 'else if' not 'elif'",
			"elsif":  "use 'else if' - Victoria uses 'else if' not 'elsif'",
			"elseif": "use 'else if' (with space) - Victoria uses 'else if'",
			// Operators as identifiers (user error)
			"and": "'and' is an operator, not a function - use it between expressions: a and b",
			"or":  "'or' is an operator, not a function - use it between expressions: a or b",
			"not": "use '!' for negation - Victoria uses !value not not(value)",
			// OOP patterns
			"self": "Victoria uses the struct instance name directly in methods",
			"this": "Victoria uses the struct instance name directly in methods",
			// Module system
			"require": "use 'include' instead - Victoria uses include \"filename\"",
			"import":  "use 'include' instead - Victoria uses include \"filename\"",
		}

		if suggestion, ok := commonBuiltins[name]; ok {
			richErr.WithHelp(suggestion)
		} else if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
			// Possibly trying to use a type as a value
			richErr.WithHelp(fmt.Sprintf("'%s' looks like a type name; did you mean to create an instance?", name))
			richErr.WithNote("use struct instantiation: new StructName { field: value }")
		} else {
			richErr.WithHelp(fmt.Sprintf("declare with: let %s = <value>", name))
			richErr.WithNote("for constants, use: const " + name + " = <value>")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// INDEX OPERATOR ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "index operator not supported") {
		richErr.WithCode("E0006")
		richErr.WithNote("indexing is only supported for arrays, strings, and hashes")

		if strings.Contains(msg, "INTEGER") {
			richErr.WithHelp("integers cannot be indexed; did you mean to use an array?")
			richErr.WithNote("example: let arr = [1, 2, 3]; arr[0] returns 1")
		} else if strings.Contains(msg, "BOOLEAN") {
			richErr.WithHelp("booleans cannot be indexed")
		} else if strings.Contains(msg, "FUNCTION") {
			richErr.WithHelp("functions cannot be indexed; call the function first with ()")
		} else if strings.Contains(msg, "NULL") {
			richErr.WithHelp("cannot index null; ensure the value is initialized")
			richErr.WithNote("this often happens when accessing a missing hash key")
		} else {
			richErr.WithHelp("only arrays, strings, and hashes support [] indexing")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// NOT A FUNCTION ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "not a function") {
		richErr.WithCode("E0005")

		if strings.Contains(msg, "INTEGER") {
			richErr.WithNote("integers cannot be called as functions")
			richErr.WithHelp("remove the parentheses, or did you mean to use a function?")
		} else if strings.Contains(msg, "STRING") {
			richErr.WithNote("strings cannot be called as functions")
			richErr.WithHelp("use string methods like split(), upper(), lower() instead")
		} else if strings.Contains(msg, "BOOLEAN") {
			richErr.WithNote("booleans cannot be called as functions")
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
			richErr.WithHelp("check that the variable contains a function")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// UNKNOWN OPERATOR ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "unknown operator") {
		richErr.WithCode("E0003")

		if strings.Contains(msg, "STRING") && strings.Contains(msg, "-") {
			richErr.WithNote("strings only support the '+' operator for concatenation")
			richErr.WithHelp("to compare strings, use == or !=")
		} else if strings.Contains(msg, "STRING") && (strings.Contains(msg, "*") || strings.Contains(msg, "/")) {
			richErr.WithNote("strings do not support multiplication or division")
			richErr.WithHelp("use string concatenation (+) or string functions instead")
		} else if strings.Contains(msg, "BOOLEAN") {
			richErr.WithNote("booleans only support comparison operators (==, !=)")
			richErr.WithHelp("use 'and', 'or', '!' for boolean logic")
		} else if strings.Contains(msg, "ARRAY") {
			richErr.WithNote("arrays only support '+' for concatenation and == for comparison")
			richErr.WithHelp("use array functions: push(), pop(), first(), last(), rest()")
		} else if strings.Contains(msg, "FUNCTION") {
			richErr.WithNote("functions can only be compared with == or !=")
			richErr.WithHelp("call the function first to operate on its return value")
		} else {
			richErr.WithNote("this operator is not supported for the given types")
			richErr.WithHelp("check the Victoria language documentation for supported operators")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// VARIABLE NOT DEFINED ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "variable not defined") {
		name := strings.TrimPrefix(msg, "variable not defined: ")
		richErr.WithCode("E0002")
		richErr.WithNote("cannot modify a variable that hasn't been declared")
		richErr.WithHelp(fmt.Sprintf("declare first: let %s = <initial_value>", name))

		// ═══════════════════════════════════════════════════════════════════════════
		// STRUCT ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "struct not found") {
		name := strings.TrimSuffix(strings.TrimPrefix(msg, "struct not found: "), "")
		richErr.WithCode("E0009")
		richErr.WithNote("structs must be defined before instantiation")
		richErr.WithHelp(fmt.Sprintf("define the struct first: struct %s { field1, field2 }", name))

	} else if strings.Contains(msg, "not a struct") {
		richErr.WithCode("E0009")
		richErr.WithNote("'new' keyword requires a struct type")
		richErr.WithHelp("use 'new StructName { field: value }' syntax")

		// ═══════════════════════════════════════════════════════════════════════════
		// PROPERTY ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "property not found") || strings.Contains(msg, "property or method not found") {
		richErr.WithCode("E0008")
		richErr.WithNote("the property or method does not exist on this type")
		richErr.WithHelp("check the spelling or use keys() to see available hash keys")

	} else if strings.Contains(msg, "expected identifier after dot") {
		richErr.WithCode("E0100")
		richErr.WithNote("dot notation requires a property name")
		richErr.WithHelp("use: object.propertyName or object.methodName()")

		// ═══════════════════════════════════════════════════════════════════════════
		// SLICE ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "slice index must be an integer") {
		richErr.WithCode("E0011")
		richErr.WithNote("slice indices must be integers")
		richErr.WithHelp("use integer values: array[0:5] or string[2:10]")

	} else if strings.Contains(msg, "slice operator not supported") {
		richErr.WithCode("E0011")
		richErr.WithNote("slicing is only supported for arrays and strings")
		richErr.WithHelp("syntax: array[start:end] or string[start:end]")

		// ═══════════════════════════════════════════════════════════════════════════
		// SPREAD OPERATOR ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "spread operator") {
		richErr.WithCode("E0012")
		if strings.Contains(msg, "requires an array") {
			richErr.WithNote("the spread operator (...) can only expand arrays")
			richErr.WithHelp("example: [...arr1, ...arr2] combines two arrays")
		} else {
			richErr.WithNote("spread operator must be used inside array literals")
			richErr.WithHelp("example: let combined = [...array1, ...array2]")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// HASH KEY ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "unusable as hash key") {
		richErr.WithCode("E0013")
		richErr.WithNote("hash keys must be hashable types: strings, integers, or booleans")
		if strings.Contains(msg, "ARRAY") {
			richErr.WithHelp("arrays cannot be hash keys; use a string identifier instead")
		} else if strings.Contains(msg, "HASH") {
			richErr.WithHelp("hashes cannot be hash keys; use a string identifier instead")
		} else if strings.Contains(msg, "FUNCTION") {
			richErr.WithHelp("functions cannot be hash keys; use a string identifier instead")
		} else {
			richErr.WithHelp("use a string, integer, or boolean as the hash key")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// WRONG NUMBER OF ARGUMENTS ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "wrong number of arguments") {
		richErr.WithCode("E0010")
		richErr.WithNote("function called with incorrect number of arguments")

		// Extract function name and provide specific help
		if strings.Contains(msg, "len") || strings.Contains(msg, "`len`") {
			richErr.WithHelp("len() takes exactly 1 argument: len(array) or len(string)")
		} else if strings.Contains(msg, "push") || strings.Contains(msg, "`push`") {
			richErr.WithHelp("push() takes exactly 2 arguments: push(array, element)")
		} else if strings.Contains(msg, "pop") || strings.Contains(msg, "`pop`") {
			richErr.WithHelp("pop() takes exactly 1 argument: pop(array)")
		} else if strings.Contains(msg, "split") || strings.Contains(msg, "`split`") {
			richErr.WithHelp("split() takes exactly 2 arguments: split(string, delimiter)")
		} else if strings.Contains(msg, "join") || strings.Contains(msg, "`join`") {
			richErr.WithHelp("join() takes exactly 2 arguments: join(array, separator)")
		} else if strings.Contains(msg, "map") || strings.Contains(msg, "`map`") {
			richErr.WithHelp("map() takes exactly 2 arguments: map(array, fn)")
		} else if strings.Contains(msg, "filter") || strings.Contains(msg, "`filter`") {
			richErr.WithHelp("filter() takes exactly 2 arguments: filter(array, predicate)")
		} else if strings.Contains(msg, "reduce") || strings.Contains(msg, "`reduce`") {
			richErr.WithHelp("reduce() takes 2 or 3 arguments: reduce(array, fn) or reduce(array, fn, initial)")
		} else if strings.Contains(msg, "range") || strings.Contains(msg, "`range`") {
			richErr.WithHelp("range() takes 1-3 arguments: range(end), range(start, end), or range(start, end, step)")
		} else if strings.Contains(msg, "format") || strings.Contains(msg, "`format`") {
			richErr.WithHelp("format() takes at least 1 argument: format(template, ...values)")
		} else {
			richErr.WithHelp("check the function signature for the correct number of arguments")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// ARGUMENT TYPE ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "argument to") && strings.Contains(msg, "must be") {
		richErr.WithCode("E0014")
		richErr.WithNote("function received an argument of the wrong type")

		if strings.Contains(msg, "must be ARRAY") {
			richErr.WithHelp("pass an array: [1, 2, 3] or a variable containing an array")
		} else if strings.Contains(msg, "must be STRING") {
			richErr.WithHelp("pass a string: \"text\" or a variable containing a string")
		} else if strings.Contains(msg, "must be INTEGER") {
			richErr.WithHelp("pass an integer: 42 or a variable containing an integer")
		} else if strings.Contains(msg, "must be FUNCTION") {
			richErr.WithHelp("pass a function: fn(x) { x * 2 } or a named function")
		} else if strings.Contains(msg, "must be HASH") {
			richErr.WithHelp("pass a hash: {\"key\": \"value\"} or a variable containing a hash")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// ITERATION ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "not iterable") {
		richErr.WithCode("E0015")
		richErr.WithNote("for-in loops require an iterable value")

		if strings.Contains(msg, "INTEGER") {
			richErr.WithHelp("use range() to iterate over numbers: for i in range(10) { ... }")
		} else if strings.Contains(msg, "STRING") {
			// Strings ARE iterable in Victoria, so this shouldn't happen, but just in case
			richErr.WithHelp("strings are iterable character by character")
		} else {
			richErr.WithHelp("iterable types: arrays, strings, hashes, and ranges")
		}

		// ═══════════════════════════════════════════════════════════════════════════
		// RANGE ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "range") && strings.Contains(msg, "must be") {
		richErr.WithCode("E0016")
		richErr.WithNote("range values must be integers")
		richErr.WithHelp("use: range(10), range(1, 10), or range(0, 100, 5)")

	} else if strings.Contains(msg, "range step cannot be zero") {
		richErr.WithCode("E0016")
		richErr.WithNote("step value determines the increment between range values")
		richErr.WithHelp("use a positive step for ascending: range(0, 10, 2)")

		// ═══════════════════════════════════════════════════════════════════════════
		// PARSE ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "could not parse") && strings.Contains(msg, "as integer") {
		richErr.WithCode("E0017")
		richErr.WithNote("string could not be converted to an integer")
		richErr.WithHelp("ensure the string contains only numeric characters: \"123\" not \"12.3\" or \"abc\"")

		// ═══════════════════════════════════════════════════════════════════════════
		// ASSIGNMENT ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "assignment to non-identifier") {
		richErr.WithCode("E0018")
		richErr.WithNote("can only assign to variable names")
		richErr.WithHelp("use: variableName = value, not expressions like (a + b) = value")

		// ═══════════════════════════════════════════════════════════════════════════
		// PREFIX/POSTFIX ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "operator on non-identifier") {
		richErr.WithCode("E0019")
		richErr.WithNote("increment/decrement operators require a variable")
		richErr.WithHelp("use: i++ or ++i where i is a variable name")

		// ═══════════════════════════════════════════════════════════════════════════
		// DOT OPERATOR ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "dot operator not supported") {
		richErr.WithCode("E0020")
		richErr.WithNote("dot notation is only for hashes, structs, and objects with methods")
		richErr.WithHelp("use hash[\"key\"] for dynamic keys, or define a struct for complex types")

		// ═══════════════════════════════════════════════════════════════════════════
		// REDUCE ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "reduce of empty array with no initial value") {
		richErr.WithCode("E0021")
		richErr.WithNote("reduce() needs an initial value when the array is empty")
		richErr.WithHelp("provide an initial value: reduce([], fn, 0)")

		// ═══════════════════════════════════════════════════════════════════════════
		// JOIN ERRORS
		// ═══════════════════════════════════════════════════════════════════════════
	} else if strings.Contains(msg, "array elements must be strings") {
		richErr.WithCode("E0022")
		richErr.WithNote("join() requires all array elements to be strings")
		richErr.WithHelp("convert elements first: map(arr, fn(x) { string(x) })")
	}

	return richErr.Format()
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// evalArrayElements evaluates array elements, handling spread expressions
func evalArrayElements(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		// Check if it's a spread expression
		if spread, ok := e.(*ast.SpreadExpression); ok {
			evaluated := Eval(spread.Right, env)
			if isError(evaluated) {
				return []object.Object{evaluated}
			}
			// The spread expression must evaluate to an array
			if arr, ok := evaluated.(*object.Array); ok {
				result = append(result, arr.Elements...)
			} else {
				return []object.Object{newError("spread operator requires an array, got %s", evaluated.Type())}
			}
		} else {
			evaluated := Eval(e, env)
			if isError(evaluated) {
				return []object.Object{evaluated}
			}
			result = append(result, evaluated)
		}
	}

	return result
}

// evalSliceExpression evaluates arr[start:end] slice expressions
func evalSliceExpression(node *ast.SliceExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	var startIdx, endIdx int64

	switch obj := left.(type) {
	case *object.Array:
		length := int64(len(obj.Elements))

		// Evaluate start index
		if node.Start != nil {
			startVal := Eval(node.Start, env)
			if isError(startVal) {
				return startVal
			}
			if intVal, ok := startVal.(*object.Integer); ok {
				startIdx = intVal.Value
				if startIdx < 0 {
					startIdx = length + startIdx // negative indexing
				}
			} else {
				return newError("slice index must be an integer, got %s", startVal.Type())
			}
		} else {
			startIdx = 0
		}

		// Evaluate end index
		if node.End != nil {
			endVal := Eval(node.End, env)
			if isError(endVal) {
				return endVal
			}
			if intVal, ok := endVal.(*object.Integer); ok {
				endIdx = intVal.Value
				if endIdx < 0 {
					endIdx = length + endIdx // negative indexing
				}
			} else {
				return newError("slice index must be an integer, got %s", endVal.Type())
			}
		} else {
			endIdx = length
		}

		// Bounds checking
		if startIdx < 0 {
			startIdx = 0
		}
		if endIdx > length {
			endIdx = length
		}
		if startIdx > endIdx {
			return &object.Array{Elements: []object.Object{}}
		}

		// Create new slice
		newElements := make([]object.Object, endIdx-startIdx)
		copy(newElements, obj.Elements[startIdx:endIdx])
		return &object.Array{Elements: newElements}

	case *object.String:
		length := int64(len(obj.Value))

		// Evaluate start index
		if node.Start != nil {
			startVal := Eval(node.Start, env)
			if isError(startVal) {
				return startVal
			}
			if intVal, ok := startVal.(*object.Integer); ok {
				startIdx = intVal.Value
				if startIdx < 0 {
					startIdx = length + startIdx
				}
			} else {
				return newError("slice index must be an integer, got %s", startVal.Type())
			}
		} else {
			startIdx = 0
		}

		// Evaluate end index
		if node.End != nil {
			endVal := Eval(node.End, env)
			if isError(endVal) {
				return endVal
			}
			if intVal, ok := endVal.(*object.Integer); ok {
				endIdx = intVal.Value
				if endIdx < 0 {
					endIdx = length + endIdx
				}
			} else {
				return newError("slice index must be an integer, got %s", endVal.Type())
			}
		} else {
			endIdx = length
		}

		// Bounds checking
		if startIdx < 0 {
			startIdx = 0
		}
		if endIdx > length {
			endIdx = length
		}
		if startIdx > endIdx {
			return &object.String{Value: ""}
		}

		return &object.String{Value: obj.Value[startIdx:endIdx]}

	default:
		return newError("slice operator not supported for: %s", left.Type())
	}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	if fn == nil {
		return newError("not a function: nil")
	}
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.ArrowFunction:
		extendedEnv := extendArrowFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		// Arrow functions implicitly return their body expression
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	// If it's a method call, we might need to inject 'self'.
	// But 'self' is not in parameters.
	// We need to handle 'self' injection in applyFunction or before.
	// See evalDotExpression.

	return env
}

func extendArrowFunctionEnv(fn *object.ArrowFunction, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		if i < len(args) {
			env.Set(param.Value, args[i])
		}
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalStructInstantiation(node *ast.StructInstantiation, env *object.Environment) object.Object {
	// Find struct definition
	sObj, ok := env.Get(node.Name.Value)
	if !ok {
		return newError("struct not found: %s", node.Name.Value)
	}

	sDef, ok := sObj.(*object.Struct)
	if !ok {
		return newError("not a struct: %s", node.Name.Value)
	}

	instance := &object.StructInstance{Struct: sDef, Fields: make(map[string]object.Object)}

	for fieldName, expr := range node.Fields {
		val := Eval(expr, env)
		if isError(val) {
			return val
		}
		instance.Fields[fieldName] = val
	}

	return instance
}

func evalDotExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	// left.right
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	// right should be identifier
	ident, ok := node.Right.(*ast.Identifier)
	if !ok {
		return newError("expected identifier after dot")
	}

	if left.Type() == object.HASH_OBJ {
		hash := left.(*object.Hash)
		key := &object.String{Value: ident.Value}
		hashed := key.HashKey()
		pair, ok := hash.Pairs[hashed]
		if ok {
			return pair.Value
		}
		return newError("property not found in hash: %s", ident.Value)
	}

	if left.Type() == object.INSTANCE_OBJ {
		instance := left.(*object.StructInstance)

		// Check fields
		if val, ok := instance.Fields[ident.Value]; ok {
			return val
		}

		// Check methods
		methodName := instance.Struct.Name + "." + ident.Value
		if method, ok := env.Get(methodName); ok {
			// It's a function. We need to bind 'self'.
			// We can return a bound function or just the function and let applyFunction handle it?
			// But applyFunction doesn't know about 'self'.
			// We can return a special "BoundMethod" object?
			// Or we can return the function, but we need to inject 'self' when it is called.
			// But here we are just evaluating the expression `obj.method`.
			// If the next thing is `()`, it will be called.
			// If we return the function as is, `applyFunction` will be called with arguments.
			// But `self` is missing.
			// We need to curry the function or something.
			// Let's create a closure that wraps the function and injects `self`.

			fn := method.(*object.Function)
			// Create a new environment for the closure
			closureEnv := object.NewEnclosedEnvironment(fn.Env)
			closureEnv.Set("self", instance)

			return &object.Function{Parameters: fn.Parameters, Env: closureEnv, Body: fn.Body}
		}

		return newError("property or method not found: %s", ident.Value)
	}

	return newError("dot operator not supported for: %s", left.Type())
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError("postfix operator on non-identifier")
	}

	// Check if the variable is a constant
	if env.IsConst(ident.Value) {
		return newErrorWithLocation("cannot reassign constant variable: "+ident.Value, node.Token.Line, node.Token.Column, node.Token.EndColumn)
	}

	currentVal, ok := env.Get(ident.Value)
	if !ok {
		return newError("variable not defined: %s", ident.Value)
	}

	var newVal object.Object
	one := &object.Integer{Value: 1}

	switch node.Operator {
	case "++":
		newVal = evalInfixExpression("+", currentVal, one)
	case "--":
		newVal = evalInfixExpression("-", currentVal, one)
	default:
		return newError("unknown operator: %s", node.Operator)
	}

	if isError(newVal) {
		return newVal
	}

	env.Update(ident.Value, newVal)
	return currentVal
}

func evalPrefixIncDec(node *ast.PrefixExpression, env *object.Environment) object.Object {
	ident, ok := node.Right.(*ast.Identifier)
	if !ok {
		return newError("prefix %s operator on non-identifier", node.Operator)
	}

	// Check if the variable is a constant
	if env.IsConst(ident.Value) {
		return newErrorWithLocation("cannot reassign constant variable: "+ident.Value, node.Token.Line, node.Token.Column, node.Token.EndColumn)
	}

	currentVal, ok := env.Get(ident.Value)
	if !ok {
		return newError("variable not defined: %s", ident.Value)
	}

	var newVal object.Object
	one := &object.Integer{Value: 1}

	switch node.Operator {
	case "++":
		newVal = evalInfixExpression("+", currentVal, one)
	case "--":
		newVal = evalInfixExpression("-", currentVal, one)
	default:
		return newError("unknown operator: %s", node.Operator)
	}

	if isError(newVal) {
		return newVal
	}

	env.Update(ident.Value, newVal)
	return newVal // Prefix returns the new value
}

func evalAssignmentExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	// left = right
	// left must be identifier
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError("assignment to non-identifier")
	}

	// Check if the variable is a constant
	if env.IsConst(ident.Value) {
		return newErrorWithLocation("cannot reassign constant variable: "+ident.Value, node.Token.Line, node.Token.Column, node.Token.EndColumn)
	}

	if node.Operator == "=" {
		val := Eval(node.Right, env)
		if isError(val) {
			return val
		}

		_, ok = env.Update(ident.Value, val)
		if !ok {
			return newError("variable not defined: %s", ident.Value)
		}

		return val
	}

	// Compound assignment: +=, -=, *=, /=
	currentVal, ok := env.Get(ident.Value)
	if !ok {
		return newError("variable not defined: %s", ident.Value)
	}

	rightVal := Eval(node.Right, env)
	if isError(rightVal) {
		return rightVal
	}

	var newVal object.Object
	switch node.Operator {
	case "+=":
		newVal = evalInfixExpression("+", currentVal, rightVal)
	case "-=":
		newVal = evalInfixExpression("-", currentVal, rightVal)
	case "*=":
		newVal = evalInfixExpression("*", currentVal, rightVal)
	case "/=":
		newVal = evalInfixExpression("/", currentVal, rightVal)
	case "%=":
		newVal = evalInfixExpression("%", currentVal, rightVal)
	}

	if isError(newVal) {
		return newVal
	}

	env.Update(ident.Value, newVal)
	return newVal
}

func evalWhileExpression(node *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(node.Body, env)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
			if result.Type() == object.CONTINUE_OBJ {
				continue
			}
		}
	}

	return result
}

func evalForExpression(node *ast.ForExpression, env *object.Environment) object.Object {
	// for item in list
	iterable := Eval(node.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	var elements []object.Object

	switch iterable := iterable.(type) {
	case *object.Array:
		elements = iterable.Elements
	case *object.String:
		for _, char := range iterable.Value {
			elements = append(elements, &object.String{Value: string(char)})
		}
	case *object.Hash:
		for _, pair := range iterable.Pairs {
			elements = append(elements, pair.Key)
		}
	case *object.Range:
		for i := iterable.Start; i < iterable.End; i++ {
			elements = append(elements, &object.Integer{Value: i})
		}
	default:
		return newError("not iterable: %s", iterable.Type())
	}

	var result object.Object = NULL

	for _, elem := range elements {
		loopEnv := object.NewEnclosedEnvironment(env)
		loopEnv.Set(node.Item.Value, elem)

		result = evalBlockStatement(node.Body, loopEnv)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
			if result.Type() == object.CONTINUE_OBJ {
				continue
			}
		}
	}

	return result
}

func evalForInIndexExpression(node *ast.ForInIndexExpression, env *object.Environment) object.Object {
	iterable := Eval(node.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	var result object.Object = NULL

	switch iterable := iterable.(type) {
	case *object.Array:
		for i, elem := range iterable.Elements {
			loopEnv := object.NewEnclosedEnvironment(env)
			loopEnv.Set(node.Index.Value, &object.Integer{Value: int64(i)})
			loopEnv.Set(node.Value.Value, elem)

			result = evalBlockStatement(node.Body, loopEnv)

			if result != nil {
				if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
					return result
				}
				if result.Type() == object.BREAK_OBJ {
					break
				}
				if result.Type() == object.CONTINUE_OBJ {
					continue
				}
			}
		}
	case *object.Hash:
		for _, pair := range iterable.Pairs {
			loopEnv := object.NewEnclosedEnvironment(env)
			loopEnv.Set(node.Index.Value, pair.Key)
			loopEnv.Set(node.Value.Value, pair.Value)

			result = evalBlockStatement(node.Body, loopEnv)

			if result != nil {
				if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
					return result
				}
				if result.Type() == object.BREAK_OBJ {
					break
				}
				if result.Type() == object.CONTINUE_OBJ {
					continue
				}
			}
		}
	case *object.String:
		for i, char := range iterable.Value {
			loopEnv := object.NewEnclosedEnvironment(env)
			loopEnv.Set(node.Index.Value, &object.Integer{Value: int64(i)})
			loopEnv.Set(node.Value.Value, &object.String{Value: string(char)})

			result = evalBlockStatement(node.Body, loopEnv)

			if result != nil {
				if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
					return result
				}
				if result.Type() == object.BREAK_OBJ {
					break
				}
				if result.Type() == object.CONTINUE_OBJ {
					continue
				}
			}
		}
	default:
		return newError("not iterable: %s", iterable.Type())
	}

	return result
}

func evalCForExpression(node *ast.CForExpression, env *object.Environment) object.Object {
	// for (init; cond; update)

	// Create a scope for the loop
	loopEnv := object.NewEnclosedEnvironment(env)

	// Init
	if node.Init != nil {
		initResult := Eval(node.Init, loopEnv)
		if isError(initResult) {
			return initResult
		}
	}

	var result object.Object = NULL

	for {
		// Condition
		if node.Condition != nil {
			cond := Eval(node.Condition, loopEnv)
			if isError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				break
			}
		}

		// Body
		// Body should run in the loopEnv (so it sees init vars)
		// But if body is a block, should it create ANOTHER scope?
		// If `evalBlockStatement` doesn't create scope, then it uses `loopEnv`.
		// This is correct for `for (int i=0...) { let x = 1; }`. `x` should be in `loopEnv` or inner?
		// If `evalBlockStatement` doesn't create scope, `x` is in `loopEnv`.
		// If we want `x` to be local to the block, `evalBlockStatement` should create scope.
		// Let's modify `evalBlockStatement` to create scope?
		// If I do that, `if (true) { let x = 1 } print(x)` will fail.
		// In C, `if (1) { int x = 1; }` -> x is not visible outside.
		// So yes, blocks should create scope.
		// But I need to be careful.
		// Let's stick to: `evalBlockStatement` does NOT create scope, but callers do if needed.
		// For `CFor`, the `Init` variable `i` should be visible in `Body`.
		// So `Body` should use `loopEnv`.
		// But if `Body` creates variables, they should be in a scope inside `loopEnv`?
		// Yes.
		// So `evalBlockStatement` SHOULD create a scope.
		// Let's change `evalBlockStatement` to create a scope.

		// Wait, if I change `evalBlockStatement` to create scope, then `fn` body also creates scope.
		// `applyFunction` creates an env and passes it to `Eval(body)`.
		// If `Eval(body)` (which is a BlockStatement) creates ANOTHER env, it's fine.

		// Let's change `evalBlockStatement` to create scope.
		// But wait, `evalIfExpression` calls `Eval(consequence, env)`.
		// If `consequence` is BlockStatement, it creates scope.
		// This is correct for C-like languages.

		// However, for `CFor`, we have `loopEnv` which contains `i`.
		// The body should be evaluated in a scope enclosed by `loopEnv`.

		// Let's do it here manually for now to avoid breaking other things if I'm unsure.
		bodyEnv := object.NewEnclosedEnvironment(loopEnv)
		result = evalBlockStatement(node.Body, bodyEnv)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
			if result.Type() == object.CONTINUE_OBJ {
				// fall through to update
			}
		}

		// Update
		if node.Update != nil {
			updateResult := Eval(node.Update, loopEnv)
			if isError(updateResult) {
				return updateResult
			}
		}
	}

	return result
}

func evalSwitchExpression(node *ast.SwitchExpression, env *object.Environment) object.Object {
	value := Eval(node.Value, env)
	if isError(value) {
		return value
	}

	for _, caseExpr := range node.Cases {
		caseValue := Eval(caseExpr.Value, env)
		if isError(caseValue) {
			return caseValue
		}

		// Compare values
		if compareObjects(value, caseValue) {
			return Eval(caseExpr.Body, env)
		}
	}

	if node.Default != nil {
		return Eval(node.Default, env)
	}

	return NULL
}

func compareObjects(a, b object.Object) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch a := a.(type) {
	case *object.Integer:
		return a.Value == b.(*object.Integer).Value
	case *object.String:
		return a.Value == b.(*object.String).Value
	case *object.Boolean:
		return a.Value == b.(*object.Boolean).Value
	case *object.Float:
		return a.Value == b.(*object.Float).Value
	}
	return a == b
}

func evalTernaryExpression(node *ast.TernaryExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	}
	return Eval(node.Alternative, env)
}

func evalRangeExpression(node *ast.RangeExpression, env *object.Environment) object.Object {
	start := Eval(node.Start, env)
	if isError(start) {
		return start
	}

	end := Eval(node.End, env)
	if isError(end) {
		return end
	}

	startInt, ok := start.(*object.Integer)
	if !ok {
		return newError("range start must be an integer, got %s", start.Type())
	}

	endInt, ok := end.(*object.Integer)
	if !ok {
		return newError("range end must be an integer, got %s", end.Type())
	}

	return &object.Range{Start: startInt.Value, End: endInt.Value}
}

// unwrapObject converts a Victoria object to a Go interface{}
// This is useful for functions like fmt.Sprintf that take interface{} arguments
func unwrapObject(obj object.Object) interface{} {
	switch obj := obj.(type) {
	case *object.Integer:
		return obj.Value
	case *object.Boolean:
		return obj.Value
	case *object.String:
		return obj.Value
	case *object.Null:
		return nil
	default:
		return obj.Inspect()
	}
}

// builtins is a map of built-in functions available in the Victoria language
var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"range": {
		Fn: func(args ...object.Object) object.Object {
			var start, end, step int64

			switch len(args) {
			case 1:
				// range(end) - from 0 to end
				endVal, ok := args[0].(*object.Integer)
				if !ok {
					return newError("argument to `range` must be INTEGER, got %s", args[0].Type())
				}
				start, end, step = 0, endVal.Value, 1
			case 2:
				// range(start, end) - from start to end
				startVal, ok1 := args[0].(*object.Integer)
				endVal, ok2 := args[1].(*object.Integer)
				if !ok1 || !ok2 {
					return newError("arguments to `range` must be integers")
				}
				start, end, step = startVal.Value, endVal.Value, 1
			case 3:
				// range(start, end, step)
				startVal, ok1 := args[0].(*object.Integer)
				endVal, ok2 := args[1].(*object.Integer)
				stepVal, ok3 := args[2].(*object.Integer)
				if !ok1 || !ok2 || !ok3 {
					return newError("arguments to `range` must be integers")
				}
				if stepVal.Value == 0 {
					return newError("range step cannot be zero")
				}
				start, end, step = startVal.Value, endVal.Value, stepVal.Value
			default:
				return newError("wrong number of arguments. got=%d, want=1, 2, or 3", len(args))
			}

			elements := []object.Object{}
			if step > 0 {
				for i := start; i < end; i += step {
					elements = append(elements, &object.Integer{Value: i})
				}
			} else {
				for i := start; i > end; i += step {
					elements = append(elements, &object.Integer{Value: i})
				}
			}

			return &object.Array{Elements: elements}
		},
	},
	"format": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("wrong number of arguments. got=%d, want=at least 1", len(args))
			}

			formatStr, ok := args[0].(*object.String)
			if !ok {
				return newError("argument 1 to `format` must be STRING, got %s", args[0].Type())
			}

			var fmtArgs []interface{}
			for _, arg := range args[1:] {
				fmtArgs = append(fmtArgs, unwrapObject(arg))
			}

			return &object.String{Value: fmt.Sprintf(formatStr.Value, fmtArgs...)}
		},
	},
	"input": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 {
				return newError("wrong number of arguments. got=%d, want=0 or 1", len(args))
			}

			if len(args) == 1 {
				fmt.Print(args[0].Inspect())
			}

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)

			return &object.String{Value: text}
		},
	},
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return arg
			case *object.String:
				val, err := strconv.ParseInt(arg.Value, 0, 64)
				if err != nil {
					return newError("could not parse %q as integer", arg.Value)
				}
				return &object.Integer{Value: val}
			case *object.Boolean:
				if arg.Value {
					return &object.Integer{Value: 1}
				}
				return &object.Integer{Value: 0}
			default:
				return newError("argument to `int` not supported, got %s", args[0].Type())
			}
		},
	},
	"string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: args[0].Inspect()}
		},
	},
	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"pop": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `pop` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[:length-1])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument 1 to `split` must be STRING, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("argument 2 to `split` must be STRING, got %s", args[1].Type())
			}
			str := args[0].(*object.String).Value
			sep := args[1].(*object.String).Value
			parts := strings.Split(str, sep)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument 1 to `join` must be ARRAY, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("argument 2 to `join` must be STRING, got %s", args[1].Type())
			}
			arr := args[0].(*object.Array)
			sep := args[1].(*object.String).Value
			parts := make([]string, len(arr.Elements))
			for i, e := range arr.Elements {
				if e.Type() != object.STRING_OBJ {
					return newError("array elements must be strings")
				}
				parts[i] = e.(*object.String).Value
			}
			return &object.String{Value: strings.Join(parts, sep)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			switch container := args[0].(type) {
			case *object.Array:
				for _, e := range container.Elements {
					if e.Type() == args[1].Type() {
						if e.Inspect() == args[1].Inspect() {
							return TRUE
						}
					}
				}
				return FALSE
			case *object.String:
				if args[1].Type() != object.STRING_OBJ {
					return newError("argument 2 to `contains` on string must be STRING")
				}
				if strings.Contains(container.Value, args[1].(*object.String).Value) {
					return TRUE
				}
				return FALSE
			default:
				return newError("argument 1 to `contains` must be ARRAY or STRING, got %s", args[0].Type())
			}
		},
	},
	"index": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			switch container := args[0].(type) {
			case *object.Array:
				for i, e := range container.Elements {
					if e.Type() == args[1].Type() && e.Inspect() == args[1].Inspect() {
						return &object.Integer{Value: int64(i)}
					}
				}
				return &object.Integer{Value: -1}
			case *object.String:
				if args[1].Type() != object.STRING_OBJ {
					return newError("argument 2 to `index` on string must be STRING")
				}
				idx := strings.Index(container.Value, args[1].(*object.String).Value)
				return &object.Integer{Value: int64(idx)}
			default:
				return newError("argument 1 to `index` must be ARRAY or STRING, got %s", args[0].Type())
			}
		},
	},
	"upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `upper` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
		},
	},
	"lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `lower` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(args[0].(*object.String).Value)}
		},
	},
	"keys": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.HASH_OBJ {
				return newError("argument to `keys` must be HASH, got %s", args[0].Type())
			}
			hash := args[0].(*object.Hash)
			elements := []object.Object{}
			for _, pair := range hash.Pairs {
				elements = append(elements, pair.Key)
			}
			return &object.Array{Elements: elements}
		},
	},
	"values": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.HASH_OBJ {
				return newError("argument to `values` must be HASH, got %s", args[0].Type())
			}
			hash := args[0].(*object.Hash)
			elements := []object.Object{}
			for _, pair := range hash.Pairs {
				elements = append(elements, pair.Value)
			}
			return &object.Array{Elements: elements}
		},
	},
	"map":    nil, // initialized in init()
	"filter": nil, // initialized in init()
	"reduce": nil, // initialized in init()
}

// isCallable checks if an object can be called as a function
func isCallable(obj object.Object) bool {
	switch obj.Type() {
	case object.FUNCTION_OBJ, object.ARROW_FUNCTION_OBJ, object.BUILTIN_OBJ:
		return true
	default:
		return false
	}
}

// getParamCount returns the number of parameters for a callable
func getParamCount(obj object.Object) int {
	switch fn := obj.(type) {
	case *object.Function:
		return len(fn.Parameters)
	case *object.ArrowFunction:
		return len(fn.Parameters)
	default:
		return 0
	}
}

func init() {
	builtins["map"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument 1 to `map` must be ARRAY, got %s", args[0].Type())
			}
			if !isCallable(args[1]) {
				return newError("argument 2 to `map` must be FUNCTION, got %s", args[1].Type())
			}
			arr := args[0].(*object.Array)
			fn := args[1]
			paramCount := getParamCount(fn)
			elements := make([]object.Object, len(arr.Elements))
			for i, e := range arr.Elements {
				fnArgs := []object.Object{e}
				if paramCount > 1 {
					fnArgs = append(fnArgs, &object.Integer{Value: int64(i)})
				}
				result := applyFunction(fn, fnArgs)
				if isError(result) {
					return result
				}
				elements[i] = result
			}
			return &object.Array{Elements: elements}
		},
	}
	builtins["filter"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument 1 to `filter` must be ARRAY, got %s", args[0].Type())
			}
			if !isCallable(args[1]) {
				return newError("argument 2 to `filter` must be FUNCTION, got %s", args[1].Type())
			}
			arr := args[0].(*object.Array)
			fn := args[1]
			paramCount := getParamCount(fn)
			elements := []object.Object{}
			for i, e := range arr.Elements {
				fnArgs := []object.Object{e}
				if paramCount > 1 {
					fnArgs = append(fnArgs, &object.Integer{Value: int64(i)})
				}
				result := applyFunction(fn, fnArgs)
				if isError(result) {
					return result
				}
				if isTruthy(result) {
					elements = append(elements, e)
				}
			}
			return &object.Array{Elements: elements}
		},
	}
	builtins["reduce"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 || len(args) > 3 {
				return newError("wrong number of arguments. got=%d, want=2 or 3", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument 1 to `reduce` must be ARRAY, got %s", args[0].Type())
			}
			if !isCallable(args[1]) {
				return newError("argument 2 to `reduce` must be FUNCTION, got %s", args[1].Type())
			}
			arr := args[0].(*object.Array)
			fn := args[1]
			paramCount := getParamCount(fn)

			var accumulator object.Object
			startIdx := 0

			if len(args) == 3 {
				accumulator = args[2]
			} else if len(arr.Elements) > 0 {
				accumulator = arr.Elements[0]
				startIdx = 1
			} else {
				return newError("reduce of empty array with no initial value")
			}

			for i := startIdx; i < len(arr.Elements); i++ {
				fnArgs := []object.Object{accumulator, arr.Elements[i]}
				if paramCount > 2 {
					fnArgs = append(fnArgs, &object.Integer{Value: int64(i)})
				}
				result := applyFunction(fn, fnArgs)
				if isError(result) {
					return result
				}
				accumulator = result
			}
			return accumulator
		},
	}
}

func createModule(props map[string]object.Object) *object.Hash {
	pairs := make(map[object.HashKey]object.HashPair)
	for name, val := range props {
		key := &object.String{Value: name}
		hashed := key.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: val}
	}
	return &object.Hash{Pairs: pairs}
}

func createSocketObject(conn net.Conn) *object.Hash {
	connMethods := map[string]object.Object{
		"write": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.STRING_OBJ {
					return newError("argument to `write` must be STRING, got %s", args[0].Type())
				}
				data := args[0].(*object.String).Value
				_, err := conn.Write([]byte(data))
				if err != nil {
					return newError("Write error: %s", err)
				}
				return TRUE
			},
		},
		"read": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.INTEGER_OBJ {
					return newError("argument to `read` must be INTEGER (buffer size), got %s", args[0].Type())
				}
				size := args[0].(*object.Integer).Value
				buf := make([]byte, size)
				n, err := conn.Read(buf)
				if n > 0 {
					return &object.String{Value: string(buf[:n])}
				}
				if err != nil {
					if err.Error() == "EOF" {
						return &object.String{Value: ""}
					}
					return newError("Read error: %s", err)
				}
				return &object.String{Value: string(buf[:n])}
			},
		},
		"close": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				conn.Close()
				return TRUE
			},
		},
		"remoteAddr": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				return &object.String{Value: conn.RemoteAddr().String()}
			},
		},
	}
	return createModule(connMethods)
}

var moduleRegistry = map[string]func() *object.Hash{}

func RegisterBuiltinModules() {
	// OS Module
	moduleRegistry["os"] = func() *object.Hash {
		osMethods := map[string]object.Object{
			"readFile": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `readFile` must be STRING, got %s", args[0].Type())
					}
					filename := args[0].(*object.String).Value
					content, err := ioutil.ReadFile(filename)
					if err != nil {
						return newError("IO error: %s", err)
					}
					return &object.String{Value: string(content)}
				},
			},
			"writeFile": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument 1 to `writeFile` must be STRING, got %s", args[0].Type())
					}
					if args[1].Type() != object.STRING_OBJ {
						return newError("argument 2 to `writeFile` must be STRING, got %s", args[1].Type())
					}
					filename := args[0].(*object.String).Value
					content := args[1].(*object.String).Value
					err := ioutil.WriteFile(filename, []byte(content), 0644)
					if err != nil {
						return newError("IO error: %s", err)
					}
					return &object.Boolean{Value: true}
				},
			},
			"remove": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `remove` must be STRING, got %s", args[0].Type())
					}
					filename := args[0].(*object.String).Value
					err := os.Remove(filename)
					if err != nil {
						return newError("IO error: %s", err)
					}
					return TRUE
				},
			},
			"exists": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `exists` must be STRING, got %s", args[0].Type())
					}
					filename := args[0].(*object.String).Value
					_, err := os.Stat(filename)
					if os.IsNotExist(err) {
						return FALSE
					}
					return TRUE
				},
			},
			"exit": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `exit` must be INTEGER, got %s", args[0].Type())
					}
					code := args[0].(*object.Integer).Value
					os.Exit(int(code))
					return NULL
				},
			},
			"mkdir": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `mkdir` must be STRING, got %s", args[0].Type())
					}
					path := args[0].(*object.String).Value
					err := os.MkdirAll(path, 0755)
					if err != nil {
						return newError("IO error: %s", err)
					}
					return TRUE
				},
			},
			"readDir": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `readDir` must be STRING, got %s", args[0].Type())
					}
					path := args[0].(*object.String).Value
					files, err := ioutil.ReadDir(path)
					if err != nil {
						return newError("IO error: %s", err)
					}
					elements := []object.Object{}
					for _, file := range files {
						elements = append(elements, &object.String{Value: file.Name()})
					}
					return &object.Array{Elements: elements}
				},
			},
			"stat": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `stat` must be STRING, got %s", args[0].Type())
					}
					path := args[0].(*object.String).Value
					info, err := os.Stat(path)
					if err != nil {
						return newError("IO error: %s", err)
					}

					statMap := make(map[object.HashKey]object.HashPair)

					nameKey := &object.String{Value: "name"}
					nameVal := &object.String{Value: info.Name()}
					statMap[nameKey.HashKey()] = object.HashPair{Key: nameKey, Value: nameVal}

					sizeKey := &object.String{Value: "size"}
					sizeVal := &object.Integer{Value: info.Size()}
					statMap[sizeKey.HashKey()] = object.HashPair{Key: sizeKey, Value: sizeVal}

					dirKey := &object.String{Value: "isDir"}
					dirVal := nativeBoolToBooleanObject(info.IsDir())
					statMap[dirKey.HashKey()] = object.HashPair{Key: dirKey, Value: dirVal}

					return &object.Hash{Pairs: statMap}
				},
			},
			"rename": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
						return newError("arguments to `rename` must be STRING")
					}
					oldPath := args[0].(*object.String).Value
					newPath := args[1].(*object.String).Value
					err := os.Rename(oldPath, newPath)
					if err != nil {
						return newError("IO error: %s", err)
					}
					return TRUE
				},
			},
			"getwd": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 0 {
						return newError("wrong number of arguments. got=%d, want=0", len(args))
					}
					dir, err := os.Getwd()
					if err != nil {
						return newError("IO error: %s", err)
					}
					return &object.String{Value: dir}
				},
			},
			"chdir": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `chdir` must be STRING, got %s", args[0].Type())
					}
					path := args[0].(*object.String).Value
					err := os.Chdir(path)
					if err != nil {
						return newError("IO error: %s", err)
					}
					return TRUE
				},
			},
		}

		// Environment variables
		envMap := make(map[object.HashKey]object.HashPair)
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			if len(pair) == 2 {
				key := &object.String{Value: pair[0]}
				val := &object.String{Value: pair[1]}
				envMap[key.HashKey()] = object.HashPair{Key: key, Value: val}
			}
		}
		osMethods["env"] = &object.Hash{Pairs: envMap}

		// Add args array to os module
		argsElements := []object.Object{}
		for _, arg := range os.Args {
			argsElements = append(argsElements, &object.String{Value: arg})
		}
		osMethods["args"] = &object.Array{Elements: argsElements}

		return createModule(osMethods)
	}

	// Net Module
	moduleRegistry["net"] = func() *object.Hash {
		netMethods := map[string]object.Object{
			"get": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `get` must be STRING, got %s", args[0].Type())
					}
					url := args[0].(*object.String).Value
					resp, err := http.Get(url)
					if err != nil {
						return newError("HTTP error: %s", err)
					}
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return newError("HTTP read error: %s", err)
					}
					return &object.String{Value: string(body)}
				},
			},
			"post": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 3 {
						return newError("wrong number of arguments. got=%d, want=3", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument 1 to `post` must be STRING (url), got %s", args[0].Type())
					}
					if args[1].Type() != object.STRING_OBJ {
						return newError("argument 2 to `post` must be STRING (contentType), got %s", args[1].Type())
					}
					if args[2].Type() != object.STRING_OBJ {
						return newError("argument 3 to `post` must be STRING (body), got %s", args[2].Type())
					}

					url := args[0].(*object.String).Value
					contentType := args[1].(*object.String).Value
					body := args[2].(*object.String).Value

					resp, err := http.Post(url, contentType, strings.NewReader(body))
					if err != nil {
						return newError("HTTP error: %s", err)
					}
					defer resp.Body.Close()

					respBody, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return newError("HTTP read error: %s", err)
					}
					return &object.String{Value: string(respBody)}
				},
			},
			"head": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `head` must be STRING, got %s", args[0].Type())
					}
					url := args[0].(*object.String).Value
					resp, err := http.Head(url)
					if err != nil {
						return newError("HTTP error: %s", err)
					}
					defer resp.Body.Close()

					// Return headers as hash? Or just status?
					// Let's return status code for now, or a hash with status and headers.
					// For simplicity, let's return status code as integer.
					return &object.Integer{Value: int64(resp.StatusCode)}
				},
			},
			"delete": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `delete` must be STRING, got %s", args[0].Type())
					}
					url := args[0].(*object.String).Value
					req, err := http.NewRequest("DELETE", url, nil)
					if err != nil {
						return newError("HTTP error: %s", err)
					}
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						return newError("HTTP error: %s", err)
					}
					defer resp.Body.Close()

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return newError("HTTP read error: %s", err)
					}
					return &object.String{Value: string(body)}
				},
			},
			"put": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 3 {
						return newError("wrong number of arguments. got=%d, want=3", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument 1 to `put` must be STRING (url), got %s", args[0].Type())
					}
					if args[1].Type() != object.STRING_OBJ {
						return newError("argument 2 to `put` must be STRING (contentType), got %s", args[1].Type())
					}
					if args[2].Type() != object.STRING_OBJ {
						return newError("argument 3 to `put` must be STRING (body), got %s", args[2].Type())
					}

					url := args[0].(*object.String).Value
					contentType := args[1].(*object.String).Value
					body := args[2].(*object.String).Value

					req, err := http.NewRequest("PUT", url, strings.NewReader(body))
					if err != nil {
						return newError("HTTP error: %s", err)
					}
					req.Header.Set("Content-Type", contentType)

					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						return newError("HTTP error: %s", err)
					}
					defer resp.Body.Close()

					respBody, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return newError("HTTP read error: %s", err)
					}
					return &object.String{Value: string(respBody)}
				},
			},
			"dial": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument 1 to `dial` must be STRING (network), got %s", args[0].Type())
					}
					if args[1].Type() != object.STRING_OBJ {
						return newError("argument 2 to `dial` must be STRING (address), got %s", args[1].Type())
					}

					network := args[0].(*object.String).Value
					address := args[1].(*object.String).Value

					conn, err := net.Dial(network, address)
					if err != nil {
						return newError("Dial error: %s", err)
					}

					return createSocketObject(conn)
				},
			},
			"listenTcp": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument 1 to `listenTcp` must be STRING (address), got %s", args[0].Type())
					}
					// args[1] is handler function
					handlerFn := args[1]
					if handlerFn.Type() != object.FUNCTION_OBJ && handlerFn.Type() != object.BUILTIN_OBJ {
						return newError("argument 2 to `listenTcp` must be FUNCTION, got %s", args[1].Type())
					}

					addr := args[0].(*object.String).Value
					listener, err := net.Listen("tcp", addr)
					if err != nil {
						return newError("Listen error: %s", err)
					}

					fmt.Printf("TCP Server listening on %s\n", addr)

					for {
						conn, err := listener.Accept()
						if err != nil {
							fmt.Printf("Accept error: %s\n", err)
							continue
						}

						// Handle connection in goroutine
						go func(c net.Conn) {
							connObj := createSocketObject(c)
							applyFunction(handlerFn, []object.Object{connObj})
							c.Close()
						}(conn)
					}
				},
			},
			"listen": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument 1 to `listen` must be STRING (address), got %s", args[0].Type())
					}
					// args[1] should be a function
					handlerFn := args[1]
					if handlerFn.Type() != object.FUNCTION_OBJ && handlerFn.Type() != object.BUILTIN_OBJ {
						return newError("argument 2 to `listen` must be FUNCTION, got %s", args[1].Type())
					}

					addr := args[0].(*object.String).Value

					mux := http.NewServeMux()
					mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						// Prepare request object
						reqMap := make(map[object.HashKey]object.HashPair)

						methodKey := &object.String{Value: "method"}
						methodVal := &object.String{Value: r.Method}
						reqMap[methodKey.HashKey()] = object.HashPair{Key: methodKey, Value: methodVal}

						urlKey := &object.String{Value: "url"}
						urlVal := &object.String{Value: r.URL.String()}
						reqMap[urlKey.HashKey()] = object.HashPair{Key: urlKey, Value: urlVal}

						pathKey := &object.String{Value: "path"}
						pathVal := &object.String{Value: r.URL.Path}
						reqMap[pathKey.HashKey()] = object.HashPair{Key: pathKey, Value: pathVal}

						// Read body
						bodyBytes, _ := ioutil.ReadAll(r.Body)
						bodyKey := &object.String{Value: "body"}
						bodyVal := &object.String{Value: string(bodyBytes)}
						reqMap[bodyKey.HashKey()] = object.HashPair{Key: bodyKey, Value: bodyVal}

						reqObj := &object.Hash{Pairs: reqMap}

						// Call handler
						result := applyFunction(handlerFn, []object.Object{reqObj})

						// Process result
						if isError(result) {
							http.Error(w, result.(*object.Error).Message, http.StatusInternalServerError)
							return
						}

						// Default response
						status := http.StatusOK
						body := ""
						contentType := "text/plain"

						if result.Type() == object.HASH_OBJ {
							hash := result.(*object.Hash)

							// Status
							statusKey := &object.String{Value: "status"}
							if pair, ok := hash.Pairs[statusKey.HashKey()]; ok {
								if pair.Value.Type() == object.INTEGER_OBJ {
									status = int(pair.Value.(*object.Integer).Value)
								}
							}

							// Body
							respBodyKey := &object.String{Value: "body"}
							if pair, ok := hash.Pairs[respBodyKey.HashKey()]; ok {
								if pair.Value.Type() == object.STRING_OBJ {
									body = pair.Value.(*object.String).Value
								}
							}

							// Content-Type
							ctKey := &object.String{Value: "content_type"}
							if pair, ok := hash.Pairs[ctKey.HashKey()]; ok {
								if pair.Value.Type() == object.STRING_OBJ {
									contentType = pair.Value.(*object.String).Value
								}
							}
						} else if result.Type() == object.STRING_OBJ {
							body = result.(*object.String).Value
							contentType = "text/html" // Default to HTML if string returned? Or plain? Let's say HTML for convenience.
						}

						w.Header().Set("Content-Type", contentType)
						w.WriteHeader(status)
						w.Write([]byte(body))
					})

					fmt.Printf("Listening on %s...\n", addr)
					err := http.ListenAndServe(addr, mux)
					if err != nil {
						return newError("Server error: %s", err)
					}
					return NULL
				},
			},
		}
		return createModule(netMethods)
	}

	// Std Module
	moduleRegistry["std"] = func() *object.Hash {
		return createModule(map[string]object.Object{
			"version":  &object.String{Value: "1.0.0"},
			"first":    builtins["first"],
			"last":     builtins["last"],
			"rest":     builtins["rest"],
			"push":     builtins["push"],
			"pop":      builtins["pop"],
			"split":    builtins["split"],
			"join":     builtins["join"],
			"contains": builtins["contains"],
			"index":    builtins["index"],
			"upper":    builtins["upper"],
			"lower":    builtins["lower"],
			"keys":     builtins["keys"],
			"values":   builtins["values"],
		})
	}

	// Math Module
	moduleRegistry["math"] = func() *object.Hash {
		mathMethods := map[string]object.Object{
			"pi": &object.Float{Value: math.Pi},
			"abs": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					switch arg := args[0].(type) {
					case *object.Integer:
						if arg.Value < 0 {
							return &object.Integer{Value: -arg.Value}
						}
						return arg
					case *object.Float:
						return &object.Float{Value: math.Abs(arg.Value)}
					default:
						return newError("argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
					}
				},
			},
			"sin": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.FLOAT_OBJ && args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `sin` must be FLOAT or INTEGER")
					}
					val := 0.0
					if args[0].Type() == object.INTEGER_OBJ {
						val = float64(args[0].(*object.Integer).Value)
					} else {
						val = args[0].(*object.Float).Value
					}
					return &object.Float{Value: math.Sin(val)}
				},
			},
			"cos": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.FLOAT_OBJ && args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `cos` must be FLOAT or INTEGER")
					}
					val := 0.0
					if args[0].Type() == object.INTEGER_OBJ {
						val = float64(args[0].(*object.Integer).Value)
					} else {
						val = args[0].(*object.Float).Value
					}
					return &object.Float{Value: math.Cos(val)}
				},
			},
			"sqrt": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.FLOAT_OBJ && args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `sqrt` must be FLOAT or INTEGER")
					}
					val := 0.0
					if args[0].Type() == object.INTEGER_OBJ {
						val = float64(args[0].(*object.Integer).Value)
					} else {
						val = args[0].(*object.Float).Value
					}
					return &object.Float{Value: math.Sqrt(val)}
				},
			},
			"pow": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					x := 0.0
					y := 0.0

					if args[0].Type() == object.INTEGER_OBJ {
						x = float64(args[0].(*object.Integer).Value)
					} else if args[0].Type() == object.FLOAT_OBJ {
						x = args[0].(*object.Float).Value
					} else {
						return newError("argument 1 to `pow` must be FLOAT or INTEGER")
					}

					if args[1].Type() == object.INTEGER_OBJ {
						y = float64(args[1].(*object.Integer).Value)
					} else if args[1].Type() == object.FLOAT_OBJ {
						y = args[1].(*object.Float).Value
					} else {
						return newError("argument 2 to `pow` must be FLOAT or INTEGER")
					}

					return &object.Float{Value: math.Pow(x, y)}
				},
			},
			"floor": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					switch arg := args[0].(type) {
					case *object.Integer:
						return arg
					case *object.Float:
						return &object.Integer{Value: int64(math.Floor(arg.Value))}
					default:
						return newError("argument to `floor` must be INTEGER or FLOAT, got %s", args[0].Type())
					}
				},
			},
			"ceil": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					switch arg := args[0].(type) {
					case *object.Integer:
						return arg
					case *object.Float:
						return &object.Integer{Value: int64(math.Ceil(arg.Value))}
					default:
						return newError("argument to `ceil` must be INTEGER or FLOAT, got %s", args[0].Type())
					}
				},
			},
			"round": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					switch arg := args[0].(type) {
					case *object.Integer:
						return arg
					case *object.Float:
						return &object.Integer{Value: int64(math.Round(arg.Value))}
					default:
						return newError("argument to `round` must be INTEGER or FLOAT, got %s", args[0].Type())
					}
				},
			},
			"min": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 2 {
						return newError("wrong number of arguments. got=%d, want=at least 2", len(args))
					}
					// Convert all to float for comparison
					minVal := math.MaxFloat64
					isAllIntegers := true
					for _, arg := range args {
						var val float64
						switch a := arg.(type) {
						case *object.Integer:
							val = float64(a.Value)
						case *object.Float:
							val = a.Value
							isAllIntegers = false
						default:
							return newError("arguments to `min` must be INTEGER or FLOAT, got %s", arg.Type())
						}
						if val < minVal {
							minVal = val
						}
					}
					if isAllIntegers {
						return &object.Integer{Value: int64(minVal)}
					}
					return &object.Float{Value: minVal}
				},
			},
			"max": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 2 {
						return newError("wrong number of arguments. got=%d, want=at least 2", len(args))
					}
					// Convert all to float for comparison
					maxVal := -math.MaxFloat64
					isAllIntegers := true
					for _, arg := range args {
						var val float64
						switch a := arg.(type) {
						case *object.Integer:
							val = float64(a.Value)
						case *object.Float:
							val = a.Value
							isAllIntegers = false
						default:
							return newError("arguments to `max` must be INTEGER or FLOAT, got %s", arg.Type())
						}
						if val > maxVal {
							maxVal = val
						}
					}
					if isAllIntegers {
						return &object.Integer{Value: int64(maxVal)}
					}
					return &object.Float{Value: maxVal}
				},
			},
			"random": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) == 0 {
						// Return random float between 0 and 1
						return &object.Float{Value: rand.Float64()}
					} else if len(args) == 1 {
						// Return random int between 0 and n-1
						if args[0].Type() != object.INTEGER_OBJ {
							return newError("argument to `random` must be INTEGER, got %s", args[0].Type())
						}
						n := args[0].(*object.Integer).Value
						if n <= 0 {
							return newError("argument to `random` must be positive")
						}
						return &object.Integer{Value: rand.Int63n(n)}
					} else if len(args) == 2 {
						// Return random int between min and max (inclusive)
						if args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
							return newError("arguments to `random` must be INTEGER")
						}
						min := args[0].(*object.Integer).Value
						max := args[1].(*object.Integer).Value
						if max < min {
							return newError("max must be >= min in random(min, max)")
						}
						return &object.Integer{Value: min + rand.Int63n(max-min+1)}
					}
					return newError("wrong number of arguments. got=%d, want=0, 1, or 2", len(args))
				},
			},
			"e": &object.Float{Value: math.E},
			"tan": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					val := getNumericValue(args[0])
					if val == nil {
						return newError("argument to `tan` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Tan(*val)}
				},
			},
			"log": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					val := getNumericValue(args[0])
					if val == nil {
						return newError("argument to `log` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Log(*val)}
				},
			},
			"log10": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					val := getNumericValue(args[0])
					if val == nil {
						return newError("argument to `log10` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Log10(*val)}
				},
			},
		}
		return createModule(mathMethods)
	}

	// JSON Module
	moduleRegistry["json"] = func() *object.Hash {
		jsonMethods := map[string]object.Object{
			"parse": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `json.parse` must be STRING, got %s", args[0].Type())
					}
					jsonStr := args[0].(*object.String).Value
					return parseJSON(jsonStr)
				},
			},
			"stringify": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 1 || len(args) > 2 {
						return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
					}
					indent := ""
					if len(args) == 2 {
						if args[1].Type() == object.INTEGER_OBJ {
							spaces := args[1].(*object.Integer).Value
							for i := int64(0); i < spaces; i++ {
								indent += " "
							}
						} else if args[1].Type() == object.STRING_OBJ {
							indent = args[1].(*object.String).Value
						}
					}
					return stringifyJSON(args[0], indent)
				},
			},
			"valid": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `json.valid` must be STRING, got %s", args[0].Type())
					}
					jsonStr := args[0].(*object.String).Value
					var js interface{}
					if err := json.Unmarshal([]byte(jsonStr), &js); err != nil {
						return FALSE
					}
					return TRUE
				},
			},
		}
		return createModule(jsonMethods)
	}

	// Time Module
	moduleRegistry["time"] = func() *object.Hash {
		timeMethods := map[string]object.Object{
			"now": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					return &object.Integer{Value: time.Now().Unix()}
				},
			},
			"nowMs": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
				},
			},
			"format": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 1 || len(args) > 2 {
						return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
					}
					var t time.Time
					if args[0].Type() == object.INTEGER_OBJ {
						t = time.Unix(args[0].(*object.Integer).Value, 0)
					} else {
						return newError("first argument to `time.format` must be INTEGER (unix timestamp)")
					}
					layout := "2006-01-02 15:04:05"
					if len(args) == 2 {
						if args[1].Type() != object.STRING_OBJ {
							return newError("second argument to `time.format` must be STRING (format)")
						}
						layout = convertTimeFormat(args[1].(*object.String).Value)
					}
					return &object.String{Value: t.Format(layout)}
				},
			},
			"parse": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 1 || len(args) > 2 {
						return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("first argument to `time.parse` must be STRING")
					}
					dateStr := args[0].(*object.String).Value
					layout := "2006-01-02 15:04:05"
					if len(args) == 2 {
						if args[1].Type() != object.STRING_OBJ {
							return newError("second argument to `time.parse` must be STRING (format)")
						}
						layout = convertTimeFormat(args[1].(*object.String).Value)
					}
					t, err := time.Parse(layout, dateStr)
					if err != nil {
						return newError("failed to parse time: %s", err.Error())
					}
					return &object.Integer{Value: t.Unix()}
				},
			},
			"year": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.year` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Year())}
				},
			},
			"month": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.month` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Month())}
				},
			},
			"day": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.day` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Day())}
				},
			},
			"hour": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.hour` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Hour())}
				},
			},
			"minute": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.minute` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Minute())}
				},
			},
			"second": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.second` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Second())}
				},
			},
			"weekday": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.weekday` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Weekday())}
				},
			},
			"sleep": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.sleep` must be INTEGER (milliseconds)")
					}
					ms := args[0].(*object.Integer).Value
					time.Sleep(time.Duration(ms) * time.Millisecond)
					return NULL
				},
			},
			"date": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					// time.date() returns current date as "YYYY-MM-DD"
					// time.date(timestamp) returns date of timestamp
					var t time.Time
					if len(args) == 0 {
						t = time.Now()
					} else if len(args) == 1 && args[0].Type() == object.INTEGER_OBJ {
						t = time.Unix(args[0].(*object.Integer).Value, 0)
					} else {
						return newError("argument to `time.date` must be INTEGER (unix timestamp) or no arguments")
					}
					return &object.String{Value: t.Format("2006-01-02")}
				},
			},
			"time": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					// time.time() returns current time as "HH:MM:SS"
					// time.time(timestamp) returns time of timestamp
					var t time.Time
					if len(args) == 0 {
						t = time.Now()
					} else if len(args) == 1 && args[0].Type() == object.INTEGER_OBJ {
						t = time.Unix(args[0].(*object.Integer).Value, 0)
					} else {
						return newError("argument to `time.time` must be INTEGER (unix timestamp) or no arguments")
					}
					return &object.String{Value: t.Format("15:04:05")}
				},
			},
		}
		return createModule(timeMethods)
	}
}

// Helper function to get numeric value as float64 pointer
func getNumericValue(obj object.Object) *float64 {
	switch o := obj.(type) {
	case *object.Integer:
		val := float64(o.Value)
		return &val
	case *object.Float:
		return &o.Value
	default:
		return nil
	}
}

// parseJSON converts a JSON string to Victoria objects
func parseJSON(jsonStr string) object.Object {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return newError("failed to parse JSON: %s", err.Error())
	}
	return jsonToObject(data)
}

// jsonToObject converts a Go interface{} to Victoria object
func jsonToObject(data interface{}) object.Object {
	if data == nil {
		return NULL
	}
	switch v := data.(type) {
	case bool:
		if v {
			return TRUE
		}
		return FALSE
	case float64:
		// Check if it's actually an integer
		if v == float64(int64(v)) {
			return &object.Integer{Value: int64(v)}
		}
		return &object.Float{Value: v}
	case string:
		return &object.String{Value: v}
	case []interface{}:
		elements := make([]object.Object, len(v))
		for i, elem := range v {
			elements[i] = jsonToObject(elem)
		}
		return &object.Array{Elements: elements}
	case map[string]interface{}:
		pairs := make(map[object.HashKey]object.HashPair)
		for key, val := range v {
			keyObj := &object.String{Value: key}
			valObj := jsonToObject(val)
			pairs[keyObj.HashKey()] = object.HashPair{Key: keyObj, Value: valObj}
		}
		return &object.Hash{Pairs: pairs}
	default:
		return newError("unsupported JSON type: %T", v)
	}
}

// stringifyJSON converts a Victoria object to JSON string
func stringifyJSON(obj object.Object, indent string) object.Object {
	goVal := objectToGo(obj)
	var jsonBytes []byte
	var err error
	if indent != "" {
		jsonBytes, err = json.MarshalIndent(goVal, "", indent)
	} else {
		jsonBytes, err = json.Marshal(goVal)
	}
	if err != nil {
		return newError("failed to stringify JSON: %s", err.Error())
	}
	return &object.String{Value: string(jsonBytes)}
}

// objectToGo converts a Victoria object to Go interface{}
func objectToGo(obj object.Object) interface{} {
	switch o := obj.(type) {
	case *object.Integer:
		return o.Value
	case *object.Float:
		return o.Value
	case *object.Boolean:
		return o.Value
	case *object.String:
		return o.Value
	case *object.Null:
		return nil
	case *object.Array:
		result := make([]interface{}, len(o.Elements))
		for i, elem := range o.Elements {
			result[i] = objectToGo(elem)
		}
		return result
	case *object.Hash:
		result := make(map[string]interface{})
		for _, pair := range o.Pairs {
			key := pair.Key.Inspect()
			result[key] = objectToGo(pair.Value)
		}
		return result
	default:
		return obj.Inspect()
	}
}

// convertTimeFormat converts common format tokens to Go's time format
func convertTimeFormat(format string) string {
	// Common format tokens: YYYY, MM, DD, HH, mm, ss, etc.
	replacements := map[string]string{
		"YYYY": "2006",
		"YY":   "06",
		"MM":   "01",
		"DD":   "02",
		"HH":   "15",
		"hh":   "03",
		"mm":   "04",
		"ss":   "05",
		"SSS":  "000",
		"Z":    "-0700",
		"A":    "PM",
		"a":    "pm",
	}
	result := format
	for token, goFormat := range replacements {
		result = strings.Replace(result, token, goFormat, -1)
	}
	return result
}

func evalIncludeStatement(node *ast.IncludeStatement, env *object.Environment) object.Object {
	for _, moduleName := range node.Modules {
		if factory, ok := moduleRegistry[moduleName]; ok {
			module := factory()
			env.Set(moduleName, module)
		} else {
			// Try to load as file
			filename := moduleName
			if !strings.HasSuffix(filename, ".vc") {
				// Check if file exists as is, if not try adding .vc
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					filename = filename + ".vc"
				}
			}

			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return newError("module or file not found: %s", moduleName)
			}

			l := lexer.New(string(content))
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				msg := fmt.Sprintf("parser errors in %s:\n", filename)
				for _, msgErr := range p.Errors() {
					msg += "\t" + msgErr + "\n"
				}
				return newError(msg)
			}

			// Evaluate the program in the CURRENT environment
			result := Eval(program, env)
			if isError(result) {
				return result
			}
		}
	}
	return NULL
}
