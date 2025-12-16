package token

import "testing"

func TestLookupIdent(t *testing.T) {
	tests := []struct {
		ident    string
		expected TokenType
	}{
		{"define", FUNCTION},
		{"let", LET},
		{"true", TRUE},
		{"false", FALSE},
		{"if", IF},
		{"else", ELSE},
		{"return", RETURN},
		{"struct", STRUCT},
		{"while", WHILE},
		{"for", FOR},
		{"in", IN},
		{"and", AND},
		{"or", OR},
		{"not", NOT},
		{"include", INCLUDE},
		{"try", TRY},
		{"catch", CATCH},
		{"break", BREAK},
		{"continue", CONTINUE},
		{"switch", SWITCH},
		{"case", CASE},
		{"default", DEFAULT},
		{"const", CONST},
		// Non-keywords should return IDENT
		{"foo", IDENT},
		{"bar", IDENT},
		{"myVar", IDENT},
		{"x", IDENT},
		{"_underscore", IDENT},
	}

	for _, tt := range tests {
		t.Run(tt.ident, func(t *testing.T) {
			got := LookupIdent(tt.ident)
			if got != tt.expected {
				t.Errorf("LookupIdent(%q) = %q, want %q", tt.ident, got, tt.expected)
			}
		})
	}
}

func TestTokenTypeConstants(t *testing.T) {
	// Test that token types are properly defined as unique strings
	tokenTypes := []TokenType{
		ILLEGAL, EOF, IDENT, INT, FLOAT, STRING,
		ASSIGN, PLUS, MINUS, BANG, ASTERISK, SLASH, MODULO,
		PLUS_ASSIGN, MINUS_ASSIGN, ASTERISK_ASSIGN, SLASH_ASSIGN, MODULO_ASSIGN,
		INC, DEC,
		LT, GT, EQ, NOT_EQ, LTE, GTE,
		AND, OR, NOT, AND_AND, OR_OR,
		QUESTION, RANGE, ARROW,
		COMMA, SEMICOLON, COLON, DOT,
		LPAREN, RPAREN, LBRACE, RBRACE, LBRACKET, RBRACKET,
		FUNCTION, LET, TRUE, FALSE, IF, ELSE, RETURN, STRUCT,
		WHILE, FOR, IN, INCLUDE, TRY, CATCH, BREAK, CONTINUE,
		SWITCH, CASE, DEFAULT, CONST, SPREAD,
	}

	seen := make(map[TokenType]bool)
	for _, tt := range tokenTypes {
		if tt == "" {
			t.Error("Empty token type found")
		}
		if seen[tt] {
			t.Errorf("Duplicate token type: %q", tt)
		}
		seen[tt] = true
	}
}
