package parser

import (
	"fmt"
	"testing"

	"github.com/hasnainzeenwala/hzglox/lexer"
	"github.com/hasnainzeenwala/hzglox/sources"
)

func TestParser(t *testing.T) {
	type testCase struct {
		name        string
		program     string
		expectedAst string	
		e           error
	}

	for _, r := range []testCase {
		{
			"sanity",
			`1 + 1 == 2`,
			`( == ( + 1 1 ) 2 )`,
			nil,
		},
		{
			"more complex",
			` 3 < 2 + -4 == 7 < (8 == 9)`,
			`( == ( < 3 ( + 2 ( -4 ) ) ) ( < 7 ( group ( == 8 9 ) ) ) )`,
			nil,
		},
		{
			"Whole expression can't be parsed",
			`3 < 2 {`,
			`( < 3 2 )`,
			fmt.Errorf("Encountered an unparseable token (Type: LeftBrace Lexeme: { LineNo: 1)"),
		},
	} {
		t.Run(r.name, func(t *testing.T) {
			p := NewParser(lexer.NewLexer(sources.NewStringToChars(r.program)))
			tree, err := p.Parse()
			if r.e == nil && err != nil {
				t.Fatalf("unexpected error = %v", err)
			}
			gotAst := tree.PrintAst()
			if gotAst != r.expectedAst {
				t.Fatalf("Expected ast = %s . Got = %s", r.expectedAst, gotAst)
			}
			if r.e != nil && r.e.Error() != err.Error() {
				t.Fatalf("Expected error = %q got = %q", r.e, err)
			}
		})
	}
}