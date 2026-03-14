package lexer

import "fmt"

type CharSource interface {
	GetNextChar() (val byte, hasNext bool, err error)
}

type Lexer struct {
	source CharSource
}

func NewLexer(s CharSource) *Lexer {
	return &Lexer{source: s}
}

func (l *Lexer) PrintAllChars() error {
	for {
		ch, hasNext, err := l.source.GetNextChar()
		if err != nil {
			return err
		}

		fmt.Printf("%c", ch)

		if !hasNext {
			return nil
		}
	}
}