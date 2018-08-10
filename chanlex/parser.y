%{
package jsonparser

type pair struct {
  obj map[string]interface{}
  key string
  val interface{}
}

func setResult(l yyLexer, v map[string]interface{}) {
  l.(*lex).result = v
}
%}

%union{
  obj map[string]interface{}
  pair pair
  list []interface{}
  val interface{}
}

%token LexError
%token <val> String Number Literal

%type <obj> object members

%type <val> array
%type <list> elements
%type <pair> pair
%type <val> value

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
      $1.key: $1.val,
    }
  }
| members ',' pair
  {
    $$ = $1
    $$[$3.key] = $3.val
  }

pair: String ':' value
  {
    $$ = pair{key: $1.(string), val: $3}
  }

array: '[' elements ']'
  {
    $$ = $2
  }

elements:
  {
    $$ = []interface{}{}
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
| Number
| object
  {
    $$ = $1
  }
| array
| Literal
