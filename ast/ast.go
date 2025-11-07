package ast // abstract syntax tree

import (
	"bytes"
	"monkey/token"
)

// every node is to implement node interface by providing tokenliteral method
// this method is used for debugging and testing, and it returns the literal value associated witht the token

// Node here is like a base interface
type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// we are retrieving the first one to see if it is a let or an identifier. Each have their
// own implementations of TokenLiteral method defined below (NOT RECURSIVE is what im tryna say)
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// ========== Statements ==========
type LetStatement struct {
	Token token.Token //this is token.LET
	Name  *Identifier
	Value Expression
}

// ls *LetStatement is a receiver (aka the struct the method belongs to)
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type Identifier struct {
	Token token.Token // this is token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {} // implements expression, by using its method
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// ========== Expressions ==========

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }

// expression strings 2.6

func (l *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.Value)

	if l.Value != nil {
		out.WriteString(" = " + l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(" = " + rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (idnt *Identifier) String() string {
	return idnt.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}



type InfixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}
func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {return ie.Token.Literal}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" "+ie.Operator+" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// Boolean Literals (2.8)

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {return b.Token.Literal}
func (b *Boolean) String() string {return b.Token.Literal}

// ============ IF EXPRESSIONS ============

type IfExpression struct {
	Token token.Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {return ie.Token.Literal}
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {return bs.Token.Literal}
func (bs * BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}