%{
package jsonparser

import (
  "strconv"
)

func setResult(l yyLexer, v Result) {
  l.(*lex).result = v
}
%}

%union{
  ch byte
  bytes []byte
  value Result
}

%token <ch> Digit

%type <value> object number
%type <bytes> optSign digits fracOpt expOpt

%start object

%%

object: number
  {
    $$ = $1
    setResult(yylex, $1)
  }

number: optSign digits fracOpt expOpt
  {
    bval := append(append(append($1, $2...), $3...), $4...)
    val, err := strconv.ParseFloat(string(bval), 64)
    if err != nil {
      yylex.Error(err.Error())
      return 1
    }
    $$ = val
  }

optSign:
  {
    $$ = nil
  }
| '+'
  {
    $$ = []byte{'+'}
  }
| '-'
  {
    $$ = []byte{'-'}
  }

// TODO(sougou): first digit of a non-zero number cannot be zero.
// Add a rule to handle that case.
digits:
  Digit
  {
    $$ = []byte{$1}
  }
| digits Digit
  {
    $$ = append($1, $2)
  }

fracOpt:
  {
    $$ = nil
  }
| '.' digits
  {
    $$ = append([]byte{'.'}, $2...)
  }

expOpt:
  {
    $$ = nil
  }
| e optSign digits
  {
    $$ = append([]byte{'e'}, $2...)
    $$ = append($$, $3...)
  }

e:
  'e'
| 'E'
