package jsonparser

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

//go:generate goyacc -l -o parser.go parser.y

// This is a channel based implementation of lex.

// Parse parses the input and returs the result.
func Parse(input []byte) (map[string]interface{}, error) {
	l := newLex(input)
	defer func() { close(l.done) }()

	_ = yyParse(l)
	return l.result, l.err
}

type lex struct {
	input  []byte
	pos    int
	next   chan token
	done   chan struct{}
	result map[string]interface{}
	err    error
}

type token struct {
	typ int
	val interface{}
}

func newLex(input []byte) *lex {
	l := &lex{
		input: input,
		next:  make(chan token),
		done:  make(chan struct{}),
	}
	go l.run()
	return l
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	tok := <-l.next
	lval.val = tok.val
	return tok.typ
}

func (l *lex) run() {
	defer close(l.next)

	for b := l.nextb(); b != 0; b = l.nextb() {
		switch {
		case unicode.IsSpace(rune(b)):
			continue
		case b == '"':
			l.scanString()
		case unicode.IsDigit(rune(b)) || b == '+' || b == '-':
			l.backup()
			l.scanNum()
		case unicode.IsLetter(rune(b)):
			l.backup()
			l.scanLiteral()
		default:
			l.send(int(b), nil)
		}
	}
}

var escape = map[byte]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func (l *lex) scanString() {
	buf := bytes.NewBuffer(nil)
	for b := l.nextb(); b != 0; b = l.nextb() {
		switch b {
		case '\\':
			// TODO(sougou): handle \uxxxx construct.
			b2 := escape[l.nextb()]
			if b2 == 0 {
				l.send(LexError, nil)
				return
			}
			buf.WriteByte(b2)
		case '"':
			l.send(String, buf.String())
			return
		default:
			buf.WriteByte(b)
		}
	}
	l.send(LexError, nil)
}

func (l *lex) scanNum() {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.nextb()
		switch {
		case unicode.IsDigit(rune(b)):
			buf.WriteByte(b)
		case strings.IndexByte(".+-eE", b) != -1:
			buf.WriteByte(b)
		default:
			l.backup()
			val, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				l.send(LexError, nil)
				return
			}
			l.send(Number, val)
			return
		}
	}
}

var literal = map[string]interface{}{
	"true":  true,
	"false": false,
	"null":  nil,
}

func (l *lex) scanLiteral() {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.nextb()
		switch {
		case unicode.IsLetter(rune(b)):
			buf.WriteByte(b)
		default:
			l.backup()
			val, ok := literal[buf.String()]
			if !ok {
				l.send(LexError, nil)
				return
			}
			l.send(Literal, val)
			return
		}
	}
}

func (l *lex) send(typ int, val interface{}) {
	select {
	case l.next <- token{typ: typ, val: val}:
	case <-l.done:
	}
}

func (l *lex) backup() {
	if l.pos == -1 {
		return
	}
	l.pos--
}

func (l *lex) nextb() byte {
	select {
	case <-l.done:
		return 0
	default:
	}

	if l.pos >= len(l.input) || l.pos == -1 {
		l.pos = -1
		return 0
	}
	l.pos++
	return l.input[l.pos-1]
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}
