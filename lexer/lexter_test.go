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
		{
			name: "keywordsOperatorsLoopsBraces",
			inputString: `for {
				if (x > 0) and (y <= 23.24) or (blah < 44) {
					while {
						print "Abcd xyz"
					}
				} else {
					return y
				}
			}`,
			expectedTokens: []Token{
				{TType: For, Lexeme: "for", LineNo: 1},
				{TType: LeftBrace, Lexeme: "{", LineNo: 1},
				{TType: If, Lexeme: "if", LineNo: 2},
				{TType: LeftParen, Lexeme: "(", LineNo: 2},
				{TType: Identifier, Lexeme: "x", LineNo: 2},
				{TType: Greater, Lexeme: ">", LineNo: 2},
				{TType: Number, Lexeme: "0", LineNo: 2, Literal: 0.0},
				{TType: RightParen, Lexeme: ")", LineNo: 2},
				{TType: And, Lexeme: "and", LineNo: 2},
				{TType: LeftParen, Lexeme: "(", LineNo: 2},
				{TType: Identifier, Lexeme: "y", LineNo: 2},
				{TType: LessEqual, Lexeme: "<=", LineNo: 2},
				{TType: Number, Lexeme: "23.24", LineNo: 2, Literal: 23.24},
				{TType: RightParen, Lexeme: ")", LineNo: 2},
				{TType: Or, Lexeme: "or", LineNo: 2},
				{TType: LeftParen, Lexeme: "(", LineNo: 2},
				{TType: Identifier, Lexeme: "blah", LineNo: 2},
				{TType: Less, Lexeme: "<", LineNo: 2},
				{TType: Number, Lexeme: "44", LineNo: 2, Literal: 44.0},
				{TType: RightParen, Lexeme: ")", LineNo: 2},
				{TType: LeftBrace, Lexeme: "{", LineNo: 2},
				{TType: While, Lexeme: "while", LineNo: 3},
				{TType: LeftBrace, Lexeme: "{", LineNo: 3},
				{TType: Print, Lexeme: "print", LineNo: 4},
				{TType: String, Lexeme: "\"Abcd xyz\"", LineNo: 4, Literal: "Abcd xyz"},
				{TType: RightBrace, Lexeme: "}", LineNo: 5},
				{TType: RightBrace, Lexeme: "}", LineNo: 6},
				{TType: Else, Lexeme: "else", LineNo: 6},
				{TType: LeftBrace, Lexeme: "{", LineNo: 6},
				{TType: Return, Lexeme: "return", LineNo: 7},
				{TType: Identifier, Lexeme: "y", LineNo: 7},
				{TType: RightBrace, Lexeme: "}", LineNo: 8},
				{TType: RightBrace, Lexeme: "}", LineNo: 9},
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
