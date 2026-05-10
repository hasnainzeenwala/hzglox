package ast

import (
	"testing"

	"github.com/hasnainzeenwala/hzglox/lexer"
)

func TestAstPrint(t *testing.T) {
	type testcase struct {
		name string
		e    Node
		s    string
	}
	for _, tt := range []testcase{
		{
			name: "literal",
			e:    &LiteralNode{lexer.Token{lexer.Number, "1.12", 0, 1.12}},
			s:    "1.12",
		},
		{
			name: "grouping",
			e: &GroupingNode{
				E: &LiteralNode{lexer.Token{lexer.True, "true", 0, nil}},
			},
			s: "( group true )",
		},
		{
			name: "unary",
			e: &UnaryNode{
				T: lexer.Token{lexer.Bang, "!", 0, nil},
				E: &LiteralNode{lexer.Token{lexer.False, "false", 0, nil}},
			},
			s: "( !false )",
		},
		{
			name: "binary",
			e: &Binary{
				Le: &LiteralNode{lexer.Token{lexer.Number, "1", 0, 1.0}},
				Op: lexer.Token{lexer.Plus, "+", 0, nil},
				Re: &LiteralNode{lexer.Token{lexer.Number, "2", 0, 2.0}},
			},
			s: "( + 1 2 )",
		},
		{
			name: "nested",
			e: &Binary{
				Le: &UnaryNode{
					T: lexer.Token{lexer.Minus, "-", 0, nil},
					E: &GroupingNode{
						E: &LiteralNode{lexer.Token{lexer.Number, "123", 0, 123.0}},
					},
				},
				Op: lexer.Token{lexer.Star, "*", 0, nil},
				Re: &LiteralNode{lexer.Token{lexer.Number, "45.67", 0, 45.67}},
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
