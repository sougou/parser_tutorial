%{
package jsonparser

func setResult(l yyLexer, v Result) {
  l.(*lex).result = v
}
%}

%union{
}

%start main

%%

main:
  {
    setResult(yylex, 0)
  }
