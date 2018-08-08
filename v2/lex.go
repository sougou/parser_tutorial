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
	result Result
	err    error
}

func newLex(input []byte) *lex {
	return &lex{}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	return 0
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}
