package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

// bunch of constants representing token classifiers/types
const (
	ILLEGAL = "ILLEGAL" // A TOKEN/CHARACTER WE DO NOT KNOW ABOUT
	EOF = "EOF"
	IDENT = "IDENT"
	INT = "INT"

	// operators
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"
	LT = "<"
	GT = ">"

	COMMA ="COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET  = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"

	EQ = "=="
	NOT_EQ = "!="
)

var keywords = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if value, exists := keywords[ident]; exists{
		return value
	}
	return IDENT // this will be the token type for all user-defined identifiers
}