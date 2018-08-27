%{
package jsonparser

func setResult(l yyLexer, v Result) {
  l.(*lex).result = v
}
%}

%union{
  result Result
  part string
  ch byte
}

%token <ch> D

%type <result> phone
%type <part> area part1 part2

%start main

%%

main: phone
  {
    setResult(yylex, $1)
  }

phone:
  area part1 part2
  {
    $$ = Result{area: $1, part1: $2, part2: $3}
  }
| area '-' part1 '-' part2
  {
    $$ = Result{area: $1, part1: $3, part2: $5}
  }

area: D D D
  {
    $$ = cat($1, $2, $3)
  }

part1: D D D
  {
    $$ = cat($1, $2, $3)
  }

part2: D D D D
  {
    $$ = cat($1, $2, $3, $4)
  }
