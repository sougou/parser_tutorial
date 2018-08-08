package jsonparser

import (
	"errors"
)

//go:generate goyacc -l -o parser.go parser.y

// Result is the type of the parser result
type Result int

// Parse parses the string and returs the result.
func Parse(input []byte) (Result, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.result, l.err
}

type lex struct {
	next   chan int
	result Result
	err    error
}

func newLex(input []byte) *lex {
	l := &lex{
		next: make(chan int),
	}
	go l.scan(input)
	return l
}

func (l *lex) scan(input []byte) {
	defer close(l.next)

	for _, ch := range input {
		switch {
		case '0' <= ch && ch <= '9':
			l.next <- Digit
		default:
			l.next <- int(ch)
		}
	}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	return <-l.next
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}
