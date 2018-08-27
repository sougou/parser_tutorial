%{
package jsonparser

func setResult(l yyLexer, v Result) {
  l.(*lex).result = v
}
%}

%union{
}

%token String Number True False Null

%start object

%%

object: '{' members '}'

members:
| pair
| members ',' pair

pair: String ':' value

array: '[' elements ']'

elements:
| value
| elements ',' value

value:
  String
| Number
| object
| array
| True
| False
| Null
