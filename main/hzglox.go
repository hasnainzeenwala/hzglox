package main

import (
	"fmt"
	"os"
	"github.com/hasnainzeenwala/hzglox/lexer"
	"github.com/hasnainzeenwala/hzglox/sources"
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

		chars := sources.NewFileCharSource(file)

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
		lineReader := sources.NewLineReader()
		for {
			val, isEOF, err := lineReader.ReadLine()
		
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
				os.Exit(1)
			}

			if isEOF {
				os.Exit(0)
			}

			chars := sources.NewStringToChars(val)

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
