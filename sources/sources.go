package sources

import (
	"bufio"
	"os"
	"io"
	"fmt"
)

type StringToChars struct {
	line string
	p    int
}

func NewStringToChars(line string) *StringToChars {
	return &StringToChars{line: line, p: 0}
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