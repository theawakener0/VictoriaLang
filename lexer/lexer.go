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
	column       int // current column position (1-indexed)
	lineStart    int // position where current line starts
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0, lineStart: 0}
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
	l.column = l.position - l.lineStart + 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	startCol := l.column
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.ASSIGN, l.ch, l.line, startCol)
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUS_ASSIGN, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.PLUS, l.ch, l.line, startCol)
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MINUS_ASSIGN, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.MINUS, l.ch, l.line, startCol)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.BANG, l.ch, l.line, startCol)
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
			tok = token.Token{Type: token.SLASH_ASSIGN, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.SLASH, l.ch, l.line, startCol)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.ASTERISK_ASSIGN, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.ASTERISK, l.ch, l.line, startCol)
		}
	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MODULO_ASSIGN, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.MODULO, l.ch, l.line, startCol)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND_AND, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.ILLEGAL, l.ch, l.line, startCol)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR_OR, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.ILLEGAL, l.ch, l.line, startCol)
		}
	case '?':
		tok = newTokenWithCol(token.QUESTION, l.ch, l.line, startCol)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.LT, l.ch, l.line, startCol)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else {
			tok = newTokenWithCol(token.GT, l.ch, l.line, startCol)
		}
	case ';':
		tok = newTokenWithCol(token.SEMICOLON, l.ch, l.line, startCol)
	case ':':
		tok = newTokenWithCol(token.COLON, l.ch, l.line, startCol)
	case ',':
		tok = newTokenWithCol(token.COMMA, l.ch, l.line, startCol)
	case '.':
		// Check if it's a range operator ..
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.RANGE, Literal: literal, Line: l.line, Column: startCol, EndColumn: l.column + 1}
		} else if isDigit(l.peekChar()) {
			// Check if it's a float starting with dot like .5
			startCol := l.column
			tok.Literal = l.readNumber()
			tok.Type = token.FLOAT
			tok.Line = l.line
			tok.Column = startCol
			tok.EndColumn = l.column
			return tok
		} else {
			tok = newTokenWithCol(token.DOT, l.ch, l.line, startCol)
		}
	case '{':
		tok = newTokenWithCol(token.LBRACE, l.ch, l.line, startCol)
	case '}':
		tok = newTokenWithCol(token.RBRACE, l.ch, l.line, startCol)
	case '(':
		tok = newTokenWithCol(token.LPAREN, l.ch, l.line, startCol)
	case ')':
		tok = newTokenWithCol(token.RPAREN, l.ch, l.line, startCol)
	case '[':
		tok = newTokenWithCol(token.LBRACKET, l.ch, l.line, startCol)
	case ']':
		tok = newTokenWithCol(token.RBRACKET, l.ch, l.line, startCol)
	case '"':
		tok.Column = startCol
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = l.line
		tok.EndColumn = l.column + 1
	case '`':
		tok.Column = startCol
		tok.Type = token.STRING
		tok.Literal = l.readMultiLineString()
		tok.Line = l.line
		tok.EndColumn = l.column + 1
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Column = l.column
		tok.EndColumn = l.column
	default:
		if isLetter(l.ch) {
			startCol := l.column
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = startCol
			tok.EndColumn = l.column
			return tok
		} else if isDigit(l.ch) {
			startCol := l.column
			tok.Literal = l.readNumber()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Line = l.line
			tok.Column = startCol
			tok.EndColumn = l.column
			return tok
		} else {
			tok = newTokenWithCol(token.ILLEGAL, l.ch, l.line, startCol)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.lineStart = l.readPosition
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
			l.lineStart = l.readPosition
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
			l.lineStart = l.readPosition
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
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: 1, EndColumn: 2}
}

func newTokenWithCol(tokenType token.TokenType, ch byte, line int, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col, EndColumn: col + 1}
}
