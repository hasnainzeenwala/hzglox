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
	buffer Token
}

func NewLexer(s CharSource) *Lexer {
	return &Lexer{source: s, lineNo: 1}
}

func (l *Lexer) newError(action string, ch byte, lexeme string, reason string, cause error) error {
	if cause != nil {
		return fmt.Errorf(
			"lexer error while %s: line=%d ch=%q byte=%d lexeme=%q reason=%s: %w",
			action,
			l.lineNo,
			ch,
			ch,
			lexeme,
			reason,
			cause,
		)
	}
	return fmt.Errorf(
		"lexer error while %s: line=%d ch=%q byte=%d lexeme=%q reason=%s",
		action,
		l.lineNo,
		ch,
		ch,
		lexeme,
		reason,
	)
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

func (l *Lexer) EmitToken() (t Token, e error) {
	if l.buffer.TType != Undefined {
		t = l.buffer
		l.buffer = Token{}
		return
	}
	var lexeme strings.Builder
	var ch byte
	var err error


	// skip whitespace before starting, and load the first character
	SkipWhiteSpace:
	for {
		ch, _, err = l.source.GetNextChar()
		if err != nil {
			t = Token{}
			e = l.newError("reading next character", ch, lexeme.String(), "failed to get next character", err)
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
			e = l.newError("peeking after slash", ch, lexeme.String(), "failed to peek next character", err)
			return
		}
		if nextCh == '/' {
			// keep consuming till the end of line because it's a comment
			for nextCh, err = l.source.PeekNextChar();
				nextCh != '\n' && err == nil && nextCh != 0;
				nextCh, err = l.source.PeekNextChar() {
				
				_, _, err = l.source.GetNextChar()
			}
			if err != nil {
				e = l.newError("consuming comment", ch, lexeme.String(), "failed while skipping comment contents", err)
				return
			}
			// Code comment is not a token so we call EmitToken() again to get the next valid token
			return l.EmitToken()
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
			e = l.newError("peeking after bang", ch, lexeme.String(), "failed to peek next character", err)
			return
		}
		if nextCh == '=' {
			ch, _, err := l.source.GetNextChar()
			if err != nil {
				e = l.newError("building bang-equal lexeme", ch, lexeme.String(), "failed to consume '='", err)
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
			e = l.newError("peeking after equal", ch, lexeme.String(), "failed to peek next character", err)
			return
		}
		if nextCh == '=' {
			ch, _, err := l.source.GetNextChar()
			if err != nil {
				e = l.newError("building equal-equal lexeme", ch, lexeme.String(), "failed to consume '='", err)
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
			e = l.newError("peeking after less-than", ch, lexeme.String(), "failed to peek next character", err)
			return
		}
		if nextCh == '=' {
			ch, _, err := l.source.GetNextChar()
			if err != nil {
				e = l.newError("building less-equal lexeme", ch, lexeme.String(), "failed to consume '='", err)
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
			e = l.newError("peeking after greater-than", ch, lexeme.String(), "failed to peek next character", err)
			return
		}
		if nextCh == '=' {
			ch, _, err := l.source.GetNextChar()
			if err != nil {
				e = l.newError("building greater-equal lexeme", ch, lexeme.String(), "failed to consume '='", err)
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
			e = l.newError("building number literal", ch, lexeme.String(), "failed while extending numeric lexeme", err)
			return
		}

		// for the lexeme to be a valid literal the last char has to be a number
		lexemeVal := lexeme.String()
		lastCh := lexemeVal[len(lexemeVal) - 1]
		if !(lastCh >= '0' && lastCh <= '9') {
			e = l.newError("validating number literal", lastCh, lexemeVal, "a valid number must end with a digit", nil)
			return
		}
		f, err := strconv.ParseFloat(lexemeVal, 64)
		if err != nil {
			e = l.newError("parsing number literal", ch, lexemeVal, "failed to parse number", err)
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
				e = l.newError("reading string literal", ch, lexeme.String(), "failed to get next character", err)
				return
			}
			lexeme.WriteByte(ch)
			if ch == 0 {
				e = l.newError("reading string literal", ch, lexeme.String(), "did not find closing quote", nil)
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
				e = l.newError("peeking identifier continuation", ch, lexeme.String(), "failed to peek next character", err)
				return
			}

			if (nextCh >= 'a' && nextCh <= 'z') || (nextCh >= 'A' && nextCh <= 'Z') ||
				(nextCh == '_') || (nextCh >= '0' && nextCh <= '9') {
			
				ch, _, err = l.source.GetNextChar()
				if err != nil {
					e = l.newError("reading identifier continuation", ch, lexeme.String(), "failed to consume identifier character", err)
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
		e = l.newError("classifying token", ch, lexeme.String(), "unidentified character", nil)
	}

	return
}

func (l *Lexer) Peek() (t Token, e error) {
	t, e = l.EmitToken()
	if e != nil {
		l.buffer = t
	}
	return
}