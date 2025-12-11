package errors

import (
	"fmt"
	"strings"
)

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

// getSourceLine extracts a specific line from source code
func getSourceLine(source string, lineNum int) string {
	lines := strings.Split(source, "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return ""
	}
	return lines[lineNum-1]
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

	return sb.String()
}

// FormatPlain returns the error message without colors (for logging)
func (e *VictoriaError) FormatPlain() string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("%s", e.Kind.String()))
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

// Common error constructors for Victoria language

// TypeMismatchError creates a type mismatch error
func TypeMismatchError(left, operator, right string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
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
		},
		Help: fmt.Sprintf("consider converting one operand to match the other type"),
	}
}

// UndefinedVariableError creates an undefined variable error
func UndefinedVariableError(name string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0002",
		Message:    fmt.Sprintf("undefined variable: '%s'", name),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "not found in this scope", Primary: true},
		},
		Help: fmt.Sprintf("did you mean to declare it with 'let %s = ...'?", name),
	}
}

// UnknownOperatorError creates an unknown operator error
func UnknownOperatorError(operator, typeName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0003",
		Message:    fmt.Sprintf("unknown operator: '%s' for type %s", operator, typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "unsupported operator", Primary: true},
		},
	}
}

// UnexpectedTokenError creates an unexpected token error
func UnexpectedTokenError(expected, found string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0004",
		Message:    fmt.Sprintf("expected '%s' but found '%s'", expected, found),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("expected '%s' here", expected), Primary: true},
		},
	}
}

// NotAFunctionError creates a not a function error
func NotAFunctionError(typeName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0005",
		Message:    fmt.Sprintf("'%s' is not a function", typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "not callable", Primary: true},
		},
		Help: "only functions and builtin functions can be called",
	}
}

// IndexOutOfBoundsError creates an index out of bounds error
func IndexOutOfBoundsError(index, length int64, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0006",
		Message:    fmt.Sprintf("index out of bounds: index is %d but length is %d", index, length),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "index out of range", Primary: true},
		},
		Notes: []string{
			fmt.Sprintf("valid indices for this array are 0 to %d", length-1),
		},
	}
}

// DivisionByZeroError creates a division by zero error
func DivisionByZeroError(loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0007",
		Message:    "division by zero",
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "cannot divide by zero", Primary: true},
		},
		Help: "ensure the divisor is not zero before dividing",
	}
}

// PropertyNotFoundError creates a property not found error
func PropertyNotFoundError(property, typeName string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0008",
		Message:    fmt.Sprintf("property '%s' not found on type %s", property, typeName),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "unknown property", Primary: true},
		},
	}
}

// StructNotFoundError creates a struct not found error
func StructNotFoundError(name string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0009",
		Message:    fmt.Sprintf("struct '%s' not found", name),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "undefined struct", Primary: true},
		},
		Help: fmt.Sprintf("did you forget to define 'struct %s { ... }'?", name),
	}
}

// InvalidArgumentError creates an invalid argument error
func InvalidArgumentError(funcName string, expected, got int, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0010",
		Message:    fmt.Sprintf("wrong number of arguments: %s expects %d, got %d", funcName, expected, got),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: fmt.Sprintf("expected %d argument(s)", expected), Primary: true},
		},
	}
}

// ParseError creates a generic parse error
func ParseError(message string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0100",
		Message:    message,
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "", Primary: true},
		},
	}
}

// IllegalCharacterError creates an illegal character error
func IllegalCharacterError(char string, loc SourceLocation, source string) *VictoriaError {
	return &VictoriaError{
		Kind:       KindError,
		Code:       "E0101",
		Message:    fmt.Sprintf("illegal character: '%s'", char),
		SourceCode: source,
		Labels: []Label{
			{Location: loc, Message: "unexpected character", Primary: true},
		},
	}
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
		Help: "add a closing quote '\"' to terminate the string",
	}
}
