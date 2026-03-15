package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"github.com/hasnainzeenwala/hzglox/lexer"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: main <optional_filepath>")
		os.Exit(65)
	}

	if len(args) == 1 {
		// **************************************
		// Run file through the interpreter
		// **************************************
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %q: %v\n", args[0], err)
			os.Exit(1)
		}
		defer file.Close()

		chars := NewFileCharSource(file)

		// Create a new instance of the interpreter (JUST LEXER FOR NOW)
		lex := lexer.NewLexer(chars)

		if err := lex.PrintAllChars(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to read %q: %v\n", args[0], err)
			os.Exit(1)
		}
		return
	} else {
		// **************************************
		// Interactive mode
		// **************************************
		lineReader := NewLineReader()
		for {
			val, isEOF, err := lineReader.ReadLine()
		
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
				os.Exit(1)
			}

			if isEOF {
				os.Exit(0)
			}

			chars := NewStringToChars(val)

			// Create a new instance of the interpreter (JUST LEXER FOR NOW)
			lex := lexer.NewLexer(chars)
			
			// send chars to lexer
			if err := lex.PrintAllChars(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to process input: %v\n", err)
				os.Exit(1)
			}
		}
	}

}

type StringToChars struct {
	line string
	p int
}

func NewStringToChars(line string) *StringToChars {
	return &StringToChars{ line: line, p: 0 }
}

func (s *StringToChars) GetNextChar() (val byte, hasNext bool, err error) {
	if s.p >= len(s.line) {
		return byte(0), false, nil
	}
	ch := s.line[s.p]
	s.p++
	return ch, s.p < len(s.line), nil
}

func (s *StringToChars) PeekNextChar() (ch byte, err error) {
	if s.p >= len(s.line) {
		return byte(0), nil
	}
	return s.line[s.p], nil
}

type FileCharSource struct {
	reader *bufio.Reader
}

func NewFileCharSource(file *os.File) *FileCharSource {
	return &FileCharSource{
		reader: bufio.NewReader(file),
	}
}

func (f *FileCharSource) GetNextChar() (val byte, hasNext bool, err error) {
	ch, err := f.reader.ReadByte()
	if err == io.EOF {
		return byte(0), false, nil
	}
	if err != nil {
		return byte(0), false, err
	}
	return ch, true, nil
}

func (f *FileCharSource) PeekNextChar() (ch byte, err error) {
	bytes, err := f.reader.Peek(1)
	if err == io.EOF {
		return byte(0), nil
	}
	if err != nil {
		return byte(0), err
	}
	return bytes[0], nil
}

type LineReader struct {
	scanner *bufio.Scanner
}

func NewLineReader() *LineReader {
	return &LineReader{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (l *LineReader) ReadLine() (string, bool, error) {
	fmt.Print("> ")
	notEOF := l.scanner.Scan()
	val := l.scanner.Text()
	return val, !notEOF, l.scanner.Err()
}
