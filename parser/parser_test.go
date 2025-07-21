package parser

import (
	"testing"
	"monkey/ast"
	"monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input  := `
	let x 5;
	let = 10;
	let 838383;`

	lexer := lexer.New(input)
	parser := New(lexer) //in the same package, so calling from parser

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	// if program == nil {
	// 	t.Fatalf("ParseProgram returned nil")
	// }

	// if len(program.Statements) != 3 {
	// 	t.Fatalf("program.Statements does not contain 3 statements, %d instead", len(program.Statements))
	// }

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier){
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral Not let, got %q", s.TokenLiteral())
		return false
	}

	letStatement, ok := s.(*ast.LetStatement) // type assertion, checking if s can be treated as LetStatement
	if !ok {
		t.Errorf("s not *ast.LetStatement, got %T", s)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not %s, got %s", name, letStatement.Name.Value)
		return false
	}
	
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral not %s, got %s", name, letStatement.Name)
		return false
	}

	return true

}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)

	}
	t.FailNow()
}