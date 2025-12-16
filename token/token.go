package token

type TokenType string

type Token struct {
	Type      TokenType
	Literal   string
	Line      int
	Column    int
	EndColumn int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"   // 1343456
	FLOAT  = "FLOAT" // 12.34
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"

	PLUS_ASSIGN     = "+="
	MINUS_ASSIGN    = "-="
	ASTERISK_ASSIGN = "*="
	SLASH_ASSIGN    = "/="
	MODULO_ASSIGN   = "%="

	INC = "++"
	DEC = "--"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="
	LTE    = "<="
	GTE    = ">="

	AND     = "and"
	OR      = "or"
	NOT     = "not"
	AND_AND = "&&"
	OR_OR   = "||"

	QUESTION     = "?"
	RANGE        = ".."
	ARROW        = "=>"
	ARROW_RETURN = "->" // For function return type annotation

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	STRUCT   = "STRUCT"
	WHILE    = "WHILE"
	FOR      = "FOR"
	IN       = "IN"
	INCLUDE  = "INCLUDE"
	TRY      = "TRY"
	CATCH    = "CATCH"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	SWITCH   = "SWITCH"
	CASE     = "CASE"
	DEFAULT  = "DEFAULT"
	CONST    = "CONST"

	// Type keywords
	TYPE_INT    = "TYPE_INT"
	TYPE_FLOAT  = "TYPE_FLOAT"
	TYPE_STRING = "TYPE_STRING"
	TYPE_BOOL   = "TYPE_BOOL"
	TYPE_CHAR   = "TYPE_CHAR"
	TYPE_ARRAY  = "TYPE_ARRAY"
	TYPE_MAP    = "TYPE_MAP"
	TYPE_ANY    = "TYPE_ANY"
	TYPE_VOID   = "TYPE_VOID"

	// Special operators
	SPREAD = "..."
)

var keywords = map[string]TokenType{
	"define":   FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"struct":   STRUCT,
	"while":    WHILE,
	"for":      FOR,
	"in":       IN,
	"and":      AND,
	"or":       OR,
	"not":      NOT,
	"include":  INCLUDE,
	"try":      TRY,
	"catch":    CATCH,
	"break":    BREAK,
	"continue": CONTINUE,
	"switch":   SWITCH,
	"case":     CASE,
	"default":  DEFAULT,
	"const":    CONST,
	// Type keywords
	"int":    TYPE_INT,
	"float":  TYPE_FLOAT,
	"string": TYPE_STRING,
	"bool":   TYPE_BOOL,
	"char":   TYPE_CHAR,
	"array":  TYPE_ARRAY,
	"map":    TYPE_MAP,
	"any":    TYPE_ANY,
	"void":   TYPE_VOID,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// IsTypeKeyword returns true if the token type is a type keyword
func IsTypeKeyword(t TokenType) bool {
	switch t {
	case TYPE_INT, TYPE_FLOAT, TYPE_STRING, TYPE_BOOL, TYPE_CHAR, TYPE_ARRAY, TYPE_MAP, TYPE_ANY, TYPE_VOID:
		return true
	}
	return false
}

// TypeKeywordToString returns the string representation of a type keyword
func TypeKeywordToString(t TokenType) string {
	switch t {
	case TYPE_INT:
		return "int"
	case TYPE_FLOAT:
		return "float"
	case TYPE_STRING:
		return "string"
	case TYPE_BOOL:
		return "bool"
	case TYPE_CHAR:
		return "char"
	case TYPE_ARRAY:
		return "array"
	case TYPE_MAP:
		return "map"
	case TYPE_ANY:
		return "any"
	case TYPE_VOID:
		return "void"
	}
	return ""
}
