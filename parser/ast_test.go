package parser

import (
	"testing"

	"github.com/hasnainzeenwala/hzglox/lexer"
)

func TestNewLiteral(t *testing.T) {
	valid := lexer.Token{TType: lexer.Number, Lexeme: "123", LineNo: 1, Literal: 123.0}
	invalid := lexer.Token{TType: lexer.Identifier, Lexeme: "abc", LineNo: 1}

	if _, match := NewLiteral(valid); !match {
		t.Fatalf("NewLiteral(valid) match = false, want true")
	}

	if _, match := NewLiteral(invalid); match {
		t.Fatalf("NewLiteral(invalid) match = true, want false")
	}
}

func TestNewGrouping(t *testing.T) {
	inner, _ := NewLiteral(lexer.Token{TType: lexer.True, Lexeme: "true", LineNo: 1})

	if _, match := NewGrouping(
		lexer.Token{TType: lexer.LeftParen, Lexeme: "(", LineNo: 1},
		inner,
		lexer.Token{TType: lexer.RightParen, Lexeme: ")", LineNo: 1},
	); !match {
		t.Fatalf("NewGrouping(valid) match = false, want true")
	}

	if _, match := NewGrouping(
		lexer.Token{TType: lexer.LeftBrace, Lexeme: "{", LineNo: 1},
		inner,
		lexer.Token{TType: lexer.RightParen, Lexeme: ")", LineNo: 1},
	); match {
		t.Fatalf("NewGrouping(invalid) match = true, want false")
	}

	if _, match := NewGrouping(
		lexer.Token{TType: lexer.LeftParen, Lexeme: "(", LineNo: 1},
		nil,
		lexer.Token{TType: lexer.RightParen, Lexeme: ")", LineNo: 1},
	); match {
		t.Fatalf("NewGrouping(nil expr) match = true, want false")
	}
}

func TestNewUnary(t *testing.T) {
	expr, _ := NewLiteral(lexer.Token{TType: lexer.False, Lexeme: "false", LineNo: 1})

	if _, match := NewUnary(
		lexer.Token{TType: lexer.Bang, Lexeme: "!", LineNo: 1},
		expr,
	); !match {
		t.Fatalf("NewUnary(valid) match = false, want true")
	}

	if _, match := NewUnary(
		lexer.Token{TType: lexer.Plus, Lexeme: "+", LineNo: 1},
		expr,
	); match {
		t.Fatalf("NewUnary(invalid token) match = true, want false")
	}

	if _, match := NewUnary(
		lexer.Token{TType: lexer.Minus, Lexeme: "-", LineNo: 1},
		nil,
	); match {
		t.Fatalf("NewUnary(nil expr) match = true, want false")
	}
}

func TestNewBinary(t *testing.T) {
	left, _ := NewLiteral(lexer.Token{TType: lexer.Number, Lexeme: "1", LineNo: 1, Literal: 1.0})
	right, _ := NewLiteral(lexer.Token{TType: lexer.Number, Lexeme: "2", LineNo: 1, Literal: 2.0})

	if _, match := NewBinary(
		left,
		lexer.Token{TType: lexer.Plus, Lexeme: "+", LineNo: 1},
		right,
	); !match {
		t.Fatalf("NewBinary(valid) match = false, want true")
	}

	if _, match := NewBinary(
		left,
		lexer.Token{TType: lexer.Equal, Lexeme: "=", LineNo: 1},
		right,
	); match {
		t.Fatalf("NewBinary(invalid operator) match = true, want false")
	}

	if _, match := NewBinary(
		nil,
		lexer.Token{TType: lexer.Star, Lexeme: "*", LineNo: 1},
		right,
	); match {
		t.Fatalf("NewBinary(nil left expr) match = true, want false")
	}

	if _, match := NewBinary(
		left,
		lexer.Token{TType: lexer.Star, Lexeme: "*", LineNo: 1},
		nil,
	); match {
		t.Fatalf("NewBinary(nil right expr) match = true, want false")
	}
}

func TestAstPrint(t *testing.T) {
	type testcase struct {
		name string
		e Expr
		s string
	}
	for _, tt := range []testcase {
		{
			name: "literal",
			e: Literal{lexer.Token{lexer.Number, "1.12", 0, 1.12}},
			s: "1.12",
		},
		{
			name: "grouping",
			e: Grouping{
				E: Literal{lexer.Token{lexer.True, "true", 0, nil}},
			},
			s: "( group true )",
		},
		{
			name: "unary",
			e: Unary{
				T: lexer.Token{lexer.Bang, "!", 0, nil},
				E: Literal{lexer.Token{lexer.False, "false", 0, nil}},
			},
			s: "( !false )",
		},
		{
			name: "binary",
			e: Binary{
				Le: Literal{lexer.Token{lexer.Number, "1", 0, 1.0}},
				Op: lexer.Token{lexer.Plus, "+", 0, nil},
				Re: Literal{lexer.Token{lexer.Number, "2", 0, 2.0}},
			},
			s: "( + 1 2 )",
		},
		{
			name: "nested",
			e: Binary{
				Le: Unary{
					T: lexer.Token{lexer.Minus, "-", 0, nil},
					E: Grouping{
						E: Literal{lexer.Token{lexer.Number, "123", 0, 123.0}},
					},
				},
				Op: lexer.Token{lexer.Star, "*", 0, nil},
				Re: Literal{lexer.Token{lexer.Number, "45.67", 0, 45.67}},
			},
			s: "( * ( -( group 123 ) ) 45.67 )",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.PrintAst(); got != tt.s {
				t.Fatalf("PrintAst() = %q, want %q", got, tt.s)
			}
		})
	}
}
