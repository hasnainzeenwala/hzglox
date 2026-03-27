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

func (l *Literal) IsExpr() {}

func (l *Literal) PrintAst() string {
	return l.T.Lexeme
}

// -- Grouping -> "(" expression ")" ; --
type Grouping struct{
	E Expr
}

func (g *Grouping) IsExpr() {}

func (g *Grouping) PrintAst() string {
	return "( group " + g.E.PrintAst() + " )"
}

// -- Unary -> ( "-" | "!" ) expression ; --
type Unary struct{
	T lexer.Token
	E Expr
}
func (u *Unary) IsExpr() {}

func (u *Unary) PrintAst() string {
	return "( " + u.T.Lexeme + u.E.PrintAst() + " )"
}

// -- Binary -> expression op expression ; --
type Binary struct {
	Le Expr
	Op lexer.Token
	Re Expr
}

func (b *Binary) IsExpr() {}

func (b *Binary) PrintAst() string {
	return "( " + b.Op.Lexeme + " " + b.Le.PrintAst() + " " + b.Re.PrintAst() + " )"
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
