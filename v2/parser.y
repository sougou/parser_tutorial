%{
package jsonparser

func setResult(l yyLexer, v Result) {
  l.(*lex).result = v
}
%}

%union{
}

%token D

%start main

%%

main: phone
  {
    setResult(yylex, 0)
  }

phone:
  area part1 part2
| area '-' part1 '-' part2

area: D D D

part1: D D D

part2: D D D D
