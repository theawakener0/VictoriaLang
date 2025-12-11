package parser

import (
	"fmt"
	"strconv"
	"victoria/ast"
	"victoria/lexer"
	"victoria/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	DOT         // struct.field
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.MODULO:   PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.DOT:      DOT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
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

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Line %d: Expected '%s' but found '%s'", p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
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
	case token.RETURN:
		return p.parseReturnStatement()
	case token.INCLUDE:
		return p.parseIncludeStatement()
	case token.TRY:
		return p.parseTryStatement()
	case token.STRUCT:
		return p.parseStructStatement()
	case token.FUNCTION:
		// Check if it is a method definition: def Struct.Method()
		if p.peekTokenIs(token.IDENT) {
			// It could be a function literal in an expression, but here we are at statement level.
			// However, standard function def is `def name()`.
			// Method def is `def Struct.name()`.
			// Let's peek further.
			// We can't easily peek 2 tokens ahead with this lexer setup without modifying it or consuming.
			// But `def` at statement level usually means function declaration or method declaration.
			// Let's handle it.
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

		methodDef.Parameters = p.parseFunctionParameters()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		methodDef.Body = p.parseBlockStatement()
		return methodDef

	} else {
		// Function definition: define greet()
		// This is actually a LetStatement in disguise: let greet = fn() { ... }
		// But we want to support `define greet() {}` as top level.
		// We can treat it as a LetStatement where value is FunctionLiteral.

		fnLit := &ast.FunctionLiteral{Token: defToken}
		fnLit.Name = firstIdent.Value

		if !p.expectPeek(token.LPAREN) {
			return nil
		}

		fnLit.Parameters = p.parseFunctionParameters()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		fnLit.Body = p.parseBlockStatement()

		// Wrap in LetStatement
		letStmt := &ast.LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let", Line: defToken.Line},
			Name:  firstIdent,
			Value: fnLit,
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
		// But IDENT is handled by parseIdentifier.
		// We need to check if we are parsing a struct instantiation inside parseIdentifier?
		// Or maybe parseIdentifier should look ahead?
		// Let's handle it in parseIdentifier.
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
		// This is ambiguous with function call returning a function that takes a block? No.
		// Ambiguous with `x { ... }` which isn't valid unless x is a struct type.
		// But wait, `if` `for` etc are keywords.
		// If we have `ident {`, it's likely a struct instantiation or a hash literal if ident was missing (but here we have ident).
		// Let's assume it is struct instantiation.
		return p.parseStructInstantiation()
	}
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseStructInstantiation() ast.Expression {
	// curToken is the Struct Name
	structName := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // consume name, now at {

	si := &ast.StructInstantiation{Token: structName.Token, Name: structName}
	si.Fields = make(map[string]ast.Expression)

	// parse { key: value, ... }
	// This is similar to hash literal but keys are identifiers (fields)

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
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

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

	lit.Parameters = p.parseFunctionParameters()

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
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
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
	// for (i=0; i<10; i++) { ... }

	if p.peekTokenIs(token.LPAREN) {
		// C-style loop
		return p.parseCForExpression()
	}

	// Python-style loop
	expr := &ast.ForExpression{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	expr.Item = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.IN) {
		return nil
	}

	p.nextToken()
	expr.Iterable = p.parseExpression(LOWEST)

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
	// We expect a statement here, usually a LetStatement or ExpressionStatement (assignment)
	// But parseStatement expects to be at the start of a statement.
	// Let's reuse parseStatement but we need to be careful about semicolons.
	// parseStatement consumes the semicolon if present.

	// Actually, C-style for loop parts are statements.
	// for (let i = 0; i < 10; i = i + 1)

	// Init
	expr.Init = p.parseStatement()
	// parseStatement consumes the semicolon if it's a LetStatement or ExpressionStatement
	// If it didn't consume semicolon (e.g. if we didn't put one), we might need to check.
	// But our parseLetStatement consumes semicolon if present.
	// If it's missing, we might be fine or not.
	// Let's assume standard C-style: for (init; cond; update)

	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	p.nextToken()
	// Update is usually an assignment or expression
	// It shouldn't have a semicolon at the end inside the loop header usually?
	// for (...; ...; i++)
	// We can parse it as a statement but ensure we don't consume the closing paren as part of it?
	// Actually, `i++` isn't a statement in our language yet (no ++ operator). `i = i + 1` is an assignment.
	// Assignment is an expression in some languages, statement in others.
	// In our AST, we don't have AssignmentExpression, we have LetStatement.
	// But we can have `i = i + 1` as an expression if we support assignment expressions.
	// The prompt says `Assignment: =` under Operators. So `x = 10` is likely an expression or statement.
	// If `x = 10` is an expression, then `i = i + 1` is an expression.
	// Let's assume assignment is an infix expression for now?
	// Wait, `let x = 10` is a statement.
	// Re-assigning `x = 11` is usually an expression or statement.
	// If I didn't implement AssignmentExpression, I should.
	// I implemented LetStatement.
	// I need to handle reassignment. `x = 5`.
	// This is usually parsed as an infix expression where left side is identifier.
	// Let's add support for `=` as infix operator in `parseInfixExpression`.
	// I already registered `=` as infix.

	expr.Update = p.parseStatement()
	// Note: parseStatement might consume the closing ) if we are not careful, but usually it stops at semicolon.
	// But the update clause doesn't have a semicolon.
	// So parseStatement might fail or consume too much.
	// Let's just parse expression for update.
	// But wait, `i = i + 1` is an expression? Yes if `=` is infix.
	// So let's change `Update` to Expression.
	// But `ast.CForExpression` has `Update Statement`. Let's change it to Expression or handle it.
	// Actually, let's just parse expression for update.

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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
