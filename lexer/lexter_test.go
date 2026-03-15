package lexer

import (
	"reflect"
	"testing"
	"github.com/hasnainzeenwala/hzglox/sources"
)

func TestLexer(t *testing.T) {
	type testCase struct {
		name string
		inputString string
		expectedTokens []Token
	}
	for _, tt := range []testCase {
		{
			name: "sanity",
			inputString: `print "Hello, World!"`,
			expectedTokens: []Token{
				{
					TType: Print,
					Lexeme: "print",
					LineNo: 1,
				},
				{
					TType: String,
					Lexeme: "\"Hello, World!\"",
					LineNo: 1,
					Literal: "Hello, World!",
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			chars := sources.NewStringToChars(tt.inputString)
			lexer := NewLexer(chars)
			tokens := make([]Token, 0)

			for {
				tok, done, err := lexer.EmitToken()
				if err != nil {
					t.Fatalf("EmitToken() error = %v", err)
				}
				if done {
					break
				}
				tokens = append(tokens, tok)
			}

			if !reflect.DeepEqual(tokens, tt.expectedTokens) {
				t.Fatalf("tokens mismatch\n got: %#v\nwant: %#v", tokens, tt.expectedTokens)
			}
		})
	}
}
