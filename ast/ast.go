package ast

import (
	"bytes"
	"go-compiler/main/token"
	"strings"
)

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

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token *token.Token // token.Let
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Lexeme
}
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type AssignStatement struct {
	Token    token.Token
	Name     *Identifier
	Value    Expression
	EnvIndex int
	EnvDepth int
}

func (as *AssignStatement) expressionNode()      {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Lexeme }
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString("(" + as.Name.Value + " = " + as.Value.String() + ")")
	return out.String()
}

type Identifier struct {
	Token *token.Token // token.IDENTIFIER
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Lexeme
}
func (i *Identifier) String() string {
	return i.Value
}

type NumberLiteral struct {
	Token *token.Token // token.NUMBER
	Value float64
}

func (i *NumberLiteral) expressionNode() {}
func (i *NumberLiteral) TokenLiteral() string {
	return i.Token.Lexeme
}
func (i *NumberLiteral) String() string {
	return i.Token.Lexeme
}

type StringLiteral struct {
	Token *token.Token // token.NUMBER
	Value string
}

func (i *StringLiteral) expressionNode() {}
func (i *StringLiteral) TokenLiteral() string {
	return i.Token.Lexeme
}
func (i *StringLiteral) String() string {
	return i.Token.Lexeme
}

type ReturnStatement struct {
	Token       *token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Lexeme
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      *token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Lexeme
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type PrefixExpression struct {
	Token    *token.Token // - or !, can add more later
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Lexeme
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + pe.Operator + pe.Right.String() + ")")

	return out.String()
}

type InfixExpression struct {
	Token    *token.Token // - or !, can add more later
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Lexeme
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")")

	return out.String()
}

type Boolean struct {
	Token *token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Lexeme
}
func (b *Boolean) String() string {
	return b.Token.Lexeme
}

type IfExpression struct {
	Token     *token.Token // if token
	Condition Expression
	Then      *BlockStatment
	Else      *BlockStatment
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Lexeme
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if" + ie.Condition.String() + " " + ie.Then.String())
	if ie.Else != nil {
		out.WriteString("else " + ie.Else.String())
	}
	return ie.Token.Lexeme
}

type WhileExpression struct {
	Token     *token.Token // if token
	Condition Expression
	Body      *BlockStatment
}

func (we *WhileExpression) expressionNode() {}
func (we *WhileExpression) TokenLiteral() string {
	return we.Token.Lexeme
}
func (we *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while" + we.Condition.String() + " " + we.Body.String())
	return we.Token.Lexeme
}

type BlockStatment struct {
	Token      token.Token // { token
	Statements []Statement
}

func (bs *BlockStatment) statementNode()       {}
func (bs *BlockStatment) TokenLiteral() string { return bs.Token.Lexeme }
func (bs *BlockStatment) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      *token.Token // The Function token
	Parameters []*Identifier
	Body       *BlockStatment
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Lexeme
}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral() + "(" + strings.Join(params, ",") + ")" + fl.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     *token.Token // '(' Token
	Function  Expression
	Arguments []Expression // Identifier or FunctionLiteral
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Lexeme
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String() + "(" + strings.Join(args, ", ") + ")")
	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token // '[' Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Lexeme }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, expr := range al.Elements {
		elements = append(elements, expr.String())
	}

	out.WriteString("[" + strings.Join(elements, ", ") + "]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Lexeme }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(" + ie.Left.String() + "[" + ie.Index.String() + "])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Lexeme }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{" + strings.Join(pairs, ", ") + "}")
	return out.String()
}
