package ast

import (
	"reflect"
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
			e: &BinaryNode{
				Le: &LiteralNode{lexer.Token{lexer.Number, "1", 0, 1.0}},
				Op: lexer.Token{lexer.Plus, "+", 0, nil},
				Re: &LiteralNode{lexer.Token{lexer.Number, "2", 0, 2.0}},
			},
			s: "( + 1 2 )",
		},
		{
			name: "nested",
			e: &BinaryNode{
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

// Test if a literal node evaluates to itself or not
// Literal types:
//    *Number
//    *String
//    *true
//    *false
//    *nil
func TestAstLiteralInterpret(t *testing.T) {

	// interpreting the node should yield "expectedVal"
	type testCase struct {
		expectedVal  any
		node         *LiteralNode
		err          error
	}

	for _, tt := range []testCase{
		{
			expectedVal: 2.3,
			node       : &LiteralNode{
				T: lexer.NewToken(lexer.Number, "2.3", 1, 2.3),
			},
		},
		{
			expectedVal: "this is some string",
			node       : &LiteralNode{
				T: lexer.NewToken(lexer.String, "\"this is some string\"", 1, "this is some string"),
			},
		},
		{
			expectedVal: true,
			node       : &LiteralNode{
				T: lexer.NewToken(lexer.True, "true", 1, true),
			},
		},
		{
			expectedVal: false,
			node       : &LiteralNode{
				T: lexer.NewToken(lexer.False, "false", 1, false),
			},
		},
		{
			expectedVal: nil,
			node       : &LiteralNode{
				T: lexer.NewToken(lexer.Nil, "nil", 1, nil),
			},
		},
	} {
		if tt.err == nil {
			gotVal, err := tt.node.Interpret()
			if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}
			if !reflect.DeepEqual(gotVal, tt.expectedVal) {
				t.Fatalf("Expected: %v but got: %v", tt.expectedVal, gotVal)
			}
		}
	}
}

func TestInterpretTrees(t *testing.T) {
	type testCase struct {
		name        string
		expectedVal any
		tree        Node
		expectedErr error
	}
	for _, tt := range []testCase {
		{
			name: "Test Addition",
			expectedVal: float64(3),
			tree: &BinaryNode{
				Le: &LiteralNode{
					T: lexer.NewToken(lexer.Number, "1", 1, 1),
				},
				Op: lexer.NewToken(lexer.Plus, "+", 1, nil),
				Re: &LiteralNode{
					T: lexer.NewToken(lexer.Number, "2", 1, 2),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.tree.Interpret()
			if tt.expectedErr != nil{
				if err.Error() != tt.expectedErr.Error() {
					t.Fatalf("Expected error: %v Got error: %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			} else {
				if !reflect.DeepEqual(tt.expectedVal, gotVal) {
					t.Fatalf("Expected value: %v but got value: %v", tt.expectedVal, gotVal)
				}
			}
		})

	}
}