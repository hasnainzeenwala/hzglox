package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

type CharSource interface {
	GetNextChar() (val byte, hasNext bool, err error)
	PeekNextChar() (val byte, err error)
}

type Lexer struct {
	source CharSource
	lineNo int
}

func NewLexer(s CharSource) *Lexer {
	return &Lexer{source: s, lineNo: 1}
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

func (l *Lexer) EmitToken() (t Token, done bool, e error) {
	var lexeme strings.Builder
	var ch byte
	var err error


	// skip whitespace before starting, and load the first character
	SkipWhiteSpace:
	for {
		ch, _, err = l.source.GetNextChar()
		if err != nil {
			t = Token{}
			done = true
			e = fmt.Errorf("ran into lexing error at line %d with ch=%q: %w", l.lineNo, ch, err)
			return
		}
		switch ch {
		case '\n':
			l.lineNo++
		case '\t', ' ', '\r':
		default:
			break SkipWhiteSpace
		}
	}

	lexeme.WriteByte(ch)
	switch {
	// ****************************
	// EOF
	// ****************************
	case ch == 0:
		t = Token{
			TType: Eof,
		}
		done = true
	
	// ******************************
	// ALl single char tokens
	// ******************************
	case ch == '(':
		t = Token{
			TType: LeftParen,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	case ch == ')':
		t = Token{
			TType:  RightParen,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	case ch == '{':
		t = Token{
			TType:  LeftBrace,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	case ch == '}':
		t = Token{
			TType:  RightBrace,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	case ch == ',':
		t = Token{
			TType:  Comma,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	case ch == '.':
		t = Token{
			TType:  Dot,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	case ch == '-':
		t = Token{
			TType:  Minus,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	case ch == '+':
		t = Token{
			TType:  Plus,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	case ch == ';':
		t = Token{
			TType:  Semicolon,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	case ch == '/':
		// need to handle the comments case
		nextCh, err := l.source.PeekNextChar()
		if err != nil {
			e = fmt.Errorf("error while peeking character at line %d with ch=%q: %w", l.lineNo, ch, err)
			done = true
			return
		}
		if nextCh == '/' {
			// keep consuming till the end of line because it's a comment
			for nextCh, err = l.source.PeekNextChar();
				nextCh != '\n' && err == nil;
				nextCh, err = l.source.PeekNextChar() {
				
				_, _, err = l.source.GetNextChar()
			}
			if err != nil {
				e = fmt.Errorf("Error while consuming a comment line no: %d, error: %w", l.lineNo, err)
				done = true
				return
			}
			e = nil
			done = false
			return
		} else {
			t = Token{
				TType:  Slash,
				Lexeme: lexeme.String(),
				LineNo: l.lineNo,
			}
		}
	case ch == '*':
		t = Token{
			TType:  Star,
			Lexeme: lexeme.String(),
			LineNo: l.lineNo,
		}
	
	// **********************************************
	// Double or single char special char lexemes
	// **********************************************
	case ch == '!':
		nextCh, err := l.source.PeekNextChar()
		if err != nil{
			e = fmt.Errorf("error while peeking character at line %d with ch=%q: %w", l.lineNo, ch, err)
			done = true
			return
		}
		if nextCh == '=' {
			ch, _, err := l.source.GetNextChar()
			if err != nil {
				e = fmt.Errorf("error while creating lexeme at line %d with ch=%q: %w", l.lineNo, ch, err)
				done = true
				return

			}
			lexeme.WriteByte(ch)
			t = Token{
				TType: BangEqual,
				Lexeme: lexeme.String(),
				LineNo: l.lineNo,
			}
		} else {
			t = Token{
				TType: Bang,
				Lexeme: lexeme.String(),
				LineNo: l.lineNo,
			}
		}
	case ch == '=':
		nextCh, err := l.source.PeekNextChar()
		if err != nil{
			e = fmt.Errorf("error while peeking character at line %d with ch=%q: %w", l.lineNo, ch, err)
			done = true
			return
		}
		if nextCh == '=' {
			ch, _, err := l.source.GetNextChar()
			if err != nil {
				e = fmt.Errorf("error while creating lexeme at line %d with ch=%q: %w", l.lineNo, ch, err)
				done = true
				return

			}
			lexeme.WriteByte(ch)
			t = Token{
				TType: EqualEqual,
				Lexeme: lexeme.String(),
				LineNo: l.lineNo,
			}
		} else {
			t = Token{
				TType: Equal,
				Lexeme: lexeme.String(),
				LineNo: l.lineNo,
			}
		}
	case ch == '<':
		nextCh, err := l.source.PeekNextChar()
		if err != nil{
			e = fmt.Errorf("error while peeking character at line %d with ch=%q: %w", l.lineNo, ch, err)
			done = true
			return
		}
		if nextCh == '=' {
			ch, _, err := l.source.GetNextChar()
			if err != nil {
				e = fmt.Errorf("error while creating lexeme at line %d with ch=%q: %w", l.lineNo, ch, err)
				done = true
				return

			}
			lexeme.WriteByte(ch)
			t = Token{
				TType: LessEqual,
				Lexeme: lexeme.String(),
				LineNo: l.lineNo,
			}
		} else {
			t = Token{
				TType: Less,
				Lexeme: lexeme.String(),
				LineNo: l.lineNo,
			}
		}
	case ch == '>':
		nextCh, err := l.source.PeekNextChar()
		if err != nil{
			e = fmt.Errorf("error while peeking character at line %d with ch=%q: %w", l.lineNo, ch, err)
			done = true
			return
		}
		if nextCh == '=' {
			ch, _, err := l.source.GetNextChar()
			if err != nil {
				e = fmt.Errorf("error while creating lexeme at line %d with ch=%q: %w", l.lineNo, ch, err)
				done = true
				return

			}
			lexeme.WriteByte(ch)
			t = Token{
				TType: GreaterEqual,
				Lexeme: lexeme.String(),
				LineNo: l.lineNo,
			}
		} else {
			t = Token{
				TType: Greater,
				Lexeme: lexeme.String(),
				LineNo: l.lineNo,
			}
		}

	// ****************************
	// Literals
	// ****************************

	// If they start with a number then keep going till you are finding numbers
	// 1 period in the middle is allowed for "double"
	case ch >= '0' && ch <= '9':
		sawAPeriod := false
		var nextCh byte
		var err error

		// Keep adding characters to the lexeme till you are finding numbers
		// or 1 period
		for nextCh, err = l.source.PeekNextChar();
			((nextCh >= '0' && nextCh <= '9') ||
				(nextCh == '.' && sawAPeriod == false)) &&
			(err == nil); 
			nextCh, err = l.source.PeekNextChar() {
			ch, _, err = l.source.GetNextChar()
			if ch == '.' {
				sawAPeriod = true
			}
			lexeme.WriteByte(ch)
		}
		if err != nil {
			e = fmt.Errorf("Error while creating a number literal lexeme at line no: %d, lexeme: %s, error: %w", l.lineNo, lexeme.String(), err)
			done = true
			return
		}

		// for the lexeme to be a valid literal the last char has to be a number
		lexemeVal := lexeme.String()
		lastCh := lexemeVal[len(lexemeVal) - 1]
		if !(lastCh >= '0' && lastCh <= '9') {
			e = fmt.Errorf("Error while creating a lexeme at line no: %d, lexeme: %s, error: a valid number must end with a digit", l.lineNo, lexemeVal)
			done = true
			return
		}
		f, err := strconv.ParseFloat(lexemeVal, 64)
		if err != nil {
			e = fmt.Errorf("Error while creating a lexeme at line no: %d, lexeme: %s, error: a valid number must end with a digit", l.lineNo, lexemeVal)
			done = true
			return
		}

		t = Token{
			TType: Number,
			LineNo: l.lineNo,
			Lexeme: lexemeVal,
			Literal: f,
		}

	// Parse string literals multi-line string are allowed
	// we are not supporting escape sequences, i.e., as soon as we encounter a '"'
	// the string will consider ended. All the \t \n will be taken literally.
	case ch == '"':
		beginningLineNo := l.lineNo
		for {
			ch, _, err = l.source.GetNextChar()
			if err != nil {
				e = fmt.Errorf("Error while reading string on line no: %d, lexeme: %s, error: %w", l.lineNo, lexeme.String(), err)
				done = true
				return
			}
			lexeme.WriteByte(ch)
			if ch == 0 {
				e = fmt.Errorf("Error while reading string on line no: %d, lexeme: %s, error: didn't find closing \"", l.lineNo, lexeme.String())
				done = true
				return
			}
			if ch == '"' {
				break
			}
			if ch == '\n' {
				l.lineNo++
			}
		}
		literal := ""
		lexemeStr := lexeme.String()
		if len(lexemeStr) > 2 {
			literal = lexemeStr[1:(len(lexemeStr)-1)]	
		}

		t = Token{
			TType: String,
			LineNo: beginningLineNo,
			Literal: literal,
			Lexeme: lexeme.String(),
		}
	
	// If it starts with a letter keep going till you have alph-numerics
	// Then check if it is a reserved keyword
	case (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z'):
		for {
			nextCh, err := l.source.PeekNextChar()
			if err != nil {
				done = true
				e = fmt.Errorf("Error file parsing keyword/identifier on line: %d lexeme: %s error: %w", l.lineNo, lexeme.String(), err)
				return
			}

			if (nextCh >= 'a' && nextCh <= 'z') || (nextCh >= 'A' && nextCh <= 'Z') ||
				(nextCh == '_') || (nextCh >= '0' && nextCh <= '9') {
			
				ch, _, err = l.source.GetNextChar()
				if err != nil {
					done = true
					e = fmt.Errorf("Error file parsing keyword/identifier on line: %d lexeme: %s error: %w", l.lineNo, lexeme.String(), err)
					return
				}
				lexeme.WriteByte(ch)
			} else {
				break
			}
		}
		lexemeStr := lexeme.String()
		if tp, ok := Keywords[lexemeStr]; ok {
			t = Token{
				TType: tp,
				Lexeme: lexemeStr,
				LineNo: l.lineNo,
			}
		} else {
			t = Token{
				TType: Identifier,
				Lexeme: lexemeStr,
				LineNo: l.lineNo,
			}
		}

	
	// *******************************
	// Unidentified
	// *******************************
	default:
		done = true
		e = fmt.Errorf("Unidentified character on line no: %d lexeme: %s", l.lineNo, lexeme.String())
	}

	return
}
