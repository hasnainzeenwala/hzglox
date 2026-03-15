package lexer

import "fmt"

type TokenType int

const (
	// UNDEFIEND
	Undefined TokenType = iota

	// Single-character tokens.
	LeftParen
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// One or two character tokens.
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals.
	Identifier
	String
	Number

	// Keywords.
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	// End of file.
	Eof
)

var Keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func (t TokenType) String() string {
	switch t {
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case LeftBrace:
		return "LeftBrace"
	case RightBrace:
		return "RightBrace"
	case Comma:
		return "Comma"
	case Dot:
		return "Dot"
	case Minus:
		return "Minus"
	case Plus:
		return "Plus"
	case Semicolon:
		return "Semicolon"
	case Slash:
		return "Slash"
	case Star:
		return "Star"
	case Bang:
		return "Bang"
	case BangEqual:
		return "BangEqual"
	case Equal:
		return "Equal"
	case EqualEqual:
		return "EqualEqual"
	case Greater:
		return "Greater"
	case GreaterEqual:
		return "GreaterEqual"
	case Less:
		return "Less"
	case LessEqual:
		return "LessEqual"
	case Identifier:
		return "Identifier"
	case String:
		return "String"
	case Number:
		return "Number"
	case And:
		return "and"
	case Class:
		return "class"
	case Else:
		return "else"
	case False:
		return "false"
	case Fun:
		return "fun"
	case For:
		return "for"
	case If:
		return "if"
	case Nil:
		return "nil"
	case Or:
		return "or"
	case Print:
		return "print"
	case Return:
		return "return"
	case Super:
		return "super"
	case This:
		return "this"
	case True:
		return "true"
	case Var:
		return "var"
	case While:
		return "while"
	case Eof:
		return "Eof"
	default:
		return "UnknownToken"
	}
}

type Token struct {
	TType TokenType
	Lexeme string
	LineNo int
	Literal any
}

func NewToken(ttype TokenType, lexeme string, lineNo int, literal any) Token {
	return Token{TType: ttype, Lexeme: lexeme, LineNo: lineNo, Literal: literal}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", t.TType, t.Lexeme, t.Literal)
}
