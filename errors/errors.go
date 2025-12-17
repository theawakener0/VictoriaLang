package errors

import (
	"fmt"
	"math/rand"
	"strings"
)

// Programming jokes to lighten the mood when errors occur
var programmerJokes = []string{
	"Why do programmers prefer dark mode? Because light attracts bugs!",
	"There are only 2 types of people: those who understand binary and those who don't.",
	"A SQL query walks into a bar, walks up to two tables and asks... 'Can I join you?'",
	"Why do Java developers wear glasses? Because they don't C#!",
	"How many programmers does it take to change a light bulb? None, that's a hardware problem.",
	"Why was the JavaScript developer sad? Because he didn't Node how to Express himself.",
	"A programmer's wife tells him: 'Go to the store and buy a loaf of bread. If they have eggs, buy a dozen.' He comes home with 12 loaves of bread.",
	"There's no place like 127.0.0.1",
	"Why do programmers always mix up Halloween and Christmas? Because Oct 31 == Dec 25!",
	"Programming is like writing a book... except if you miss a single comma on page 126, the whole thing makes no sense.",
	"99 little bugs in the code, 99 little bugs. Take one down, patch it around... 127 little bugs in the code.",
	"It works on my machine! ¯\\_(ツ)_/¯",
	"The best thing about a boolean is that even if you're wrong, you're only off by a bit.",
	"A programmer puts two glasses on his bedside table before going to sleep. A full one, in case he gets thirsty, and an empty one, in case he doesn't.",
	"To understand recursion, you must first understand recursion.",
	"I would tell you a UDP joke, but you might not get it.",
	"Why did the developer go broke? Because he used up all his cache!",
	"['hip', 'hip'] // hooray!",
	"A foo walks into a bar, takes a look around and says 'Hello World!'",
	"An SEO expert walks into a bar, bars, pub, tavern, public house, Irish pub, drinks, beer...",
	"The glass is neither half full nor half empty. It's twice as big as it needs to be.",
	"I've got a really good UDP joke to tell you but I don't know if you'll get it.",
	"If at first you don't succeed, call it version 1.0",
	"Software and cathedrals are much the same – first we build them, then we pray.",
	"Debugging: Being the detective in a crime movie where you are also the murderer.",
	"I don't always test my code, but when I do, I do it in production.",
	"In theory, there's no difference between theory and practice. In practice, there is.",
	"Real programmers count from 0.",
	"!false - It's funny because it's true.",
	// Type system jokes
	"A string walks into a bar. The bartender says, 'We don't serve your type here.'",
	"Why did the type checker break up with the dynamic language? Too many unexpected surprises!",
	"Strong typing: Because 'undefined is not a function' should never be a runtime error.",
	"Types are like vegetables – you know they're good for you, but sometimes you just want dessert.",
	"In a statically typed world, bugs are caught at compile time. In a dynamically typed world, bugs are caught in production.",
	"Type inference: because sometimes the compiler knows you better than you know yourself.",
	"Any: the type that says 'I give up, do whatever you want.'",
	"void: for when your function has commitment issues about returning values.",
}

// getRandomJoke returns a random programming joke (30% chance)
func getRandomJoke() string {
	if rand.Intn(100) < 30 {
		return programmerJokes[rand.Intn(len(programmerJokes))]
	}
	return ""
}

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"

	// Primary colors
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright colors
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
)

// SourceLocation represents a position in source code
type SourceLocation struct {
	Line      int
	Column    int
	EndLine   int
	EndColumn int
	Filename  string
}

// String returns a formatted location string
func (loc SourceLocation) String() string {
	if loc.Filename != "" {
		return fmt.Sprintf("%s:%d:%d", loc.Filename, loc.Line, loc.Column)
	}
	return fmt.Sprintf("%d:%d", loc.Line, loc.Column)
}

// ErrorKind represents the severity of an error
type ErrorKind int

const (
	KindError ErrorKind = iota
	KindWarning
	KindNote
	KindHelp
)

func (k ErrorKind) String() string {
	switch k {
	case KindError:
		return "error"
	case KindWarning:
		return "warning"
	case KindNote:
		return "note"
	case KindHelp:
		return "help"
	default:
		return "unknown"
	}
}

func (k ErrorKind) Color() string {
	switch k {
	case KindError:
		return BrightRed
	case KindWarning:
		return BrightYellow
	case KindNote:
		return BrightCyan
	case KindHelp:
		return BrightGreen
	default:
		return White
	}
}

// Label represents a labeled span in the source code
type Label struct {
	Location SourceLocation
	Message  string
	Primary  bool // If true, this is the primary label (shown in red)
}

// VictoriaError represents a rich error with context
type VictoriaError struct {
	Kind       ErrorKind
	Code       string // Error code like E0001
	Message    string
	Labels     []Label
	Notes      []string
	Help       string
	SourceCode string // The full source code for snippet extraction
}

// NewError creates a new error
func NewError(message string) *VictoriaError {
	return &VictoriaError{
		Kind:    KindError,
		Message: message,
		Labels:  []Label{},
		Notes:   []string{},
	}
}

// NewParseError creates a parse error with location
func NewParseError(message string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: message, Primary: true},
		},
	}
}

// NewRuntimeError creates a runtime error
func NewRuntimeError(message string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: message, Primary: true},
		},
	}
}

// WithCode adds an error code
func (e *VictoriaError) WithCode(code string) *VictoriaError {
	e.Code = code
	return e
}

// WithLabel adds a label to the error
func (e *VictoriaError) WithLabel(loc SourceLocation, message string, primary bool) *VictoriaError {
	e.Labels = append(e.Labels, Label{
		Location: loc,
		Message:  message,
		Primary:  primary,
	})
	return e
}

// WithNote adds a note to the error
func (e *VictoriaError) WithNote(note string) *VictoriaError {
	e.Notes = append(e.Notes, note)
	return e
}

// WithHelp adds help text to the error
func (e *VictoriaError) WithHelp(help string) *VictoriaError {
	e.Help = help
	return e
}

// WithSource sets the source code for snippet extraction
func (e *VictoriaError) WithSource(source string) *VictoriaError {
	e.SourceCode = source
	return e
}

// getSourceLines extracts multiple lines from source code
func getSourceLines(source string, startLine, endLine int) []string {
	lines := strings.Split(source, "\n")
	if startLine < 1 {
		startLine = 1
	}
	if endLine > len(lines) {
		endLine = len(lines)
	}
	if startLine > endLine || startLine > len(lines) {
		return []string{}
	}
	return lines[startLine-1 : endLine]
}

// countDigits returns the number of digits in a number
func countDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n > 0 {
		n /= 10
		count++
	}
	return count
}

// Format returns the formatted error message with colors and source snippets
func (e *VictoriaError) Format() string {
	var sb strings.Builder

	// Header: error[E0001]: message
	headerColor := e.Kind.Color()
	sb.WriteString(fmt.Sprintf("%s%s%s%s", Bold, headerColor, e.Kind.String(), Reset))
	if e.Code != "" {
		sb.WriteString(fmt.Sprintf("%s[%s]%s", Dim, e.Code, Reset))
	}
	sb.WriteString(fmt.Sprintf("%s: %s%s\n", Bold+White, e.Message, Reset))

	// Find the maximum line number for padding
	maxLine := 0
	for _, label := range e.Labels {
		if label.Location.Line > maxLine {
			maxLine = label.Location.Line
		}
		if label.Location.EndLine > maxLine {
			maxLine = label.Location.EndLine
		}
	}
	lineNumWidth := countDigits(maxLine)
	if lineNumWidth < 1 {
		lineNumWidth = 1
	}

	// Source location header
	if len(e.Labels) > 0 && e.Labels[0].Location.Line > 0 {
		loc := e.Labels[0].Location
		padding := strings.Repeat(" ", lineNumWidth)
		sb.WriteString(fmt.Sprintf("%s%s--> %s%s\n", Cyan, padding, Reset, loc.String()))
		sb.WriteString(fmt.Sprintf("%s%s |%s\n", Cyan, padding, Reset))
	}

	// Group labels by line
	labelsByLine := make(map[int][]Label)
	for _, label := range e.Labels {
		labelsByLine[label.Location.Line] = append(labelsByLine[label.Location.Line], label)
	}

	// Print source snippets with annotations
	if e.SourceCode != "" && len(e.Labels) > 0 {
		// Determine context range (show 1-2 lines before and after)
		minLine := e.Labels[0].Location.Line
		maxLine := e.Labels[0].Location.Line
		for _, label := range e.Labels {
			if label.Location.Line < minLine {
				minLine = label.Location.Line
			}
			endL := label.Location.Line
			if label.Location.EndLine > endL {
				endL = label.Location.EndLine
			}
			if endL > maxLine {
				maxLine = endL
			}
		}

		// Context lines
		contextBefore := 1
		contextAfter := 1
		startLine := minLine - contextBefore
		if startLine < 1 {
			startLine = 1
		}
		endLine := maxLine + contextAfter

		lines := getSourceLines(e.SourceCode, startLine, endLine)
		for i, line := range lines {
			currentLineNum := startLine + i
			lineNumStr := fmt.Sprintf("%*d", lineNumWidth, currentLineNum)

			// Check if this line has labels
			labels, hasLabels := labelsByLine[currentLineNum]

			// Print the source line
			sb.WriteString(fmt.Sprintf("%s%s | %s%s\n", Cyan, lineNumStr, Reset, line))

			// Print annotations for this line
			if hasLabels {
				for _, label := range labels {
					padding := strings.Repeat(" ", lineNumWidth)
					sb.WriteString(fmt.Sprintf("%s%s | %s", Cyan, padding, Reset))

					// Calculate the caret position
					col := label.Location.Column
					if col < 1 {
						col = 1
					}
					endCol := label.Location.EndColumn
					if endCol < col {
						endCol = col + 1
					}

					// Spaces before the caret
					spaces := strings.Repeat(" ", col-1)

					// The caret/underline
					underlineLen := endCol - col
					if underlineLen < 1 {
						underlineLen = 1
					}

					var underlineColor string
					var caretChar string
					if label.Primary {
						underlineColor = BrightRed
						caretChar = "^"
					} else {
						underlineColor = BrightBlue
						caretChar = "-"
					}
					underline := strings.Repeat(caretChar, underlineLen)

					sb.WriteString(fmt.Sprintf("%s%s%s%s%s", spaces, Bold, underlineColor, underline, Reset))

					// Label message
					if label.Message != "" {
						sb.WriteString(fmt.Sprintf(" %s%s%s", underlineColor, label.Message, Reset))
					}
					sb.WriteString("\n")
				}
			}
		}

		// Closing bar
		padding := strings.Repeat(" ", lineNumWidth)
		sb.WriteString(fmt.Sprintf("%s%s |%s\n", Cyan, padding, Reset))
	}

	// Notes
	for _, note := range e.Notes {
		padding := strings.Repeat(" ", lineNumWidth)
		sb.WriteString(fmt.Sprintf("%s%s = %s%snote%s: %s\n", Cyan, padding, Reset, Bold+BrightCyan, Reset, note))
	}

	// Help
	if e.Help != "" {
		padding := strings.Repeat(" ", lineNumWidth)
		sb.WriteString(fmt.Sprintf("%s%s = %s%shelp%s: %s\n", Cyan, padding, Reset, Bold+BrightGreen, Reset, e.Help))
	}

	// Random joke (30% chance to lighten the mood)
	if joke := getRandomJoke(); joke != "" {
		lineNumWidth := 1
		if len(e.Labels) > 0 {
			maxLine := e.Labels[0].Location.Line
			for _, label := range e.Labels {
				if label.Location.Line > maxLine {
					maxLine = label.Location.Line
				}
			}
			lineNumWidth = countDigits(maxLine)
		}
		padding := strings.Repeat(" ", lineNumWidth)
		sb.WriteString(fmt.Sprintf("%s%s = %s%sjoke%s: %s\n", Cyan, padding, Reset, Bold+BrightMagenta, Reset, joke))
	}

	return sb.String()
}

// FormatPlain returns the error message without colors (for logging)
func (e *VictoriaError) FormatPlain() string {
	var sb strings.Builder

	// Header
	sb.WriteString(e.Kind.String())
	if e.Code != "" {
		sb.WriteString(fmt.Sprintf("[%s]", e.Code))
	}
	sb.WriteString(fmt.Sprintf(": %s\n", e.Message))

	// Location
	if len(e.Labels) > 0 && e.Labels[0].Location.Line > 0 {
		loc := e.Labels[0].Location
		sb.WriteString(fmt.Sprintf("  --> %s\n", loc.String()))
	}

	// Notes
	for _, note := range e.Notes {
		sb.WriteString(fmt.Sprintf("  = note: %s\n", note))
	}

	// Help
	if e.Help != "" {
		sb.WriteString(fmt.Sprintf("  = help: %s\n", e.Help))
	}

	return sb.String()
}

// Error implements the error interface
func (e *VictoriaError) Error() string {
	return e.Format()
}

// ErrorReporter collects and formats multiple errors
type ErrorReporter struct {
	Errors     []*VictoriaError
	SourceCode string
	Filename   string
}

// NewErrorReporter creates a new error reporter
func NewErrorReporter(source string, filename string) *ErrorReporter {
	return &ErrorReporter{
		Errors:     []*VictoriaError{},
		SourceCode: source,
		Filename:   filename,
	}
}

// AddError adds an error to the reporter
func (r *ErrorReporter) AddError(err *VictoriaError) {
	err.SourceCode = r.SourceCode
	for i := range err.Labels {
		if err.Labels[i].Location.Filename == "" {
			err.Labels[i].Location.Filename = r.Filename
		}
	}
	r.Errors = append(r.Errors, err)
}

// HasErrors returns true if there are any errors
func (r *ErrorReporter) HasErrors() bool {
	return len(r.Errors) > 0
}

// Format returns all errors formatted
func (r *ErrorReporter) Format() string {
	var sb strings.Builder
	for i, err := range r.Errors {
		sb.WriteString(err.Format())
		if i < len(r.Errors)-1 {
			sb.WriteString("\n")
		}
	}
	if len(r.Errors) > 0 {
		sb.WriteString(fmt.Sprintf("\n%s%serror%s: could not compile due to %d previous error(s)\n",
			Bold, BrightRed, Reset, len(r.Errors)))
	}
	return sb.String()
}

// ═══════════════════════════════════════════════════════════════════════════════
// Common Error Constructors for Victoria Language
// ═══════════════════════════════════════════════════════════════════════════════

// TypeMismatchError creates a type mismatch error with contextual help
func TypeMismatchError(left, operator, right string, loc SourceLocation, source string) *VictoriaError {
	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0001",
		Message:    fmt.Sprintf("type mismatch: cannot apply '%s' to %s and %s", operator, left, right),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("'%s' cannot be applied to these types", operator), Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("left operand has type %s", left),
			fmt.Sprintf("right operand has type %s", right),
			"Victoria is dynamically typed but operators require compatible types",
		},
	}

	// Provide specific help based on the type combination
	switch {
	case (left == "STRING" && right == "INTEGER") || (left == "INTEGER" && right == "STRING"):
		err.Help = "use string() to convert integers: \"text\" + string(42)"
		err.Notes = append(err.Notes, "string concatenation requires both operands to be strings")
	case (left == "STRING" && right == "FLOAT") || (left == "FLOAT" && right == "STRING"):
		err.Help = "use string() to convert floats: \"value: \" + string(3.14)"
	case left == "STRING" && right == "STRING" && operator != "+":
		err.Help = "strings only support '+' for concatenation and '==' / '!=' for comparison"
	case (left == "BOOLEAN" && right == "INTEGER") || (left == "INTEGER" && right == "BOOLEAN"):
		err.Help = "use int() to convert boolean: int(true) returns 1, int(false) returns 0"
	case left == "ARRAY" || right == "ARRAY":
		err.Help = "arrays only support '+' for concatenation: [1, 2] + [3, 4]"
		err.Notes = append(err.Notes, "use push(), pop(), or spread operator for array manipulation")
	case left == "HASH" || right == "HASH":
		err.Help = "hashes don't support arithmetic; access values with hash[\"key\"]"
	case left == "NULL" || right == "NULL":
		err.Help = "check for null before performing operations: if value != null { ... }"
		err.Notes = append(err.Notes, "null cannot be used in arithmetic operations")
	default:
		err.Help = "convert one operand to match the other's type"
	}

	return err
}

// UndefinedVariableError creates an undefined variable error with suggestions
func UndefinedVariableError(name string, loc SourceLocation, source string) *VictoriaError {
	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0002",
		Message:    fmt.Sprintf("undefined variable: '%s'", name),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "not found in this scope", Primary: true},
		},
		Notes: []string{
			"variables must be declared before use",
		},
	}

	// Common misspellings and alternatives
	suggestions := map[string]string{
		// Output functions
		"println": "did you mean 'print'? Victoria uses print() for output",
		"printf":  "did you mean 'format'? Victoria uses format() for formatted strings",
		"console": "did you mean 'print'? Victoria uses print() for output",
		"log":     "did you mean 'print'? Victoria uses print() for output",
		"echo":    "did you mean 'print'? Victoria uses print() for output",
		"puts":    "did you mean 'print'? Victoria uses print() for output",
		"write":   "did you mean 'print'? Victoria uses print() for output",
		// Type conversions
		"str":      "did you mean 'string'? Victoria uses string() for conversion",
		"toString": "did you mean 'string'? Victoria uses string() for conversion",
		"String":   "did you mean 'string'? Victoria uses lowercase string()",
		"parseInt": "did you mean 'int'? Victoria uses int() for conversion",
		"toInt":    "did you mean 'int'? Victoria uses int() for conversion",
		"Int":      "did you mean 'int'? Victoria uses lowercase int()",
		"float":    "did you mean an integer? Victoria supports float literals like 3.14",
		// Length functions
		"size":   "did you mean 'len'? Victoria uses len() for length",
		"length": "did you mean 'len'? Victoria uses len() for length",
		"count":  "did you mean 'len'? Victoria uses len() for length",
		"sizeof": "did you mean 'len'? Victoria uses len() for length",
		// Array functions
		"append":  "did you mean 'push'? Victoria uses push(array, element)",
		"add":     "did you mean 'push'? Victoria uses push(array, element)",
		"insert":  "did you mean 'push'? Victoria uses push() for the end; use slicing for other positions",
		"shift":   "did you mean 'rest'? Victoria uses rest(array) to skip the first element",
		"unshift": "use array concatenation: [newElement, ...array]",
		"remove":  "did you mean 'pop'? Victoria uses pop(array) for the last element",
		"delete":  "use filter() to create a new array without elements",
		"concat":  "use the + operator or spread: [...arr1, ...arr2]",
		// Null/nil/none
		"nil":       "did you mean 'null'? Victoria uses null for no value",
		"none":      "did you mean 'null'? Victoria uses null for no value",
		"None":      "did you mean 'null'? Victoria uses null for no value",
		"undefined": "did you mean 'null'? Victoria uses null for no value",
		"NULL":      "did you mean 'null'? Keywords are lowercase in Victoria",
		"Null":      "did you mean 'null'? Keywords are lowercase in Victoria",
		// Booleans
		"True":  "did you mean 'true'? Booleans are lowercase in Victoria",
		"False": "did you mean 'false'? Booleans are lowercase in Victoria",
		"TRUE":  "did you mean 'true'? Booleans are lowercase in Victoria",
		"FALSE": "did you mean 'false'? Booleans are lowercase in Victoria",
		// Function definition
		"fn":       "did you mean 'define'? Victoria uses 'define' for functions",
		"func":     "did you mean 'define'? Victoria uses 'define' for functions",
		"function": "did you mean 'define'? Victoria uses 'define' for functions",
		"lambda":   "did you mean 'define'? Victoria uses 'define' for functions",
		"def":      "did you mean 'define'? Victoria uses 'define' for functions",
		// Variable declaration
		"var": "did you mean 'let'? Victoria uses 'let' for variable declaration",
		// String methods
		"substr":     "use string slicing: str[start:end]",
		"substring":  "use string slicing: str[start:end]",
		"charAt":     "use string indexing: str[index]",
		"indexOf":    "did you mean 'index'? Victoria uses index(string, substring)",
		"includes":   "did you mean 'contains'? Victoria uses contains(string, substring)",
		"startsWith": "use slicing: str[0:len(prefix)] == prefix",
		"endsWith":   "use slicing: str[len(str)-len(suffix):] == suffix",
		"trim":       "Victoria doesn't have trim() yet; use a custom function",
		"replace":    "Victoria doesn't have replace() yet; use split() and join()",
		// Iteration
		"forEach": "use a for-in loop: for item in array { ... }",
		"foreach": "use a for-in loop: for item in array { ... }",
		"each":    "use a for-in loop: for item in array { ... }",
		// Math
		"abs":    "Victoria doesn't have abs() yet; use: if x < 0 { -x } else { x }",
		"max":    "Victoria doesn't have max() yet; use: if a > b { a } else { b }",
		"min":    "Victoria doesn't have min() yet; use: if a < b { a } else { b }",
		"floor":  "use int() to truncate: int(3.7) returns 3",
		"round":  "use int() with 0.5: int(x + 0.5)",
		"sqrt":   "Victoria doesn't have sqrt() yet",
		"pow":    "Victoria doesn't have pow() yet; use repeated multiplication",
		"random": "Victoria doesn't have random() yet",
		// System
		"exit":    "Victoria doesn't have exit() yet",
		"sleep":   "Victoria doesn't have sleep() yet",
		"require": "did you mean 'include'? Victoria uses include \"filename\"",
		"import":  "did you mean 'include'? Victoria uses include \"filename\"",
	}

	if suggestion, ok := suggestions[name]; ok {
		err.Help = suggestion
	} else {
		err.Help = fmt.Sprintf("declare with 'let %s = <value>' or 'const %s = <value>'", name, name)
		err.Notes = append(err.Notes, "check spelling and ensure the variable is in scope")
	}

	return err
}

// UnknownOperatorError creates an unknown operator error with contextual help
func UnknownOperatorError(operator, typeName string, loc SourceLocation, source string) *VictoriaError {
	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0003",
		Message:    fmt.Sprintf("unknown operator: '%s' for type %s", operator, typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "unsupported operator for this type", Primary: true},
		},
	}

	// Type-specific operator help
	switch typeName {
	case "STRING":
		err.Notes = []string{
			"strings support: + (concatenation), == and != (comparison)",
		}
		err.Help = "use string functions for other operations: upper(), lower(), split(), contains()"
	case "BOOLEAN":
		err.Notes = []string{
			"booleans support: == and != (comparison), and, or, ! (logical)",
		}
		err.Help = "use 'and', 'or', '!' for boolean logic, not arithmetic operators"
	case "ARRAY":
		err.Notes = []string{
			"arrays support: + (concatenation), == and != (comparison)",
		}
		err.Help = "use array functions: push(), pop(), map(), filter(), reduce()"
	case "HASH":
		err.Notes = []string{
			"hashes support: == and != (comparison) only",
		}
		err.Help = "access hash values with hash[\"key\"] or hash.key"
	case "FUNCTION":
		err.Notes = []string{
			"functions can only be compared with == and !=",
		}
		err.Help = "call the function first to operate on its return value: fn() + 1"
	default:
		err.Notes = []string{
			fmt.Sprintf("the operator '%s' is not defined for type %s", operator, typeName),
		}
		err.Help = "check the language reference for supported operators"
	}

	return err
}

// UnexpectedTokenError creates an unexpected token error with helpful context
func UnexpectedTokenError(expected, found string, loc SourceLocation, source string) *VictoriaError {
	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0004",
		Message:    fmt.Sprintf("expected '%s' but found '%s'", expected, found),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("expected '%s' here", expected), Primary: true},
		},
	}

	// Context-specific help
	switch expected {
	case "=":
		err.Notes = []string{
			"variable declarations require an initializer",
		}
		err.Help = "add an initial value: let variable = value"
	case ")":
		err.Notes = []string{
			"every opening parenthesis '(' must have a matching closing ')'",
		}
		err.Help = "check for missing closing parenthesis in function calls or expressions"
	case "}":
		err.Notes = []string{
			"every opening brace '{' must have a matching closing '}'",
		}
		err.Help = "check for missing closing brace in blocks, functions, or hashes"
	case "]":
		err.Notes = []string{
			"every opening bracket '[' must have a matching closing ']'",
		}
		err.Help = "check for missing closing bracket in array literals or index expressions"
	case ";", "NEWLINE":
		err.Notes = []string{
			"statements should be separated by newlines",
		}
		err.Help = "Victoria doesn't require semicolons; just use a new line"
	case "identifier":
		err.Notes = []string{
			"an identifier (variable name) was expected here",
		}
		err.Help = "variable names must start with a letter and contain only letters, numbers, and underscores"
	default:
		err.Notes = []string{
			"the parser encountered an unexpected token",
		}
		err.Help = "check the syntax around this location"
	}

	return err
}

// NotAFunctionError creates a not a function error
func NotAFunctionError(typeName string, loc SourceLocation, source string) *VictoriaError {
	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0005",
		Message:    fmt.Sprintf("'%s' is not a function", typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "cannot be called as a function", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("found type %s, but expected FUNCTION or BUILTIN", typeName),
			"only functions defined with fn() { } can be called",
		},
	}

	switch typeName {
	case "INTEGER":
		err.Help = "remove the parentheses, or use a function that returns an integer"
	case "STRING":
		err.Help = "strings cannot be called; use string methods like upper(), lower(), split()"
	case "ARRAY":
		err.Help = "use array[index] to access elements, not array(index)"
	case "HASH":
		err.Help = "use hash[\"key\"] or hash.key to access values, not hash()"
	case "BOOLEAN":
		err.Help = "booleans cannot be called; use boolean expressions with 'and', 'or', '!'"
	default:
		err.Help = "ensure the variable contains a function before calling it"
	}

	return err
}

// IndexOutOfBoundsError creates an index out of bounds error
func IndexOutOfBoundsError(index, length int64, loc SourceLocation, source string) *VictoriaError {
	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0006",
		Message:    fmt.Sprintf("index out of bounds: index is %d but length is %d", index, length),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("index %d is out of range", index), Primary: true},
		},
	}

	if length == 0 {
		err.Notes = []string{
			"the array/string is empty (length 0)",
			"there are no valid indices for an empty collection",
		}
		err.Help = "check if the collection is empty with len() before accessing"
	} else {
		err.Notes = []string{
			fmt.Sprintf("valid indices are 0 to %d (inclusive)", length-1),
			"Victoria uses zero-based indexing",
		}
		if index < 0 {
			err.Help = "negative indices are not supported; use len(arr) - 1 for the last element"
		} else {
			err.Help = fmt.Sprintf("use an index between 0 and %d", length-1)
		}
	}

	return err
}

// DivisionByZeroError creates a division by zero error
func DivisionByZeroError(loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0007",
		Message:    "division by zero",
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "divisor is zero here", Primary: true},
		},
		Notes: []string{
			"dividing by zero is undefined in mathematics",
			"this error occurs at runtime when the divisor evaluates to 0",
		},
		Help: "add a check: if divisor != 0 { result = x / divisor }",
	}
}

// PropertyNotFoundError creates a property not found error
func PropertyNotFoundError(property, typeName string, loc SourceLocation, source string) *VictoriaError {
	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0008",
		Message:    fmt.Sprintf("property '%s' not found on type %s", property, typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "unknown property or method", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("type %s does not have a property named '%s'", typeName, property),
		},
	}

	switch typeName {
	case "HASH":
		err.Help = fmt.Sprintf("check if key exists: if hash[\"%s\"] != null { ... }", property)
		err.Notes = append(err.Notes, "use keys(hash) to see all available keys")
	case "STRUCT_INSTANCE":
		err.Help = "check the struct definition for available fields"
		err.Notes = append(err.Notes, "struct fields must be defined when the struct is created")
	case "ARRAY":
		err.Help = "arrays don't have properties; use len(), first(), last(), etc."
		err.Notes = append(err.Notes, "common array functions: push, pop, first, last, rest, len")
	case "STRING":
		err.Help = "strings don't have properties; use string functions instead"
		err.Notes = append(err.Notes, "common string functions: upper, lower, split, len, contains")
	default:
		err.Help = "check that the property name is spelled correctly"
	}

	return err
}

// StructNotFoundError creates a struct not found error
func StructNotFoundError(name string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0009",
		Message:    fmt.Sprintf("struct '%s' not found", name),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "undefined struct type", Primary: true},
		},
		Notes: []string{
			"structs must be defined before they can be instantiated",
			fmt.Sprintf("no struct named '%s' exists in the current scope", name),
		},
		Help: fmt.Sprintf("define the struct first:\n       struct %s {\n           field1,\n           field2\n       }", name),
	}
}

// InvalidArgumentError creates an invalid argument error
func InvalidArgumentError(funcName string, expected, got int, loc SourceLocation, source string) *VictoriaError {
	var message string
	if expected == got {
		message = fmt.Sprintf("wrong number of arguments to '%s'", funcName)
	} else if got < expected {
		message = fmt.Sprintf("too few arguments to '%s': expected %d, got %d", funcName, expected, got)
	} else {
		message = fmt.Sprintf("too many arguments to '%s': expected %d, got %d", funcName, expected, got)
	}

	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0010",
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("expected %d argument(s), found %d", expected, got), Primary: true},
		},
	}

	// Add function-specific help
	functionHelp := map[string]string{
		"len":      "len(collection) - returns the length of an array, string, or hash",
		"push":     "push(array, element) - adds an element to the end of an array",
		"pop":      "pop(array) - removes and returns the last element",
		"first":    "first(array) - returns the first element",
		"last":     "last(array) - returns the last element",
		"rest":     "rest(array) - returns all elements except the first",
		"split":    "split(string, delimiter) - splits a string into an array",
		"join":     "join(array, separator) - joins array elements into a string",
		"upper":    "upper(string) - converts string to uppercase",
		"lower":    "lower(string) - converts string to lowercase",
		"contains": "contains(collection, element) - checks if element exists",
		"index":    "index(collection, element) - finds the index of an element",
		"map":      "map(array, fn) - applies fn to each element",
		"filter":   "filter(array, fn) - keeps elements where fn returns true",
		"reduce":   "reduce(array, fn, [initial]) - reduces array to single value",
		"range":    "range(end) or range(start, end) or range(start, end, step)",
		"format":   "format(template, ...values) - formats a string with values",
		"int":      "int(value) - converts value to integer",
		"string":   "string(value) - converts value to string",
		"type":     "type(value) - returns the type of a value as a string",
		"keys":     "keys(hash) - returns array of hash keys",
		"values":   "values(hash) - returns array of hash values",
		"print":    "print(...values) - prints values to stdout",
		"input":    "input([prompt]) - reads a line from stdin",
	}

	if help, ok := functionHelp[funcName]; ok {
		err.Help = fmt.Sprintf("usage: %s", help)
	}

	return err
}

// ParseError creates a generic parse error
func ParseError(message string, loc SourceLocation, source string) *VictoriaError {
	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0100",
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "", Primary: true},
		},
	}

	// Add context-specific notes for common parse errors
	if strings.Contains(message, "expected") {
		err.Notes = []string{"the parser encountered an unexpected token"}
		if strings.Contains(message, "=") {
			err.Help = "check that variable declarations use 'let' or 'const'"
		} else if strings.Contains(message, ")") {
			err.Help = "ensure all opening parentheses '(' have matching closing parentheses ')'"
		} else if strings.Contains(message, "}") {
			err.Help = "ensure all opening braces '{' have matching closing braces '}'"
		} else if strings.Contains(message, "]") {
			err.Help = "ensure all opening brackets '[' have matching closing brackets ']'"
		}
	}

	return err
}

// IllegalCharacterError creates an illegal character error
func IllegalCharacterError(char string, loc SourceLocation, source string) *VictoriaError {
	err := &VictoriaError{
		Kind:       KindError,
		Code:       "E0101",
		Message:    fmt.Sprintf("illegal character: '%s'", char),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "this character is not valid here", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("the character '%s' (ASCII %d) is not recognized", char, []byte(char)[0]),
		},
	}

	// Suggest alternatives for common issues
	switch char {
	case "@":
		err.Help = "Victoria doesn't use @ for decorators; use regular function calls"
	case "$":
		err.Help = "variable names don't need $; just use: let name = value"
		err.Notes = append(err.Notes, "$ is only used inside strings for interpolation: \"${variable}\"")
	case "#":
		err.Help = "Victoria uses // for comments, not #"
	case "`":
		err.Help = "Victoria uses double quotes for strings: \"text\""
	case "?":
		err.Help = "use the ternary operator as: condition ? value1 : value2"
	default:
		err.Help = "check for copy-paste errors or encoding issues"
	}

	return err
}

// UnterminatedStringError creates an unterminated string error
func UnterminatedStringError(loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0102",
		Message:    "unterminated string literal",
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "string starts here but never ends", Primary: true},
		},
		Notes: []string{
			"strings must begin and end with double quotes",
			"multi-line strings are supported; ensure the closing quote exists",
		},
		Help: "add a closing quote '\"' to terminate the string",
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Additional Error Constructors for Better Coverage
// ═══════════════════════════════════════════════════════════════════════════════

// SliceError creates a slice-related error
func SliceError(message string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0011",
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "invalid slice operation", Primary: true},
		},
		Notes: []string{
			"slice syntax: collection[start:end]",
			"both start and end must be integers",
		},
		Help: "example: arr[0:5] or str[2:10]",
	}
}

// SpreadError creates a spread operator error
func SpreadError(message string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0012",
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "invalid spread operation", Primary: true},
		},
		Notes: []string{
			"the spread operator (...) unpacks array elements",
			"it can only be used inside array literals: [...arr]",
		},
		Help: "example: let combined = [...arr1, ...arr2]",
	}
}

// HashKeyError creates an unusable hash key error
func HashKeyError(typeName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0013",
		Message:    fmt.Sprintf("unusable as hash key: %s", typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "cannot be used as a hash key", Primary: true},
		},
		Notes: []string{
			"hash keys must be hashable (immutable) types",
			"valid key types: STRING, INTEGER, BOOLEAN",
		},
		Help: "use a string, integer, or boolean value as the key",
	}
}

// ArgumentTypeError creates an argument type error
func ArgumentTypeError(funcName, paramNum, expected, got string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0014",
		Message:    fmt.Sprintf("argument %s to '%s' must be %s, got %s", paramNum, funcName, expected, got),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("expected %s", expected), Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("'%s' requires a %s value for this argument", funcName, expected),
		},
		Help: fmt.Sprintf("convert the value to %s or use a different value", expected),
	}
}

// NotIterableError creates a not iterable error
func NotIterableError(typeName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0015",
		Message:    fmt.Sprintf("not iterable: %s", typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "cannot iterate over this value", Primary: true},
		},
		Notes: []string{
			"for-in loops require an iterable collection",
			"iterable types: ARRAY, STRING, HASH, and ranges",
		},
		Help: "use range() to iterate over numbers: for i in range(10) { ... }",
	}
}

// RangeError creates a range-related error
func RangeError(message string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0016",
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "invalid range", Primary: true},
		},
		Notes: []string{
			"range() creates a sequence of integers",
			"all range arguments must be integers",
		},
		Help: "usage: range(end), range(start, end), or range(start, end, step)",
	}
}

// ConversionError creates a type conversion error
func ConversionError(value, targetType string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0017",
		Message:    fmt.Sprintf("could not convert '%s' to %s", value, targetType),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "conversion failed", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("the value '%s' cannot be interpreted as %s", value, targetType),
		},
		Help: fmt.Sprintf("ensure the value is a valid %s representation", targetType),
	}
}

// AssignmentError creates an assignment error
func AssignmentError(message string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0018",
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "invalid assignment target", Primary: true},
		},
		Notes: []string{
			"assignments must target a variable name or index expression",
		},
		Help: "use: variable = value, array[index] = value, or hash[\"key\"] = value",
	}
}

// OperatorError creates an operator-related error
func OperatorError(message string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0019",
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "invalid operator usage", Primary: true},
		},
		Notes: []string{
			"increment (++) and decrement (--) require a variable",
		},
		Help: "use: i++ or ++i where i is a declared variable",
	}
}

// MemberAccessError creates a member access error
func MemberAccessError(message string, typeName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0020",
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "cannot access member", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("type %s does not support dot notation", typeName),
			"dot access is for hashes, structs, and objects with methods",
		},
		Help: "for dynamic keys, use bracket notation: hash[\"key\"]",
	}
}

// ============================================
// TYPE SYSTEM ERRORS
// ============================================

// TypeAnnotationMismatchError creates a type mismatch error for typed variables/parameters
func TypeAnnotationMismatchError(expected, actual, context string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0030",
		Message:    fmt.Sprintf("type mismatch: expected %s, got %s", expected, actual),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("expected %s", expected), Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("in %s: expected type '%s' but received '%s'", context, expected, actual),
			"Victoria's type system helps catch errors early",
		},
		Help: fmt.Sprintf("ensure the value is of type %s, or adjust the type annotation", expected),
	}
}

// VariableTypeMismatchError creates an error for variable type mismatches
func VariableTypeMismatchError(varName, expected, actual string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0031",
		Message:    fmt.Sprintf("cannot assign %s to variable '%s' of type %s", actual, varName, expected),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("expected %s, found %s", expected, actual), Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("variable '%s' was declared with type annotation :%s", varName, expected),
			"type annotations are enforced at runtime for type safety",
		},
		Help: fmt.Sprintf("either assign a %s value, or change the type annotation to :%s", expected, actual),
	}
}

// ParameterTypeMismatchError creates an error for function parameter type mismatches
func ParameterTypeMismatchError(paramName, funcName, expected, actual string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0032",
		Message:    fmt.Sprintf("type mismatch for parameter '%s': expected %s, got %s", paramName, expected, actual),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("wrong type for parameter '%s'", paramName), Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("function '%s' expects parameter '%s' to be of type %s", funcName, paramName, expected),
			"typed parameters enforce type checking when the function is called",
		},
		Help: fmt.Sprintf("pass a %s value, or use type conversion: %s(value)", expected, expected),
	}
}

// ReturnTypeMismatchError creates an error for function return type mismatches
func ReturnTypeMismatchError(funcName, expected, actual string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0033",
		Message:    fmt.Sprintf("return type mismatch: expected %s, got %s", expected, actual),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("returns %s, expected %s", actual, expected), Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("function '%s' has return type annotation -> %s", funcName, expected),
			"return type annotations ensure the function returns the expected type",
		},
		Help: fmt.Sprintf("return a %s value, or change the return type annotation", expected),
	}
}

// InvalidTypeAnnotationError creates an error for invalid type annotations
func InvalidTypeAnnotationError(typeName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0034",
		Message:    fmt.Sprintf("invalid type annotation: '%s'", typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "unknown type", Primary: true},
		},
		Notes: []string{
			"type annotations must be valid type names",
			"built-in types: int, float, string, bool, char, array, map, any, void",
		},
		Help: "use a built-in type or a defined struct name",
	}
}

// TypeAnnotationRequiredError creates an error when type annotation is required but missing
func TypeAnnotationRequiredError(context string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0035",
		Message:    fmt.Sprintf("type annotation required in %s", context),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "missing type annotation", Primary: true},
		},
		Notes: []string{
			"some contexts require explicit type annotations",
		},
		Help: "add a type annotation using the syntax :type (e.g., x:int)",
	}
}

// ArrayTypeMismatchError creates an error for array element type mismatches
func ArrayTypeMismatchError(expected, actual string, index int, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0036",
		Message:    fmt.Sprintf("array element type mismatch at index %d: expected %s, got %s", index, expected, actual),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("wrong element type at index %d", index), Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("typed arrays ([]%s) require all elements to be of type %s", expected, expected),
		},
		Help: fmt.Sprintf("ensure all array elements are of type %s", expected),
	}
}

// VoidReturnError creates an error for returning a value from void function
func VoidReturnError(loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0037",
		Message:    "cannot return a value from a void function",
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "unexpected return value", Primary: true},
		},
		Notes: []string{
			"functions with return type -> void should not return a value",
			"use 'return' without a value, or simply let the function end",
		},
		Help: "remove the return value, or change the return type annotation",
	}
}

// MissingReturnError creates an error for missing return in typed function
func MissingReturnError(funcName, expected string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0038",
		Message:    fmt.Sprintf("function '%s' must return a value of type %s", funcName, expected),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "missing return statement", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("function '%s' has return type -> %s", funcName, expected),
			"all code paths must return a value of the declared type",
		},
		Help: fmt.Sprintf("add a return statement that returns a %s value", expected),
	}
}

// ════════════════════════════════════════════════════════════════════════════════
// DSA-SPECIFIC ERRORS - Common mistakes in algorithms and data structures
// ════════════════════════════════════════════════════════════════════════════════

// InfiniteLoopWarning creates a warning for potential infinite loops
func InfiniteLoopWarning(reason string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindWarning,
		Code:       "W0001",
		Message:    fmt.Sprintf("potential infinite loop: %s", reason),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "this loop may never terminate", Primary: true},
		},
		Notes: []string{
			"infinite loops can freeze your program",
			"ensure your loop condition will eventually become false",
		},
		Help: "check that your loop variable is being modified inside the loop",
	}
}

// RecursionDepthError creates an error for excessive recursion
func RecursionDepthError(funcName string, depth int, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0040",
		Message:    fmt.Sprintf("maximum recursion depth exceeded in '%s' (depth: %d)", funcName, depth),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "recursion too deep", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("function '%s' called itself too many times", funcName),
			"this usually indicates missing base case or incorrect termination condition",
			"DSA tip: every recursive function needs a base case that stops the recursion",
		},
		Help: "check your base case: ensure the recursion stops for some input (e.g., n <= 0, array is empty)",
	}
}

// MemoizationSuggestion creates a suggestion to use memoization
func MemoizationSuggestion(funcName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindNote,
		Code:       "N0001",
		Message:    fmt.Sprintf("function '%s' may benefit from memoization", funcName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "consider memoizing", Primary: false},
		},
		Notes: []string{
			"memoization stores results of expensive function calls",
			"it can dramatically improve performance for recursive algorithms",
			"DSA tip: use a hash map to cache results",
		},
		Help: "add a cache: let memo = {}; if memo[key] != null { return memo[key] }; ... memo[key] = result",
	}
}

// OffByOneError creates an error for common off-by-one mistakes
func OffByOneError(context string, suggestion string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0041",
		Message:    fmt.Sprintf("off-by-one error: %s", context),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "index boundary issue", Primary: true},
		},
		Notes: []string{
			"off-by-one errors are one of the most common bugs in algorithms",
			"remember: arrays use 0-based indexing (first element is index 0)",
			"the last valid index is len(array) - 1, not len(array)",
		},
		Help: suggestion,
	}
}

// EmptyCollectionError creates an error for operations on empty collections
func EmptyCollectionError(operation string, typeName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0042",
		Message:    fmt.Sprintf("cannot %s on empty %s", operation, typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("%s is empty", typeName), Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("the %s has no elements to operate on", typeName),
			"DSA tip: always check for empty collections before accessing elements",
		},
		Help: fmt.Sprintf("add a check: if len(collection) > 0 { ... } or if len(collection) == 0 { return defaultValue }"),
	}
}

// BinarySearchError creates an error for common binary search mistakes
func BinarySearchError(mistake string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0043",
		Message:    fmt.Sprintf("binary search error: %s", mistake),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "binary search issue", Primary: true},
		},
		Notes: []string{
			"binary search requires a sorted array",
			"common mistakes: wrong mid calculation, incorrect boundary updates",
			"DSA tip: use mid = left + (right - left) / 2 to avoid integer overflow",
		},
		Help: "ensure: 1) array is sorted, 2) boundaries update correctly, 3) loop condition is left <= right",
	}
}

// GraphCycleError creates an error for cycle detection in graphs
func GraphCycleError(loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0044",
		Message:    "cycle detected in graph",
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "cycle found here", Primary: true},
		},
		Notes: []string{
			"a cycle exists in the graph structure",
			"DSA tip: use visited state tracking (UNVISITED, VISITING, VISITED) for cycle detection",
		},
		Help: "use an enum to track node states during DFS traversal",
	}
}

// SortedArrayRequiredError creates an error when sorted array is required
func SortedArrayRequiredError(operation string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0045",
		Message:    fmt.Sprintf("%s requires a sorted array", operation),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "array may not be sorted", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("'%s' assumes the input array is sorted in ascending order", operation),
			"DSA tip: sort the array first, or use a different algorithm",
		},
		Help: "ensure the array is sorted before calling this function",
	}
}

// TimeComplexityWarning creates a warning for potentially slow operations
func TimeComplexityWarning(operation string, complexity string, suggestion string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindWarning,
		Code:       "W0002",
		Message:    fmt.Sprintf("potentially slow: %s has %s complexity", operation, complexity),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("O(%s) operation", complexity), Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("this operation has %s time complexity", complexity),
			"for large inputs, this may cause performance issues",
		},
		Help: suggestion,
	}
}

// IntegerOverflowWarning creates a warning for potential integer overflow
func IntegerOverflowWarning(operation string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindWarning,
		Code:       "W0003",
		Message:    fmt.Sprintf("potential integer overflow in %s", operation),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "may overflow for large values", Primary: true},
		},
		Notes: []string{
			"integer operations can overflow for very large numbers",
			"DSA tip: use modular arithmetic to prevent overflow",
		},
		Help: "use modulo: result = (a * b) % MOD, or use #make MOD 1000000007",
	}
}

// NegativeIndexError creates an error for negative array indices
func NegativeIndexError(index int64, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0046",
		Message:    fmt.Sprintf("negative index: %d", index),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "negative indices not supported", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("index %d is negative", index),
			"Victoria arrays use zero-based positive indexing",
		},
		Help: "to access from the end, use: arr[len(arr) - 1] for the last element",
	}
}

// ConstantReassignmentError creates an error for reassigning constants
func ConstantReassignmentError(name string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0047",
		Message:    fmt.Sprintf("cannot reassign constant: '%s'", name),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "constant cannot be modified", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("'%s' was declared as a constant with 'const' or '#make'", name),
			"constants are immutable and cannot be changed after declaration",
		},
		Help: fmt.Sprintf("if you need to modify '%s', declare it with 'let' instead of 'const'", name),
	}
}

// EnumValueError creates an error for invalid enum access
func EnumValueError(enumName, valueName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0048",
		Message:    fmt.Sprintf("enum '%s' has no value '%s'", enumName, valueName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "invalid enum value", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("'%s' is not a valid member of enum '%s'", valueName, enumName),
			"DSA tip: enums are great for representing finite states (e.g., NodeState.VISITED)",
		},
		Help: fmt.Sprintf("check the enum definition for valid values: enum %s { VALUE1, VALUE2, ... }", enumName),
	}
}

// CharacterConversionError creates an error for invalid character operations
func CharacterConversionError(value string, operation string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0049",
		Message:    fmt.Sprintf("cannot %s: '%s'", operation, value),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "invalid character operation", Primary: true},
		},
		Notes: []string{
			"character operations expect single characters or valid ASCII values",
			"ord() expects a character, chr() expects an integer 0-127",
		},
		Help: "for ord(): pass a single character 'a'. For chr(): pass an integer 0-127",
	}
}

// ComparisonWithNullError creates an error for null comparisons
func ComparisonWithNullError(operator string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindWarning,
		Code:       "W0004",
		Message:    fmt.Sprintf("comparing with null using '%s' may not work as expected", operator),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "null comparison", Primary: true},
		},
		Notes: []string{
			"null is a special value representing 'no value'",
			"only == and != are meaningful for null comparisons",
		},
		Help: "use 'value == null' or 'value != null' to check for null",
	}
}

// ════════════════════════════════════════════════════════════════════════════════
// COMPETITIVE PROGRAMMING ERRORS
// ════════════════════════════════════════════════════════════════════════════════

// ModuloWithNegativeError creates an error for modulo with negative numbers
func ModuloWithNegativeError(loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindWarning,
		Code:       "W0005",
		Message:    "modulo with negative number may give unexpected results",
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "modulo with negative", Primary: true},
		},
		Notes: []string{
			"modulo behavior varies: some languages return negative, some positive",
			"DSA tip: to ensure positive result, use: ((a % m) + m) % m",
		},
		Help: "for competitive programming, normalize negative results: ((result % MOD) + MOD) % MOD",
	}
}

// MakeDirectiveError creates an error for invalid #make directive
func MakeDirectiveError(message string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0050",
		Message:    fmt.Sprintf("#make error: %s", message),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "invalid #make directive", Primary: true},
		},
		Notes: []string{
			"#make creates compile-time constants, similar to C's #define",
			"syntax: #make NAME value",
		},
		Help: "example: #make MOD 1000000007 or #make MAX_N 100005",
	}
}
