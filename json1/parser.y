%{
package jsonparser

func setResult(l yyLexer, v Result) {
  l.(*lex).result = v
}
%}

%union{
}

%token String Number Literal

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
| Literal
| object
| array
