package parser

import (
	"fmt"

	"github.com/hasnainzeenwala/hzglox/lexer"
)

// ============================================================================================
// Parser Description
// ============================================================================================

// Lox expression grammar
// Unambiguous and left recursion removed.
// "expression" is the starting rule of the grammar
// **********************************************************************************
//
// -> expression   ->  equal
//    equal        ->  comparison (("==" | "!=") comparison)*
//    comparison   ->  term ( ("<" | "<=" | ">" | ">=") term )*
//    term         ->  factor (( "+" | "-" ) factor)*
//    factor       ->  unary (("*" | "/") unary)*
//    unary        ->  ("-" | "!") unary | primary
//    primary      ->  NUMBER | STRING | true | false | nil | "(" expression ")" ;
//
// **********************************************************************************
// Parser is a recursive descent parser
// Every non-terminal/rule shall have its own parsing function. It will be named 'parse<rulename>rule()'
// The recipe to create the parsing function is the following. Read the rule from left to right.
// If you encounter a non-terminal, call the corresponding function for that non-terminal,
// if you encounter a terminal parse it appropriately. The rule might have '|' which is the OR operator,
// so deal with the terminals accordingly. Keep proceeding by matching the terminals and calling the functions for non-terminals
// in a similar fashion till the rule is done.
// For some of the rules, a generic function has been created since all of them had a very similar structure.
// But the generic function follows the same idea described above.




type Parser struct {
	l *lexer.Lexer
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{l}
}

// Main public method which parses and returns the tree
func (p *Parser) Parse() (Node, error) {
	// begin parsing from the topmost rule
	expr, err := p.parseExpressionRule()
	if err != nil {
		return expr, err
	}

	// Check what are the remaining tokens
	// Expectation is that only EOF should be remaining
	t, err := p.l.FetchNextToken()
	if err != nil {
		return expr, fmt.Errorf("Encountered unexpected error while checking if any tokens are left for parsing (%w)", err)
	}
	if t.TType != lexer.Eof {
		return expr, fmt.Errorf("Parsing failed: Encountered an unparseable token (%v)", t)
	}
	return expr, nil
}


// =============================================================================
// Parsing Rules Implementation
// =============================================================================

func (p *Parser) parseExpressionRule() (Node, error) {
	return p.parseEqualRule()
}

func (p *Parser) parseEqualRule() (Node, error) {
	ops := []lexer.Token{
		{
			TType: lexer.EqualEqual,
		},
		{
			TType: lexer.BangEqual,
		},
	}
	return p.genericParseFunctionForRuleOfTypeXopXRepeat(p.parseComparisonRule, ops)
}

func (p *Parser) parseComparisonRule() (Node, error) {
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
	return p.genericParseFunctionForRuleOfTypeXopXRepeat(p.parseTermRule, ops)
}

func (p *Parser) parseTermRule() (Node, error) {
	ops := []lexer.Token{
		{
			TType: lexer.Plus,
		},
		{
			TType: lexer.Minus,
		},
	}
	return p.genericParseFunctionForRuleOfTypeXopXRepeat(p.parseFactorRule, ops)
}

func (p *Parser) parseFactorRule() (Node, error) {
	ops := []lexer.Token{
		{
			TType: lexer.Star,
		},
		{
			TType: lexer.Slash,
		},
	}
	return p.genericParseFunctionForRuleOfTypeXopXRepeat(p.parseUnaryRule, ops)
}

func (p *Parser) parseUnaryRule() (Node, error) {
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
		r, err := p.parseUnaryRule()
		if err != nil {
			return nil, err
		}
		return &UnaryNode{
			T: t,
			E: r,
		}, nil
	default:
		return p.parsePrimaryRule()
	}
}

func (p *Parser) parsePrimaryRule() (Node, error) {
	t, err := p.l.FetchNextToken()
	if err != nil {
		return nil, err
	}
	switch t.TType {
	case lexer.Number, lexer.String, lexer.Nil, lexer.True, lexer.False:
		return &LiteralNode{
			T: t,
		}, nil

	case lexer.LeftParen:
		exp, err := p.parseExpressionRule()
		if err != nil {
			return nil, err
		}
		t, err = p.l.FetchNextToken()
		if t.TType != lexer.RightParen {
			return nil, fmt.Errorf("Expected Right Paren ')' at the end of the expression but found %s instead", t)
		}
		return &GroupingNode{exp}, nil
	default:
		return nil, fmt.Errorf("Token not recognized %s", t)
	}
}





// ========================================================================
// Helper functions
// ========================================================================

// General parsing function for all rules of type
// *********************************************************
// NT -> X ((op1 | op2 ...) X)*
// *********************************************************
// X can be a terminal or non terminal
func (p *Parser) genericParseFunctionForRuleOfTypeXopXRepeat(parseX func() (Node, error), op []lexer.Token) (Node, error) {
	// F represents the function meant to parse "X"
	left, err := parseX()
	if err != nil {
		return nil, err
	}
	// General parsing strategy is of creating a left associative tree.
	// Above, X was already parsed. Now the next token shall be checked if it is
	// of the appropriate type. If it doesn't match, our parsing is done and we return.
	// Otherwise we consume it and then we again call "f" because the next part is expected to be "X".
	// Once X is parsed, we have a binary expression of the form  {X op X}.
	// We push this to the left (lets call it "W") and proceed ahead, trying to parse the next "op X".
	// Once it is parsed, the new expression will be {W op X}. Then this expression will be pushed to the left
	// and the new "op X" will be parsed. Convince yourself that this leads to a left associative tree.
	for {
		pl, err := p.l.Peek()
		if err != nil {
			return left, err
		}

		// If token type doesn't match we are done and can return the parsed expression thusfar
		if !tokenTypeIsOneOfTheList(pl, op) {
			return left, nil
		}

		// Else we create a binary expression
		op, err := p.l.FetchNextToken()
		if err != nil {
			return left, err
		}
		r, err := parseX()
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

func tokenTypeIsOneOfTheList(pl lexer.Token, l []lexer.Token) bool {
	for _, t := range l {
		if pl.TType == t.TType {
			return true
		}
	}
	return false
}
