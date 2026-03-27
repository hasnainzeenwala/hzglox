package parser

import (
	"testing"

	"github.com/hasnainzeenwala/hzglox/lexer"
	"github.com/hasnainzeenwala/hzglox/sources"
)

func TestParser(t *testing.T) {
	type testCase struct {
		name        string
		program     string
		expectedAst string	
	}

	for _, r := range []testCase {
		{
			"sanity",
			`1 + 1 == 2`,
			`( == ( + 1 1 ) 2 )`,
		},
		{
			"more complex",
			` 3 < 2 + -4 == 7 < (8 == 9)`,
			`( == ( < 3 ( + 2 ( -4 ) ) ) ( < 7 ( group ( == 8 9 ) ) ) )`,

		},
	} {
		t.Run(r.name, func(t *testing.T) {
			p := NewParser(lexer.NewLexer(sources.NewStringToChars(r.program)))
			tree, err := p.Parse()
			if err != nil {
				t.Fatalf("unexpected error = %v", err)
			}
			gotAst := tree.PrintAst()
			if gotAst != r.expectedAst {
				t.Fatalf("Expected ast = %s . Got = %s", r.expectedAst, gotAst)
			}
		})
	}
}