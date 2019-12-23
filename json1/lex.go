package jsonparser

import (
	"bytes"
	"errors"
)

//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser.go parser.y

// Result is the type of the parser result
type Result int

// Parse parses the input and returs the result.
func Parse(input []byte) (Result, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.result, l.err
}

type lex struct {
	input  *bytes.Buffer
	result Result
	err    error
}

func newLex(input []byte) *lex {
	return &lex{
		input: bytes.NewBuffer(input),
	}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	return int(l.nextb())
}

func (l *lex) nextb() byte {
	b, err := l.input.ReadByte()
	if err != nil {
		return 0
	}
	return b
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}
