package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)
var precedences = map[token.TokenType]int{
	token.EQ:EQUALS,
	token.NOT_EQ:EQUALS,
	token.LT:LESSGREATER,
	token.GT:LESSGREATER,
	token.PLUS:SUM,
	token.MINUS:SUM,
	token.SLASH:PRODUCT,
	token.ASTERISK: PRODUCT,
}
const (
	_ int = iota
	LOWEST
	EQUALS // ==
	LESSGREATER // > or <
	SUM// +
	PRODUCT// *
	PREFIX// -X or !X
	CALL// myFunction(X)
)

type (
	prefixParseFn func() ast.Expression
	inflixParseFn func(ast.Expression) ast.Expression 
)
type Parser struct {
	l *lexer.Lexer
	// these are the same concept as in the lexer, but instead of pointing to character, they point to 
	// current and next token
	curToken token.Token // like token.INT, EOF, etc
	peekToken token.Token
	errors []string

	prefixParseFnMap map[token.TokenType]prefixParseFn
	infixParseFnMap map[token.TokenType]inflixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFnMap[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn inflixParseFn) {
	p.infixParseFnMap[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	// since p holds the address of the parser, it is a pointer
	p := &Parser{l: l,
				errors: []string{}}
	p.prefixParseFnMap = make(map[token.TokenType]prefixParseFn)
	// registering prefixes
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpressions)
	p.registerPrefix(token.MINUS, p.parsePrefixExpressions)	
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	// registering infixes
	p.infixParseFnMap = make(map[token.TokenType]inflixParseFn)	
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	
	p.nextToken()
	p.nextToken()
	return p
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
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
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
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN){
		return nil
	}
	return exp
}
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression {
		Token: p.curToken,
		Operator: p.curToken.Literal,
		Left: left,
	}
	precendence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precendence)
	return expression
}
func(p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken= p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}


func (p *Parser) Errors() []string{
	return p.errors
}
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected the next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	// consoom semicolon so that it doesn't affect the next statement
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {

	// constructing an initial node with the current token (the token.LET)
	stmt := &ast.LetStatement{Token: p.curToken}

	// Now the method advances the tokens while asserting next ones

	// since let is usually followed by an identifier of some kind, that's what 
	// we are expecting
	if !p.expectPeek(token.IDENT){
		return nil
	}
	// if it is an identifier, then we construct an identifier node
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	//expecting an assignment
	if !p.expectPeek(token.ASSIGN){
		return nil
	}
	// progressing through the expression until we see a semicolon
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	new_node := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return new_node
}
func (p *Parser) curTokenIs(t token.TokenType) bool{
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool{
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t){
		p.nextToken()
		return true
	}else{
		p.peekError(t)
		return false
	}
}


func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}


func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFnMap[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence <p.peekPrecedence() {
		infix := p.infixParseFnMap[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parsePrefixExpressions() ast.Expression {
	expression := &ast.PrefixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
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

