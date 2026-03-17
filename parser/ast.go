package parser

import (
	"github.com/hasnainzeenwala/hzglox/lexer"
)


// -- Base Production Rule --
// expression -> literal | unary | binary | grouping
type Expr interface {
	IsExpr()
	PrintAst() string
}

// -- Literal -> NUMBER | STRING | true | false | nil ; --
type Literal struct {
	T lexer.Token
}

func (l Literal) IsExpr() {}

func (l Literal) PrintAst() string {
	return l.T.Lexeme
}

func NewLiteral(T lexer.Token) (l Literal, match bool){
	match = IsLiteral(T)
	if match {
		l = Literal{T}
	}
	return
}


// -- Grouping -> "(" expression ")" ; --
type Grouping struct{
	E Expr
}

func (g Grouping) IsExpr() {}

func (g Grouping) PrintAst() string {
	return "( group " + g.E.PrintAst() + " )"
}

func NewGrouping(Lt lexer.Token, E Expr, Rt lexer.Token) (g Grouping, match bool){
	match = IsLeftParen(Lt) && E != nil && IsRightParen(Rt)
	if match {
		g = Grouping{E}
	}
	return
}

// -- Unary -> ( "-" | "!" ) expression ; --
type Unary struct{
	T lexer.Token
	E Expr
}
func (u Unary) IsExpr() {}

func (u Unary) PrintAst() string {
	return "( " + u.T.Lexeme + u.E.PrintAst() + " )"
}

func NewUnary(T lexer.Token, E Expr) (u Unary, match bool) {
	match = (IsMinus(T) || IsBang(T)) && E != nil
	if match {
		u = Unary{T: T, E: E}
	}
	return
}

// -- Binary -> expression op expression ; --
type Binary struct {
	Le Expr
	Op lexer.Token
	Re Expr
}

func (b Binary) IsExpr() {}

func (b Binary) PrintAst() string {
	return "( " + b.Op.Lexeme + " " + b.Le.PrintAst() + " " + b.Re.PrintAst() + " )"
}

func NewBinary(Le Expr, Op lexer.Token, Re Expr) (b Binary, match bool) {
	match = Le != nil && IsOperator(Op) && Re != nil
	if match {
		b = Binary{Le: Le, Op: Op, Re: Re}
	}
	return
}


// -- Helper functions to Establish the Token Type --

func IsLiteral(t lexer.Token) bool {
	return t.TType == lexer.Number || t.TType == lexer.String ||
		t.TType == lexer.True || t.TType == lexer.False || t.TType == lexer.Nil
}

func IsLeftParen(t lexer.Token) bool {
	return t.TType == lexer.LeftParen
}

func IsRightParen(t lexer.Token) bool {
	return t.TType == lexer.RightParen
}

func IsMinus(t lexer.Token) bool {
	return t.TType == lexer.Minus
}

func IsBang(t lexer.Token) bool {
	return t.TType == lexer.Bang
}

func IsOperator(t lexer.Token) bool {
	return t.TType == lexer.EqualEqual || t.TType == lexer.BangEqual ||
		t.TType == lexer.Less || t.TType == lexer.LessEqual || t.TType == lexer.Greater ||
		t.TType == lexer.GreaterEqual || t.TType == lexer.Star || t.TType == lexer.Minus ||
		t.TType == lexer.Plus || t.TType == lexer.Slash
}
