%{
package jsonparser

func setResult(l yyLexer, v Result) {
  l.(*lex).result = v
}
%}

%union{
}

%token String Digit

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
| number
| object
| array
| true
| false
| null


//------------------------------------------------------
// The parts below are often handled by lex.

number: optSign digits fracOpt expOpt

optSign:
| '+'
| '-'

digits:
  Digit
| digits Digit

fracOpt:
| '.' digits

expOpt:
| e optSign digits

e:
  'e'
| 'E'

true: 't' 'r' 'u' 'e'

false: 'f' 'a' 'l' 's' 'e'

null: 'n' 'u' 'l' 'l'
