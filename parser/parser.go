package parser

import (
	"fmt"

	"github.com/hasnainzeenwala/hzglox/lexer"
)

// Lox expression grammar
// Unambiguous and left recursion removed.
// **********************************************************************************
//
// expression   ->  equal
// equal        ->  comparison (("==" | "!=") comparison)*
// comparison   ->  term ( ("<" | "<=" | ">" | ">=") term )*
// term         ->  factor (( "+" | "-" ) factor)*
// factor       ->  unary (("*" | "/") unary)*
// unary        ->  ("-" | "!") unary | primary
// primary      ->  NUMBER | STRING | true | false | nil | "(" expression ")" ;
//
// **********************************************************************************

type Parser struct {
	l *lexer.Lexer
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{l}
}

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return expr, err
	}
	pt, e := p.l.Peek()
	if e != nil {
		return expr, e
	}
	if pt.TType != lexer.Eof {
		return expr, fmt.Errorf("Encountered an unparsable lexeme (%v)", pt)
	}
	return expr, err
}

func (p *Parser) expression() (Expr, error) {
	return p.equal()
}

func (p *Parser) equal() (Expr, error) {
	ops := []lexer.Token{
		{
			TType: lexer.EqualEqual,
		},
		{
			TType: lexer.BangEqual,
		},
	}
	return p.binFOpfRepeat(p.comparison, ops)
}

func (p *Parser) comparison() (Expr, error) {
	ops := []lexer.Token{
		{
			TType: lexer.Less,
		},
		{
			TType: lexer.LessEqual,
		},
		{
			TType: lexer.Greater,
		},
		{
			TType: lexer.GreaterEqual,
		},
	}
	return p.binFOpfRepeat(p.term, ops)
}

func (p *Parser) term() (Expr, error) {
	ops := []lexer.Token{
		{
			TType: lexer.Plus,
		},
		{
			TType: lexer.Minus,
		},
	}
	return p.binFOpfRepeat(p.factor, ops)
}

func (p *Parser) factor() (Expr, error) {
	ops := []lexer.Token{
		{
			TType: lexer.Star,
		},
		{
			TType: lexer.Slash,
		},
	}
	return p.binFOpfRepeat(p.unary, ops)
}

func (p *Parser) unary() (Expr, error) {
	pl, err := p.l.Peek()
	if err != nil {
		return nil, err
	}
	switch pl.TType {
	case lexer.Bang, lexer.Minus:
		t, err := p.l.FetchNextToken()
		if err != nil {
			return nil, err
		}
		r, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{
			T: t,
			E: r,
		}, nil
	default:
		return p.primary()
	}
}

func (p *Parser) primary() (Expr, error) {
	t, err := p.l.FetchNextToken()
	if err != nil {
		return nil, err
	}
	switch t.TType {
	case lexer.Number, lexer.String, lexer.Nil, lexer.True, lexer.False:
		return &Literal{
			T: t,
		}, nil

	case lexer.LeftParen:
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}
		t, err = p.l.FetchNextToken()
		if t.TType != lexer.RightParen {
			return nil, fmt.Errorf("Expected Right Paren ')' at the end of the expression but found %s instead", t)
		}
		return &Grouping{exp}, nil
	default:
		return nil, fmt.Errorf("Token not recognized %s", t)
	}
}

// General parsing function for all rules of type
// *********************************************************
// NT -> X ((op1 | op2 ...) X)*
// *********************************************************
// X can be a terminal or non terminal
func (p *Parser) binFOpfRepeat(f func() (Expr, error), op []lexer.Token) (Expr, error) {
	left, err := f()
	if err != nil {
		return nil, err
	}
	for {
		pl, err := p.l.Peek()
		if err != nil {
			return left, err
		}

		// If token type doesn't match we are done and can return the parsed expression thusfar
		if !tokenTypeMatches(pl, op) {
			return left, nil
		}

		// Else we create a binary expression
		op, err := p.l.FetchNextToken()
		if err != nil {
			return left, err
		}
		r, err := f()
		if err != nil {
			return left, err
		}

		exp := &Binary{
			Op: op,
			Le: left,
			Re: r,
		}

		left = exp
	}
}

func tokenTypeMatches(pl lexer.Token, l []lexer.Token) bool {
	for _, t := range l {
		if pl.TType == t.TType {
			return true
		}
	}
	return false
}
