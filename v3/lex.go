package jsonparser

import (
	"errors"
)

//go:generate goyacc -l -o parser.go parser.y

// Result is the type of the parser result
type Result interface{}

// Parse parses the string and returs the result.
func Parse(input []byte) (Result, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.result, l.lastError
}

type lex struct {
	next      chan token
	result    Result
	lastError error
}

type token struct {
	typ int
	val byte
}

func newLex(input []byte) *lex {
	l := &lex{
		next: make(chan token),
	}
	go l.scan(input)
	return l
}

func (l *lex) scan(input []byte) {
	defer close(l.next)

	for _, ch := range input {
		switch {
		case '0' <= ch && ch <= '9':
			l.send(Digit, ch)
		default:
			l.send(int(ch), 0)
		}
	}
}

func (l *lex) send(typ int, val byte) {
	l.next <- token{typ: typ, val: val}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	tok := <-l.next
	lval.ch = tok.val
	return tok.typ
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.lastError = errors.New(s)
}
