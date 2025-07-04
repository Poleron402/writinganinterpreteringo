package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // A TOKEN/CHARACTER WE DO NOT KNOW ABOUT
	EOF = "EOF"
	IDENT = "IDENT"
	INT = "INT"

	ASSIGN = "="
	PLUS = "+"

	COMMA ="COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET  = "LET"
)