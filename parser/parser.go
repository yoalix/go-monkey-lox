package parser

import (
	"fmt"
	"go-compiler/main/ast"
	"go-compiler/main/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	AND         // and OR or
	EQUALS      // ==
	LESSGREATER // > OR <
	SUM         // + OR -
	PRODUCT     // * OR /
	PREFIX      // -X OR !X
	CALL        // myFunction(X)
	INDEX
)

var precedences = map[token.TokenType]int{
	token.AND:          AND,
	token.OR:           AND,
	token.EQUAL_EQUAL:  EQUALS,
	token.BANG_EQUAL:   EQUALS,
	token.GREATER:      LESSGREATER,
	token.LESS:         LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.STAR:         PRODUCT,
	token.LEFT_PAREN:   CALL,
	token.LEFT_BRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	tokens         []*token.Token
	current        int
	currToken      *token.Token
	peekToken      *token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(tokens []*token.Token) *Parser {
	p := &Parser{tokens: tokens, current: 0, errors: []string{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LEFT_PAREN, p.parseGroupedExpressions)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	// p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LEFT_BRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LEFT_BRACE, p.parseHashLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.STAR, p.parseInfixExpression)
	p.registerInfix(token.EQUAL_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.BANG_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESS, p.parseInfixExpression)
	p.registerInfix(token.GREATER, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LEFT_PAREN, p.parseCallExpression)
	p.registerInfix(token.LEFT_BRACKET, p.parseIndexExpression)
	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	if p.current < len(p.tokens) {
		p.peekToken = p.tokens[p.current]
	} else {
		p.peekToken = nil
	}
	p.current++
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		statement := p.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()

	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Lexeme}

	if !p.expectPeek(token.EQUAL) {
		return nil
	}
	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currToken}
	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.currToken}
	statement.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecendence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currToken}
	if !p.expectPeek(token.LEFT_PAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}
	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}
	expression.Then = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LEFT_BRACE) {
			return nil
		}
		expression.Else = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.currToken}
	if !p.expectPeek(token.LEFT_PAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}
	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}
	expression.Body = p.parseBlockStatement()
	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatment {
	block := &ast.BlockStatment{Token: *p.currToken}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.currTokenIs(token.RIGHT_BRACE) && !p.currTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}
	return block

}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	expression := &ast.FunctionLiteral{Token: p.currToken}
	// consume fn
	if !p.expectPeek(token.LEFT_PAREN) {
		return nil
	}
	expression.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}
	expression.Body = p.parseBlockStatement()
	return expression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RIGHT_PAREN) {
		p.nextToken()
		return identifiers
	}
	// consume (
	p.nextToken()
	identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Lexeme}
	identifiers = append(identifiers, identifier)

	for p.peekTokenIs(token.COMMA) {
		// consume ,
		p.nextToken()
		// consume y
		p.nextToken()
		identifier = &ast.Identifier{Token: p.currToken, Value: p.currToken.Lexeme}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currToken, Function: function}
	expression.Arguments = p.parseExpressionList(token.RIGHT_PAREN)
	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Lexeme}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	f, err := strconv.ParseFloat(p.currToken.Lexeme, 64)
	if err != nil {
		message := token.TokenError(p.peekToken, fmt.Sprintf("could not parse %q as integer", p.currToken.Lexeme))
		p.errors = append(p.errors, message)
		return nil
	}
	return &ast.NumberLiteral{Token: p.currToken, Value: f}
}
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Lexeme}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpressions() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Lexeme}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{Token: p.currToken, Left: left, Operator: p.currToken.Lexeme}
	precedence := p.currPrecendence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: *p.currToken}
	arr.Elements = p.parseExpressionList(token.RIGHT_BRACKET)
	return arr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: *p.currToken, Left: left}
	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RIGHT_BRACKET) {
		return nil
	}
	return expression
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.
		parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		// consume , comma
		p.nextToken()
		// consume ' ' space
		p.nextToken()
		list = append(list, p.
			parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: *p.currToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)
	if p.peekTokenIs(token.RIGHT_BRACE) {
		p.nextToken()
		return hash
	}
	for !p.peekTokenIs(token.RIGHT_BRACE) {
		// consume "{" first time or " " following times
		p.nextToken()
		key := p.parseExpression(LOWEST)
		// consume ':'
		if !p.expectPeek(token.COLON) {
			return nil
		}
		// consume ' '
		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		if !p.peekTokenIs(token.RIGHT_BRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	// // consume "{" left brace
	// p.nextToken()
	// key := p.parseExpression(LOWEST)
	// // consume ':'
	// if !p.expectPeek(token.COLON) {
	// 	return nil
	// }
	// // consume ' '
	// p.nextToken()
	// value := p.parseExpression(LOWEST)
	// hash.Pairs[key] = value
	// for p.peekTokenIs(token.COMMA) {
	// 	// consume , comma
	// 	p.nextToken()
	// 	// consume ' ' space
	// 	p.nextToken()
	// 	key = p.parseExpression(LOWEST)
	// 	// consume ':'
	// 	if !p.expectPeek(token.COLON) {
	// 		return nil
	// 	}
	// 	// consume ' '
	// 	p.nextToken()
	// 	value = p.parseExpression(LOWEST)
	// 	hash.Pairs[key] = value
	// }
	if !p.expectPeek(token.RIGHT_BRACE) {
		return nil
	}
	return hash
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekPrecendence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecendence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekError(t token.TokenType) *ParseError {
	errorMessage := token.TokenError(p.peekToken, fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type))
	p.errors = append(p.errors, errorMessage)
	return NewParseError()
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	errorMessage := token.TokenError(p.currToken, fmt.Sprintf("no prefix parse function for %s found",
		t))
	p.errors = append(p.errors, errorMessage)
}

func (p *Parser) Errors() []string {
	return p.errors
}

type ParseError struct {
}

func NewParseError() *ParseError {
	return &ParseError{}
}
