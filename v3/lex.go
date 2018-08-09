package jsonparser

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

//go:generate goyacc -l -o parser.go parser.y

// Parse parses the string and returs the result.
func Parse(input []byte) (map[string]interface{}, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.result, l.err
}

type lex struct {
	input  *bytes.Buffer
	pushed byte
	result map[string]interface{}
	err    error
}

type token struct {
	typ int
	val interface{}
}

func newLex(input []byte) *lex {
	return &lex{
		input: bytes.NewBuffer(input),
	}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	return l.scanNormal(lval)
}

func (l *lex) scanNormal(lval *yySymType) int {
	for b := l.nextb(); b != 0; b = l.nextb() {
		switch {
		case unicode.IsSpace(rune(b)):
			continue
		case b == '"':
			return l.scanString(lval)
		case unicode.IsDigit(rune(b)) || b == '+' || b == '-':
			l.pushb(b)
			return l.scanNum(lval)
		case unicode.IsLetter(rune(b)):
			l.pushb(b)
			return l.scanLiteral(lval)
		default:
			return int(b)
		}
	}
	return 0
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

func (l *lex) scanString(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for b := l.nextb(); b != 0; b = l.nextb() {
		switch b {
		case '\\':
			// TODO(sougou): handle \uxxxx construct.
			b2 := escape[l.nextb()]
			if b2 == 0 {
				return LexError
			}
			buf.WriteByte(b2)
		case '"':
			lval.val = buf.String()
			return String
		default:
			buf.WriteByte(b)
		}
	}
	return 0
}

func (l *lex) scanNum(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for b := l.nextb(); b != 0; b = l.nextb() {
		switch {
		case unicode.IsDigit(rune(b)):
			buf.WriteByte(b)
		case strings.IndexByte(".+-eE", b) != -1:
			buf.WriteByte(b)
		default:
			l.pushb(b)
			val, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return LexError
			}
			lval.val = val
			return Number
		}
	}
	return 0
}

var literal = map[string]interface{}{
	"true":  true,
	"false": false,
	"null":  nil,
}

func (l *lex) scanLiteral(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for b := l.nextb(); b != 0; b = l.nextb() {
		switch {
		case unicode.IsLetter(rune(b)):
			buf.WriteByte(b)
		default:
			l.pushb(b)
			val, ok := literal[buf.String()]
			if !ok {
				return LexError
			}
			lval.val = val
			return Literal
		}
	}
	return 0
}

func (l *lex) pushb(b byte) {
	l.pushed = b
}

func (l *lex) nextb() byte {
	if l.pushed != 0 {
		b := l.pushed
		l.pushed = 0
		return b
	}
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
