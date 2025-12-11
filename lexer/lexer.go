package lexer

import (
	"strings"
	"victoria/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.line)
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUS_ASSIGN, Literal: literal, Line: l.line}
		} else if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.PLUS, l.ch, l.line)
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MINUS_ASSIGN, Literal: literal, Line: l.line}
		} else if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.MINUS, l.ch, l.line)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.BANG, l.ch, l.line)
		}
	case '/':
		if l.peekChar() == '/' {
			l.skipSingleLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipMultiLineComment()
			return l.NextToken()
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.SLASH_ASSIGN, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.SLASH, l.ch, l.line)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.ASTERISK_ASSIGN, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.ASTERISK, l.ch, l.line)
		}
	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MODULO_ASSIGN, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.MODULO, l.ch, l.line)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND_AND, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR_OR, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line)
		}
	case '?':
		tok = newToken(token.QUESTION, l.ch, l.line)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.LT, l.ch, l.line)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.GT, l.ch, l.line)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.line)
	case ':':
		tok = newToken(token.COLON, l.ch, l.line)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line)
	case '.':
		// Check if it's a range operator ..
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.RANGE, Literal: literal, Line: l.line}
		} else if isDigit(l.peekChar()) {
			// Check if it's a float starting with dot like .5
			tok.Type = token.FLOAT
			tok.Literal = l.readNumber()
			tok.Line = l.line
			return tok
		} else {
			tok = newToken(token.DOT, l.ch, l.line)
		}
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.line)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.line)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = l.line
	case '`':
		tok.Type = token.STRING
		tok.Literal = l.readMultiLineString()
		tok.Line = l.line
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.line
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Line = l.line
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func (l *Lexer) skipSingleLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	l.skipWhitespace()
}

func (l *Lexer) skipMultiLineComment() {
	// consume /*
	l.readChar()
	l.readChar()

	for l.ch != 0 {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			break
		}
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
	l.skipWhitespace()
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) { // Allow digits in identifiers after first char
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\\' && l.peekChar() != 0 {
			l.readChar() // skip escaped character
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readMultiLineString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '`' || l.ch == 0 {
			break
		}
		if l.ch == '\n' {
			l.line++
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func newToken(tokenType token.TokenType, ch byte, line int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line}
}
