package jsonparser

import (
	"bytes"
	"errors"
	"unicode"
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
	nextByte  chan byte
	next      chan token
	result    Result
	lastError error
}

type token struct {
	typ int
	val interface{}
}

func newLex(input []byte) *lex {
	l := &lex{
		nextByte: make(chan byte),
		next:     make(chan token),
	}
	// Pump the bytes.
	go func() {
		defer close(l.nextByte)
		for _, ch := range input {
			l.nextByte <- ch
		}
	}()
	// pump the tokens.
	go l.scan()
	return l
}

func (l *lex) scan() {
	defer close(l.next)
	l.readNormal()
}

func (l *lex) send(typ int, val interface{}) {
	l.next <- token{typ: typ, val: val}
}

func (l *lex) readNormal() {
	for ch := range l.nextByte {
		switch {
		case unicode.IsSpace(rune(ch)):
			continue
		case '0' <= ch && ch <= '9':
			l.send(Digit, ch)
			l.readNum()
			continue
		case ch == '+' || ch == '-' || ch == '.':
			l.send(int(ch), 0)
			l.readNum()
		case ch == '"':
			l.readString()
		default:
			l.send(int(ch), 0)
		}
	}
}

func (l *lex) readString() {
	buf := bytes.NewBuffer(nil)
	for ch := range l.nextByte {
		switch ch {
		case '\\':
			// TODO(sougou): handle \uxxx construct.
			escaped := escape[<-l.nextByte]
			if escaped == 0 {
				l.send(LexError, 0)
				return
			}
			buf.WriteByte(escaped)
		case '"':
			l.send(String, buf.String())
			return
		case 0:
			l.send(LexError, 0)
			return
		default:
			buf.WriteByte(ch)
		}
	}
}

func (l *lex) readNum() {
	for ch := range l.nextByte {
		if '0' <= ch && ch <= '9' {
			l.send(Digit, ch)
			continue
		}
		l.send(int(ch), 0)
		switch ch {
		case '+', '-', '.', 'e', 'E':
			// continue
		default:
			return
		}
	}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	tok := <-l.next
	switch v := tok.val.(type) {
	case byte:
		lval.ch = v
	case string:
		lval.str = v
	}
	return tok.typ
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.lastError = errors.New(s)
}

var escape = map[byte]byte{
	'b': '\b',
	'f': '\f',
	'n': '\n',
	'r': '\r',
	't': '\t',
}
