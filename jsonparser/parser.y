%{
package jsonparser

import (
  "strconv"
)

type pair struct {
  str string
  val interface{}
}

func setResult(l yyLexer, v map[string]interface{}) {
  l.(*lex).result = v
}
%}

%union{
  obj map[string]interface{}
  list []interface{}
  value interface{}
  pair pair
  ch byte
  bytes []byte
  str string
}

%token LexError
%token <str> String
%token <ch> Digit

%type <obj> object members
%type <pair> pair
%type <list> array elements
%type <value> value number
%type <bytes> signOpt digits fracOpt expOpt
%type <value> true false null

%start object

%%

object: '{' members '}'
  {
    $$ = $2
    setResult(yylex, $$)
  }

members:
  {
    $$ = map[string]interface{}{}
  }
| pair
  {
    $$ = map[string]interface{}{
      $1.str: $1.val,
    }
  }
| members ',' pair
  {
    $1[$3.str] = $3.val
    $$ = $1
  }

pair: String ':' value
  {
    $$ = pair{str: $1, val: $3}
  }

array: '[' elements ']'
  {
    $$ = $2
  }

elements:
  {
    $$ = nil
  }
| value
  {
    $$ = []interface{}{$1}
  }
| elements ',' value
  {
    $$ = append($1, $3)
  }

value:
  String
  {
    $$ = $1
  }
| number
| object
  {
    $$ = $1
  }
| array
  {
    $$ = $1
  }
| true
| false
| null


//------------------------------------------------------
// The parts below are often handled by lex.

number: signOpt digits fracOpt expOpt
  {
    bval := append(append(append($1, $2...), $3...), $4...)
    val, err := strconv.ParseFloat(string(bval), 64)
    if err != nil {
      yylex.Error(err.Error())
      return 1
    }
    $$ = val
  }

signOpt:
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
| e signOpt digits
  {
    $$ = append([]byte{'e'}, $2...)
    $$ = append($$, $3...)
  }

e:
  'e'
| 'E'

true: 't' 'r' 'u' 'e'
  {
    $$ = true
  }

false: 'f' 'a' 'l' 's' 'e'
  {
    $$ = false
  }

null: 'n' 'u' 'l' 'l'
  {
    $$ = nil
  }
