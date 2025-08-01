package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer
	// these are the same concept as in the lexer, but instead of pointing to character, they point to 
	// current and next token
	curToken token.Token // like token.INT, EOF, etc
	peekToken token.Token
	errors []string
}


func New(l *lexer.Lexer) *Parser {
	// since p holds the address of the parser, it is a pointer
	p := &Parser{l: l,
				errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
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
		return nil
	}
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