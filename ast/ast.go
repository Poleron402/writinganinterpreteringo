package ast
import "monkey/token"


// every node is to implement node interface by providing tokenliteral method
// this method is used for debugging and testing, and it returns the literal value associated witht the token


// Node here is like a base interface
type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// implementation of Node

type Program struct {
	Statements []Statement
}

// we are retrieving the first one to see if it is a let or an identifier. Each have their
// own implementations of TokenLiteral method defined below (NOT RECURSIVE is what im tryna say)
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}else{
		return ""
	}
}

type LetStatement struct {
	Token token.Token //this is token.LET
	Name *Identifier
	Value Expression
}

// ls *LetStatement is a receiver (aka the struct the method belongs to)
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {return ls.Token.Literal}


type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {return rs.Token.Literal}
type Identifier struct {
	Token token.Token // this is token.IDENT
	Value string
}

func (i *Identifier) expressionNode() {} // implements expression, by using its method
func (i *Identifier) TokenLiteral() string {return i.Token.Literal}