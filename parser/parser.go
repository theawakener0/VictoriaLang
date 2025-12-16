package parser

import (
	"fmt"
	"strconv"
	"victoria/ast"
	"victoria/errors"
	"victoria/lexer"
	"victoria/token"
)

const (
	_ int = iota
	LOWEST
	ARROW_PREC  // =>
	TERNARY     // ?:
	OR_PREC     // || or
	AND_PREC    // && and
	ASSIGN      // =
	EQUALS      // ==
	LESSGREATER // > or <
	RANGE_PREC  // ..
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	POSTFIX     // i++
	DOT         // struct.field
)

var precedences = map[token.TokenType]int{
	token.ARROW:           ARROW_PREC,
	token.EQ:              EQUALS,
	token.NOT_EQ:          EQUALS,
	token.LT:              LESSGREATER,
	token.GT:              LESSGREATER,
	token.LTE:             LESSGREATER,
	token.GTE:             LESSGREATER,
	token.PLUS:            SUM,
	token.MINUS:           SUM,
	token.SLASH:           PRODUCT,
	token.ASTERISK:        PRODUCT,
	token.MODULO:          PRODUCT,
	token.LPAREN:          CALL,
	token.LBRACKET:        INDEX,
	token.DOT:             DOT,
	token.ASSIGN:          ASSIGN,
	token.PLUS_ASSIGN:     ASSIGN,
	token.MINUS_ASSIGN:    ASSIGN,
	token.ASTERISK_ASSIGN: ASSIGN,
	token.SLASH_ASSIGN:    ASSIGN,
	token.MODULO_ASSIGN:   ASSIGN,
	token.INC:             POSTFIX,
	token.DEC:             POSTFIX,
	token.AND:             AND_PREC,
	token.OR:              OR_PREC,
	token.AND_AND:         AND_PREC,
	token.OR_OR:           OR_PREC,
	token.QUESTION:        TERNARY,
	token.RANGE:           RANGE_PREC,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l          *lexer.Lexer
	errors     []string
	richErrors []*errors.VictoriaError
	sourceCode string
	filename   string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:          l,
		errors:     []string{},
		richErrors: []*errors.VictoriaError{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral) // Can be hash or block, but in expression context usually hash or struct init? No, struct init starts with IDENT.
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.TRY, p.parseTryExpression)
	p.registerPrefix(token.SWITCH, p.parseSwitchExpression)
	p.registerPrefix(token.INC, p.parsePrefixIncDec)
	p.registerPrefix(token.DEC, p.parsePrefixIncDec)
	p.registerPrefix(token.SPREAD, p.parseSpreadExpression)

	// Allow type keywords to be used as identifiers in expression context
	// This allows builtin functions like string(), int(), float(), bool() to work
	p.registerPrefix(token.TYPE_STRING, p.parseTypeKeywordAsIdentifier)
	p.registerPrefix(token.TYPE_INT, p.parseTypeKeywordAsIdentifier)
	p.registerPrefix(token.TYPE_FLOAT, p.parseTypeKeywordAsIdentifier)
	p.registerPrefix(token.TYPE_BOOL, p.parseTypeKeywordAsIdentifier)
	p.registerPrefix(token.TYPE_CHAR, p.parseTypeKeywordAsIdentifier)
	p.registerPrefix(token.TYPE_ARRAY, p.parseTypeKeywordAsIdentifier)
	p.registerPrefix(token.TYPE_MAP, p.parseTypeKeywordAsIdentifier)
	p.registerPrefix(token.TYPE_ANY, p.parseTypeKeywordAsIdentifier)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseDotExpression) // For method calls or field access
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.PLUS_ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.MINUS_ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK_ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.SLASH_ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.MODULO_ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.AND_AND, p.parseInfixExpression)
	p.registerInfix(token.OR_OR, p.parseInfixExpression)
	p.registerInfix(token.QUESTION, p.parseTernaryExpression)
	p.registerInfix(token.RANGE, p.parseRangeExpression)
	p.registerInfix(token.ARROW, p.parseArrowFunction)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

// RichErrors returns the rich error objects for formatted output
func (p *Parser) RichErrors() []*errors.VictoriaError {
	return p.richErrors
}

// SetSource sets the source code for error reporting
func (p *Parser) SetSource(source string, filename string) {
	p.sourceCode = source
	p.filename = filename
}

// HasErrors returns true if there are any errors
func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0 || len(p.richErrors) > 0
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected '%s' but found '%s'", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)

	loc := errors.SourceLocation{
		Line:      p.peekToken.Line,
		Column:    p.peekToken.Column,
		EndColumn: p.peekToken.EndColumn,
		Filename:  p.filename,
	}
	richErr := errors.UnexpectedTokenError(string(t), string(p.peekToken.Type), loc, p.sourceCode)

	// Add context-specific help
	switch t {
	case token.RBRACE:
		_ = richErr.WithHelp("you might be missing a closing brace '}'")
	case token.RPAREN:
		_ = richErr.WithHelp("you might be missing a closing parenthesis ')'")
	case token.RBRACKET:
		_ = richErr.WithHelp("you might be missing a closing bracket ']'")
	case token.ASSIGN:
		_ = richErr.WithHelp("variable declarations require an initial value: let name = value")
	case token.LBRACE:
		_ = richErr.WithHelp("expected a block starting with '{'")
	case token.IDENT:
		_ = richErr.WithHelp("expected an identifier (variable or function name)")
	}

	// Add note about what was found
	if p.peekToken.Type == token.EOF {
		_ = richErr.WithNote("reached end of file unexpectedly")
	}

	p.richErrors = append(p.richErrors, richErr)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.CONST:
		return p.parseConstStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.INCLUDE:
		return p.parseIncludeStatement()
	case token.TRY:
		return p.parseTryStatement()
	case token.STRUCT:
		return p.parseStructStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	case token.FUNCTION:
		// Check if it is a method definition: def Struct.Method()
		if p.peekTokenIs(token.IDENT) {
			return p.parseFunctionOrMethodDeclaration()
		}
		return nil // Should not happen if syntax is correct
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for optional type annotation: let x:int = ...
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // consume ':'
		stmt.Type = p.parseTypeAnnotation()
		if stmt.Type == nil {
			return nil
		}
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for optional type annotation: const x:int = ...
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // consume ':'
		stmt.Type = p.parseTypeAnnotation()
		if stmt.Type == nil {
			return nil
		}
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIncludeStatement() *ast.IncludeStatement {
	stmt := &ast.IncludeStatement{Token: p.curToken}
	stmt.Modules = []string{}

	if p.peekTokenIs(token.LPAREN) {
		// include ("os", "net")
		p.nextToken() // consume include
		p.nextToken() // consume (

		// Parse first string
		if !p.curTokenIs(token.STRING) {
			return nil
		}
		stmt.Modules = append(stmt.Modules, p.curToken.Literal)

		for p.peekTokenIs(token.COMMA) {
			p.nextToken() // consume string
			p.nextToken() // consume comma
			if !p.curTokenIs(token.STRING) {
				return nil
			}
			stmt.Modules = append(stmt.Modules, p.curToken.Literal)
		}

		if !p.expectPeek(token.RPAREN) {
			return nil
		}
	} else {
		// include "os"
		if !p.expectPeek(token.STRING) {
			return nil
		}
		stmt.Modules = append(stmt.Modules, p.curToken.Literal)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseTryExpression() ast.Expression {
	return p.parseTryStatement()
}

func (p *Parser) parseTryStatement() *ast.TryStatement {
	stmt := &ast.TryStatement{Token: p.curToken}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Block = p.parseBlockStatement()

	if p.peekTokenIs(token.CATCH) {
		p.nextToken() // consume CATCH

		if p.peekTokenIs(token.LPAREN) {
			p.nextToken() // consume (
			if !p.expectPeek(token.IDENT) {
				return nil
			}
			stmt.CatchVar = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			if !p.expectPeek(token.RPAREN) {
				return nil
			}
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		stmt.CatchBlock = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseStructStatement() *ast.StructLiteral {
	stmt := &ast.StructLiteral{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Fields = []*ast.Identifier{}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		if p.curTokenIs(token.IDENT) {
			stmt.Fields = append(stmt.Fields, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		}
		// Optional commas or newlines (handled by lexer skipping whitespace)
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return stmt
}

func (p *Parser) parseFunctionOrMethodDeclaration() ast.Statement {
	// curToken is 'define'
	defToken := p.curToken
	p.nextToken() // consume 'define', now at Name or StructName

	firstIdent := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.DOT) {
		// Method definition: define Student.greet()
		p.nextToken() // consume Name
		p.nextToken() // consume '.'

		methodName := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		methodDef := &ast.MethodDefinition{
			Token:      defToken,
			StructName: firstIdent,
			MethodName: methodName,
		}

		if !p.expectPeek(token.LPAREN) {
			return nil
		}

		// Parse parameters with optional type annotations
		params, typedParams := p.parseTypedFunctionParameters()
		methodDef.Parameters = params
		methodDef.TypedParameters = typedParams

		// Check for return type annotation: -> type
		if p.peekTokenIs(token.ARROW_RETURN) {
			p.nextToken() // consume '->'
			methodDef.ReturnTypes = p.parseReturnTypes()
			if methodDef.ReturnTypes == nil {
				return nil
			}
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		methodDef.Body = p.parseBlockStatement()
		return methodDef

	} else {

		defineLit := &ast.FunctionLiteral{Token: defToken}
		defineLit.Name = firstIdent.Value

		if !p.expectPeek(token.LPAREN) {
			return nil
		}

		// Parse parameters with optional type annotations
		params, typedParams := p.parseTypedFunctionParameters()
		defineLit.Parameters = params
		defineLit.TypedParameters = typedParams

		// Check for return type annotation: -> type
		if p.peekTokenIs(token.ARROW_RETURN) {
			p.nextToken() // consume '->'
			defineLit.ReturnTypes = p.parseReturnTypes()
			if defineLit.ReturnTypes == nil {
				return nil
			}
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		defineLit.Body = p.parseBlockStatement()

		// Wrap in LetStatement
		letStmt := &ast.LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let", Line: defToken.Line},
			Name:  firstIdent,
			Value: defineLit,
		}
		return letStmt
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		// Special case: Struct instantiation looks like IDENT { ... }
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseIdentifier() ast.Expression {
	// Check for Struct Instantiation: Student { ... }
	if p.peekTokenIs(token.LBRACE) {
		return p.parseStructInstantiation()
	}
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseTypeKeywordAsIdentifier allows type keywords (string, int, etc.) to be used
// as identifiers in expression context, enabling builtin functions like string(), int()
func (p *Parser) parseTypeKeywordAsIdentifier() ast.Expression {
	// Convert the type keyword token to an identifier
	identToken := token.Token{
		Type:      token.IDENT,
		Literal:   token.TypeKeywordToString(p.curToken.Type),
		Line:      p.curToken.Line,
		Column:    p.curToken.Column,
		EndColumn: p.curToken.EndColumn,
	}
	return &ast.Identifier{Token: identToken, Value: identToken.Literal}
}

func (p *Parser) parseStructInstantiation() ast.Expression {
	// curToken is the Struct Name
	structName := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // consume name, now at {

	si := &ast.StructInstantiation{Token: structName.Token, Name: structName}
	si.Fields = make(map[string]ast.Expression)

	// parse { key: value, ... }

	if !p.curTokenIs(token.LBRACE) {
		return nil
	}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken() // move to key
		if !p.curTokenIs(token.IDENT) && !p.curTokenIs(token.STRING) {
			// Error
			return nil
		}
		key := p.curToken.Literal

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken() // move to value
		value := p.parseExpression(LOWEST)
		si.Fields[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return si
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)

		loc := errors.SourceLocation{
			Line:      p.curToken.Line,
			Column:    p.curToken.Column,
			EndColumn: p.curToken.EndColumn,
			Filename:  p.filename,
		}
		richErr := errors.ParseError(fmt.Sprintf("invalid integer literal '%s'", p.curToken.Literal), loc, p.sourceCode).
			WithCode("E0103").
			WithHelp("integers must be valid numeric values within the supported range")
		p.richErrors = append(p.richErrors, richErr)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)

		loc := errors.SourceLocation{
			Line:      p.curToken.Line,
			Column:    p.curToken.Column,
			EndColumn: p.curToken.EndColumn,
			Filename:  p.filename,
		}
		richErr := errors.ParseError(fmt.Sprintf("invalid float literal '%s'", p.curToken.Literal), loc, p.sourceCode).
			WithCode("E0104").
			WithHelp("floats must be valid numeric values like 3.14 or 0.5")
		p.richErrors = append(p.richErrors, richErr)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	startToken := p.curToken
	p.nextToken()

	// Check if this might be a lambda parameter list: (x) => or (x, y) =>
	// First, check if it's just an identifier followed by ) and =>
	if p.curTokenIs(token.IDENT) {
		firstIdent := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if p.peekTokenIs(token.RPAREN) {
			// Single parameter: (x) => expr
			p.nextToken() // consume )
			if p.peekTokenIs(token.ARROW) {
				p.nextToken() // consume =>
				arrowToken := p.curToken
				p.nextToken() // move to body
				body := p.parseExpression(ARROW_PREC)
				return &ast.ArrowFunction{
					Token:      arrowToken,
					Parameters: []*ast.Identifier{firstIdent},
					Body:       body,
				}
			}
			// Not a lambda, just (x), return the identifier
			return firstIdent
		} else if p.peekTokenIs(token.COMMA) {
			// Multi-parameter: (x, y, ...) => expr
			params := []*ast.Identifier{firstIdent}
			for p.peekTokenIs(token.COMMA) {
				p.nextToken() // consume comma
				p.nextToken() // move to next ident
				if !p.curTokenIs(token.IDENT) {
					// Not a valid parameter list, need to backtrack...
					// This is complex - for simplicity, we assume it's a lambda
					return nil
				}
				params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			}
			if !p.expectPeek(token.RPAREN) {
				return nil
			}
			if p.peekTokenIs(token.ARROW) {
				p.nextToken() // consume =>
				arrowToken := p.curToken
				p.nextToken() // move to body
				body := p.parseExpression(ARROW_PREC)
				return &ast.ArrowFunction{
					Token:      arrowToken,
					Parameters: params,
					Body:       body,
				}
			}
			// Not a lambda, this is an error - can't have (x, y) without =>
			return nil
		}
	}

	// Check for empty parameter list: () => expr
	if p.curTokenIs(token.RPAREN) {
		if p.peekTokenIs(token.ARROW) {
			p.nextToken() // consume =>
			arrowToken := p.curToken
			p.nextToken() // move to body
			body := p.parseExpression(ARROW_PREC)
			return &ast.ArrowFunction{
				Token:      arrowToken,
				Parameters: []*ast.Identifier{},
				Body:       body,
			}
		}
		// Empty parens without arrow - likely an error
		return nil
	}

	// Regular grouped expression
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Check if followed by => (for non-identifier expressions wrapped in parens)
	if p.peekTokenIs(token.ARROW) {
		if ident, ok := exp.(*ast.Identifier); ok {
			p.nextToken() // consume =>
			arrowToken := p.curToken
			p.nextToken() // move to body
			body := p.parseExpression(ARROW_PREC)
			return &ast.ArrowFunction{
				Token:      arrowToken,
				Parameters: []*ast.Identifier{ident},
				Body:       body,
			}
		}
	}

	_ = startToken // suppress unused warning
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			// else if
			p.nextToken() // consume 'if'
			// Recursively parse the if expression
			expression.Alternative = &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Expression: p.parseIfExpression()},
				},
			}
		} else {
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Parse parameters with optional type annotations
	params, typedParams := p.parseTypedFunctionParameters()
	lit.Parameters = params
	lit.TypedParameters = typedParams

	// Check for return type annotation: -> type
	if p.peekTokenIs(token.ARROW_RETURN) {
		p.nextToken() // consume '->'
		lit.ReturnTypes = p.parseReturnTypes()
		if lit.ReturnTypes == nil {
			return nil
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseTypeAnnotation parses a type annotation after a colon (:)
// Supports: int, string, bool, float, char, []int (arrays), map[string]int, and custom types
func (p *Parser) parseTypeAnnotation() *ast.TypeAnnotation {
	p.nextToken() // consume the current token to move to the type

	typeAnn := &ast.TypeAnnotation{Token: p.curToken}

	// Check for array type: []type
	if p.curTokenIs(token.LBRACKET) {
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
		typeAnn.IsArray = true
		p.nextToken() // move to element type
		typeAnn.ElementType = &ast.TypeAnnotation{Token: p.curToken}
		if token.IsTypeKeyword(p.curToken.Type) {
			typeAnn.ElementType.TypeName = token.TypeKeywordToString(p.curToken.Type)
		} else if p.curTokenIs(token.IDENT) {
			typeAnn.ElementType.TypeName = p.curToken.Literal
		} else {
			msg := fmt.Sprintf("expected type after '[]', got %s", p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
		return typeAnn
	}

	// Check for map type: map[keyType]valueType
	if p.curTokenIs(token.TYPE_MAP) {
		if !p.expectPeek(token.LBRACKET) {
			return nil
		}
		p.nextToken() // move to key type
		typeAnn.KeyType = &ast.TypeAnnotation{Token: p.curToken}
		if token.IsTypeKeyword(p.curToken.Type) {
			typeAnn.KeyType.TypeName = token.TypeKeywordToString(p.curToken.Type)
		} else if p.curTokenIs(token.IDENT) {
			typeAnn.KeyType.TypeName = p.curToken.Literal
		} else {
			msg := fmt.Sprintf("expected type in map key, got %s", p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
		p.nextToken() // move to value type
		typeAnn.ElementType = &ast.TypeAnnotation{Token: p.curToken}
		if token.IsTypeKeyword(p.curToken.Type) {
			typeAnn.ElementType.TypeName = token.TypeKeywordToString(p.curToken.Type)
		} else if p.curTokenIs(token.IDENT) {
			typeAnn.ElementType.TypeName = p.curToken.Literal
		} else {
			msg := fmt.Sprintf("expected type after map key type, got %s", p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
		typeAnn.TypeName = "map"
		return typeAnn
	}

	// Handle basic types or custom type names (IDENT)
	if token.IsTypeKeyword(p.curToken.Type) {
		typeAnn.TypeName = token.TypeKeywordToString(p.curToken.Type)
	} else if p.curTokenIs(token.IDENT) {
		// Custom type like a struct name
		typeAnn.TypeName = p.curToken.Literal
	} else {
		msg := fmt.Sprintf("expected type annotation, got %s", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	return typeAnn
}

// parseTypedFunctionParameters parses function parameters with type annotations
// e.g., (x:int, y:string) or (x, y) for backwards compatibility
func (p *Parser) parseTypedFunctionParameters() ([]*ast.Identifier, []*ast.TypedParameter) {
	identifiers := []*ast.Identifier{}
	typedParams := []*ast.TypedParameter{}
	hasTypes := false

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers, typedParams
	}

	p.nextToken()

	// Parse first parameter
	if !p.curTokenIs(token.IDENT) {
		return nil, nil
	}
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	// Check for type annotation
	var typeAnn *ast.TypeAnnotation
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // consume ':'
		typeAnn = p.parseTypeAnnotation()
		if typeAnn == nil {
			return nil, nil
		}
		hasTypes = true
	}
	typedParams = append(typedParams, &ast.TypedParameter{Name: ident, Type: typeAnn})

	// Parse remaining parameters
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume ','
		p.nextToken() // move to next parameter name

		if !p.curTokenIs(token.IDENT) {
			return nil, nil
		}
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)

		// Check for type annotation
		var typeAnn *ast.TypeAnnotation
		if p.peekTokenIs(token.COLON) {
			p.nextToken() // consume ':'
			typeAnn = p.parseTypeAnnotation()
			if typeAnn == nil {
				return nil, nil
			}
			hasTypes = true
		}
		typedParams = append(typedParams, &ast.TypedParameter{Name: ident, Type: typeAnn})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil, nil
	}

	// Only return typed params if at least one parameter has a type
	if hasTypes {
		return identifiers, typedParams
	}
	return identifiers, nil
}

// parseReturnTypes parses return types after -> in function definition
// Supports single type: -> int, or multiple types: -> int, bool
func (p *Parser) parseReturnTypes() []*ast.TypeAnnotation {
	returnTypes := []*ast.TypeAnnotation{}

	// First return type
	typeAnn := p.parseTypeAnnotation()
	if typeAnn == nil {
		return nil
	}
	returnTypes = append(returnTypes, typeAnn)

	// Check for multiple return types (comma-separated)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume ','
		typeAnn := p.parseTypeAnnotation()
		if typeAnn == nil {
			return nil
		}
		returnTypes = append(returnTypes, typeAnn)
	}

	return returnTypes
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseArrayElements()
	return array
}

// parseArrayElements parses array elements including spread expressions
func (p *Parser) parseArrayElements() []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	bracketToken := p.curToken

	p.nextToken()

	// Check for slice expression: arr[start:end], arr[:end], arr[start:]
	// If current token is COLON, it's [:end]
	if p.curTokenIs(token.COLON) {
		// [:end] case
		sliceExp := &ast.SliceExpression{Token: bracketToken, Left: left, Start: nil}
		p.nextToken()
		if !p.curTokenIs(token.RBRACKET) {
			sliceExp.End = p.parseExpression(LOWEST)
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		}
		return sliceExp
	}

	// Parse the first expression
	firstExpr := p.parseExpression(LOWEST)

	// Check if next token is COLON (slice) or RBRACKET (index)
	if p.peekTokenIs(token.COLON) {
		// Slice expression: arr[start:end] or arr[start:]
		sliceExp := &ast.SliceExpression{Token: bracketToken, Left: left, Start: firstExpr}
		p.nextToken() // consume COLON
		p.nextToken() // move to end expression or RBRACKET

		if !p.curTokenIs(token.RBRACKET) {
			sliceExp.End = p.parseExpression(LOWEST)
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		}
		return sliceExp
	}

	// Regular index expression
	exp := &ast.IndexExpression{Token: bracketToken, Left: left, Index: firstExpr}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

// parseSpreadExpression parses ...expression
func (p *Parser) parseSpreadExpression() ast.Expression {
	spreadToken := p.curToken
	p.nextToken()

	right := p.parseExpression(PREFIX)
	return &ast.SpreadExpression{Token: spreadToken, Right: right}
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseWhileExpression() ast.Expression {
	// while (condition) { body }
	expr := &ast.WhileExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseForExpression() ast.Expression {
	// for item in list { ... }
	// OR
	// for i, v in list { ... }
	// OR
	// for (i=0; i<10; i++) { ... }
	// OR
	// for i in 0..10 { ... }

	if p.peekTokenIs(token.LPAREN) {
		// C-style loop
		return p.parseCForExpression()
	}

	forToken := p.curToken

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	firstIdent := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for for i, v in list { ... }
	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume comma
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		secondIdent := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if !p.expectPeek(token.IN) {
			return nil
		}

		p.nextToken()
		var iterable ast.Expression
		if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.LBRACE) {
			iterable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		} else {
			iterable = p.parseExpression(LOWEST)
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		body := p.parseBlockStatement()

		return &ast.ForInIndexExpression{
			Token:    forToken,
			Index:    firstIdent,
			Value:    secondIdent,
			Iterable: iterable,
			Body:     body,
		}
	}

	// Python-style loop: for item in list { ... }
	expr := &ast.ForExpression{Token: forToken}
	expr.Item = firstIdent

	if !p.expectPeek(token.IN) {
		return nil
	}

	p.nextToken()
	// Parse the iterable, but don't let it consume the opening brace
	// Use a simple identifier parse if it's just an identifier, otherwise parse expression
	if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.LBRACE) {
		// Just an identifier followed by block
		expr.Iterable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	} else {
		expr.Iterable = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseCForExpression() ast.Expression {
	expr := &ast.CForExpression{Token: p.curToken}

	p.nextToken() // consume (

	// Init
	p.nextToken()

	// Init
	expr.Init = p.parseStatement()

	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	p.nextToken()

	expr.Update = p.parseStatement()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseDotExpression(left ast.Expression) ast.Expression {
	// left.right
	// right should be an identifier (method name or field name)

	// Precedence of DOT is high.

	// We need to parse the identifier after dot.
	p.nextToken()
	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Just field access (or method access, call will be handled by LPAREN infix)
	return &ast.InfixExpression{
		Token:    token.Token{Type: token.DOT, Literal: "."},
		Left:     left,
		Operator: ".",
		Right:    name,
	}
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
}

func (p *Parser) parsePrefixIncDec() ast.Expression {
	tok := p.curToken
	p.nextToken()
	operand := p.parseExpression(PREFIX)
	return &ast.PrefixExpression{
		Token:    tok,
		Operator: tok.Literal,
		Right:    operand,
	}
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	stmt := &ast.ContinueStatement{Token: p.curToken}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseSwitchExpression() ast.Expression {
	expr := &ast.SwitchExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expr.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Cases = []*ast.CaseExpression{}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.CASE) {
			caseExpr := &ast.CaseExpression{Token: p.curToken}
			p.nextToken()
			caseExpr.Value = p.parseExpression(LOWEST)

			if !p.expectPeek(token.COLON) {
				return nil
			}

			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			caseExpr.Body = p.parseBlockStatement()
			expr.Cases = append(expr.Cases, caseExpr)
		} else if p.curTokenIs(token.DEFAULT) {
			if !p.expectPeek(token.COLON) {
				return nil
			}
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			expr.Default = p.parseBlockStatement()
		}
		p.nextToken()
	}

	return expr
}

func (p *Parser) parseTernaryExpression(condition ast.Expression) ast.Expression {
	expr := &ast.TernaryExpression{
		Token:     p.curToken,
		Condition: condition,
	}

	p.nextToken()
	expr.Consequence = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	expr.Alternative = p.parseExpression(TERNARY)

	return expr
}

func (p *Parser) parseRangeExpression(start ast.Expression) ast.Expression {
	expr := &ast.RangeExpression{
		Token: p.curToken,
		Start: start,
	}

	p.nextToken()
	expr.End = p.parseExpression(RANGE_PREC)

	return expr
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("unexpected token '%s'", t)
	p.errors = append(p.errors, msg)

	loc := errors.SourceLocation{
		Line:      p.curToken.Line,
		Column:    p.curToken.Column,
		EndColumn: p.curToken.EndColumn,
		Filename:  p.filename,
	}

	var help string
	switch t {
	case token.RBRACE:
		help = "you might have an extra closing brace '}'"
	case token.RPAREN:
		help = "you might have an extra closing parenthesis ')'"
	case token.RBRACKET:
		help = "you might have an extra closing bracket ']'"
	case token.ASSIGN:
		help = "did you forget to declare a variable with 'let'?"
	default:
		help = "check your syntax around this location"
	}

	richErr := errors.ParseError(fmt.Sprintf("unexpected token '%s'", t), loc, p.sourceCode).WithHelp(help)
	p.richErrors = append(p.richErrors, richErr)
}

// parseArrowFunction parses lambda shorthand: x => x * 2 or (x, y) => x + y
func (p *Parser) parseArrowFunction(left ast.Expression) ast.Expression {
	arrowToken := p.curToken
	var params []*ast.Identifier

	switch l := left.(type) {
	case *ast.Identifier:
		params = []*ast.Identifier{l}
	case *ast.InfixExpression:
		return nil
	default:
		msg := fmt.Sprintf("expected identifier before '=>', got %T", left)
		p.errors = append(p.errors, msg)
		return nil
	}

	p.nextToken() // consume =>

	body := p.parseExpression(ARROW_PREC)

	return &ast.ArrowFunction{
		Token:      arrowToken,
		Parameters: params,
		Body:       body,
	}
}
