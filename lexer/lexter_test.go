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
		expectedErr error
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
				{
					TType: Eof,
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
				{TType: Eof},
			},
		},
		{
			name: "CodeWithComments",
			inputString: `// This is some description
			for blah { s = "some string"}
			// last line is a comment and there is no newline`,
			expectedTokens: []Token{
				{TType: For, Lexeme: "for", LineNo: 2},
				{TType: Identifier, Lexeme: "blah", LineNo: 2},
				{TType: LeftBrace, Lexeme: "{", LineNo: 2},
				{TType: Identifier, Lexeme: "s", LineNo: 2},
				{TType: Equal, Lexeme: "=", LineNo: 2},
				{TType: String, Lexeme: "\"some string\"", LineNo: 2, Literal: "some string"},
				{TType: RightBrace, Lexeme: "}", LineNo: 2},
				{TType: Eof},
			},
		},
		{
			name: "StringNotClosed",
			inputString: `// Things are fine
			// we have a for loop
			for x > 10 { a = "forgot to close }`,
			expectedTokens: []Token{
				{TType: For, Lexeme: "for", LineNo: 3},
				{TType: Identifier, Lexeme: "x", LineNo: 3},
				{TType: Greater, Lexeme: ">", LineNo: 3},
				{TType: Number, Lexeme: "10", LineNo: 3, Literal: 10.0},
				{TType: LeftBrace, Lexeme: "{", LineNo: 3},
				{TType: Identifier, Lexeme: "a", LineNo: 3},
				{TType: Equal, Lexeme: "=", LineNo: 3},
			},
			expectedErr: (&Lexer{lineNo: 3}).newError(
				"reading string literal",
				0,
				"\"forgot to close }",
				"did not find closing quote",
				nil,
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			chars := sources.NewStringToChars(tt.inputString)
			lexer := NewLexer(chars)
			tokens := make([]Token, 0)
			var gotErr error

			for {
				tok, err := lexer.EmitToken()
				if err != nil {
					gotErr = err
					break
				}
				tokens = append(tokens, tok)
				if tok.TType == Eof {
					break
				}
			}

			if !reflect.DeepEqual(tokens, tt.expectedTokens) {
				t.Fatalf("tokens mismatch\n got: %#v\nwant: %#v", tokens, tt.expectedTokens)
			}

			if tt.expectedErr == nil {
				if gotErr != nil {
					t.Fatalf("unexpected error = %v", gotErr)
				}
			} else {
				if gotErr == nil {
					t.Fatalf("expected error %q but got nil", tt.expectedErr.Error())
				}
				if gotErr.Error() != tt.expectedErr.Error() {
					t.Fatalf("error mismatch\n got: %q\nwant: %q", gotErr.Error(), tt.expectedErr.Error())
				}
			}
		})
	}
}
