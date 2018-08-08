%{
package jsonparser

func setResult(l yyLexer, v Result) {
  l.(*lex).result = v
}
%}

%union{
}

%token Digit

%start object

%%

object: number

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
