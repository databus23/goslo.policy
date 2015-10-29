%{
package policy

import (
  "fmt"
)

func(y yySymType) String() string {
  return fmt.Sprintf("{str:%s,...}", y.str)
}

%}

%union {
  f   rule 
  str string
  num int
  b   bool
}

%type <f> check expr

%token '(' ')' '@' '!' ':' 
%token and or not 
%left or 
%left and 
%left not 

%token variable, unquotedStr, constStr, number, boolean
%type <str> variable, unquotedStr, constStr
%type <b> boolean
%type <num> number 


%%

rule:
   {
     var f rule = func(c Context ) bool {return true }
     yylex.(*lexer).parseResult = f
   }
|
   expr 
   {
      yylex.(*lexer).parseResult = $1
   }

expr:
  not expr
  {
    f := $2
    $$ = func(c Context) bool { return !f(c) }
  }
|
  '(' expr ')'
  {
    f := $2
    $$ = func(c Context) bool { return f(c) } 
  } 
|
  expr or expr
  {
    left := $1
    right := $3
    $$ = func(c Context) bool { return left(c) || right(c) }
  }
|
  expr and expr
  {
    left := $1
    right := $3
    $$ = func(c Context) bool { return left(c) && right(c) } 
  }
|
  check
  {
    $$ = $1
  }

check:
  unquotedStr ':' unquotedStr
  {
    $$ = func(c Context) bool { return  c.genericCheck($1, $3, false) }
  }
|
  unquotedStr ':' constStr
  {
    $$ = func(c Context) bool { return  c.genericCheck($1, $3, false) }
  }
|
  unquotedStr ':' variable 
  {
    $$ = func(c Context) bool { return c.genericCheck($1, $3, true) }
  }
|
  constStr ':' variable
  {
    $$ = func(c Context) bool { return c.checkVariable( $3, $1 ) }
  }
|
  '@'
  {
    $$ = func(_ Context) bool { return true }
  }
|
  '!'
  {
    $$ = func(_ Context) bool { return false }
  }
%%

