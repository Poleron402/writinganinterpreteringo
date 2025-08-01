package lexer

import (
	"monkey/token"
)

// keep in mind that all the struct variables are initialized to their respective 0 values
type Lexer struct {
	input string // the string we are working on
	position int // points to the current char in input
	readPosition int //current reading position is after the current char (because we will need to see what comes after the char that is read)
	ch byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input} 
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // signifies either 'we havent read anything yet, or eof
	}else{
		l.ch = l.input[l.readPosition] // when this is called for the first time, readPosition is automatically 0 bc of how structs be
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType token.TokenType, ch byte) token.Token{
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch{ 
	case '=':
		if l.peekChar() == '=' {
			char := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(char)+string(l.ch)}
		}else{
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '='{
			char := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(char)+string(l.ch)}
		}else{
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch){
			tok.Literal = l.readIdentifier() // returns a whole word
			tok.Type = token.LookupIdent(tok.Literal) // check if word is def or let
			return tok // we have to return because we already readchar in readIdentifier
		}else if isDigit(l.ch){
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok // we have to return because we already readchar in readNumber
		}else{
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' '|| l.ch == '\t'|| l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch){
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch){
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte)bool{
	return '0'<=ch && ch<='9'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}