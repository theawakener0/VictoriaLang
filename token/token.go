package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
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

	QUESTION = "?"
	RANGE    = ".."

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
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
