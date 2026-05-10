package parser

import (
	"fmt"

	"github.com/hasnainzeenwala/hzglox/lexer"
	"github.com/hasnainzeenwala/hzglox/ast"
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
// Parse() function will return the whole parse tree of the source code.
// Every non-terminal/rule shall have its own parsing function. It will be named 'parse<rulename>rule()'
// The recipe to create the parsing function is the following.
//     - Read the rule from left to right.
//     - If you encounter a non-terminal, call the corresponding function for that non-terminal.
//     - if you encounter a terminal, peek and see if the next lexeme matches it.
// 	         - If it does, consume the lexeme take the appropriate action.
//           - Otherwise, return because the rule no longer applies.
//     - Keep proceeding by matching the terminals and calling the functions for non-terminals
//       in a similar fashion till the rule is done.
// For some of the rules, a generic function has been created since all of them had a very similar structure.
// But the generic function follows the same idea described above.




type Parser struct {
	l *lexer.Lexer
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{l}
}

// args: empty
// result: Parse tree of the code
func (p *Parser) Parse() (ast.Node, error) {
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

// Rule: -> expression   ->  equal
// -----------------------------------------------------------------------------
// args: Empty
// Result: Result(equal rule)
func (p *Parser) parseExpressionRule() (ast.Node, error) {
	return p.parseEqualRule()
}

// Rule:  equal        ->  comparison (("==" | "!=") comparison)*
// -----------------------------------------------------------------------------
// args: Empty
// Result: Binary Node | Result(comparison rule)
func (p *Parser) parseEqualRule() (ast.Node, error) {
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

// Rule:   comparison   ->  term ( ("<" | "<=" | ">" | ">=") term )*
// -----------------------------------------------------------------------------
// args: Empty
// Result: Binary Node | Result(term rule)
func (p *Parser) parseComparisonRule() (ast.Node, error) {
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

// Rule:   term         ->  factor (( "+" | "-" ) factor)*
// -----------------------------------------------------------------------------
// args: Empty
// Result: Binary Node | Result(factor rule)
func (p *Parser) parseTermRule() (ast.Node, error) {
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

// Rule:   factor       ->  unary (("*" | "/") unary)*
// -----------------------------------------------------------------------------
// args: Empty
// Result: Binary Node | Result (unary rule)
func (p *Parser) parseFactorRule() (ast.Node, error) {
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

// Rule:   unary        ->  ("-" | "!") unary | primary
// -----------------------------------------------------------------------------
// args: Empty
// Result: Unary Node | Result(primary rule)
func (p *Parser) parseUnaryRule() (ast.Node, error) {
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
		return &ast.UnaryNode{
			T: t,
			E: r,
		}, nil
	default:
		return p.parsePrimaryRule()
	}
}

// Rule:   primary      ->  NUMBER | STRING | true | false | nil | "(" expression ")" ;
// -----------------------------------------------------------------------------
//
// args: empty
// Result: Literal Node | Grouping Node
func (p *Parser) parsePrimaryRule() (ast.Node, error) {
	t, err := p.l.FetchNextToken()
	if err != nil {
		return nil, err
	}
	switch t.TType {
	case lexer.Number, lexer.String, lexer.Nil, lexer.True, lexer.False:
		return &ast.LiteralNode{
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
		return &ast.GroupingNode{exp}, nil
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
// X is a non terminal
// args:
//    parseX: parsing function to parse the non-terminal X
//    op: list of tokens [op1, op2, .....]
// Result:
//    Binary tree | Result(parseX)
func (p *Parser) genericParseFunctionForRuleOfTypeXopXRepeat(parseX func() (ast.Node, error), op []lexer.Token) (ast.Node, error) {
	// First parse X. This is the first "Left" expression of the binary tree
	left, err := parseX()
	if err != nil {
		return nil, err
	}
	// Above, X was already parsed. The following steps will be taken now:
	//     - The next token shall be peeked to check if it lies in the list which was passed.
	//     - If it doesn't, our parsing is done for this rule and we return.
	//     - Otherwise we consume it and then call "parseX" because the next part is expected to be "X".
	//     - Once X is parsed, create a binary expression of the form  {X op X}.
	//     - Push this to the "left" (lets call it "W").
	//     - Again at the top of the loop now
	//     - Parse the next "op X".
	//     - Create new "left" expression {W op X}
	//     - Parse the next Op X
	//     - ....
	//     - Continue till you encounter a token which isn't in the list and return the the tree present in "left"
	//     - The tree is now your parsed binary expression
	// Convince yourself that this leads to a left associative tree.
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

		exp := &ast.Binary{
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
