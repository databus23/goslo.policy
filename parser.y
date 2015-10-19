%{
package policy

import (
  "fmt"
)

type check struct {
  key string
  match string
}

func(y yySymType) String() string {
  return fmt.Sprintf("{key:%s, match:%s; f(%q)}", y.check.key, y.check.match, y.f)
}

%}

%union {
  f   rule 
  check check
}

%type <f> check

%token '(' ')' '@' '!' 
%token and or not 
%left or 
%left and 
%left not 

%token const_check role_check rule_check http_check token_var_check token_const_check 
%type <check> const_check role_check rule_check http_check token_var_check token_const_check 


%%

rule:
   {
     var f rule = func(c Context ) bool {return true }
     yylex.(*Lexer).parseResult = f
   }
|
   check 
   {
      yylex.(*Lexer).parseResult = $1
   }

check:
  not check
  {
    f := $2
    $$ = func(c Context) bool { return !f(c) }
  }
|
  '(' check ')'
  {
    f := $2
    $$ = func(c Context) bool { return f(c) } 
  } 
|
  check or check
  {
    left := $1
    right := $3
    $$ = func(c Context) bool { return left(c) || right(c) }
  }
|
  check and check
  {
    left := $1
    right := $3
    $$ = func(c Context) bool { return left(c) && right(c) } 
  }
|
  rule_check
  {
    rules := yylex.(*Lexer).rules
    $$ = func(c Context) bool { r,ok := (*rules)[$1.match]; return ok && r(c) }
  }
|
  role_check
  {
    $$ = func(c Context) bool { return c.hasRole($1.match) }
  }
|
  const_check
  {
    $$ = func(c Context) bool { return c.matchVar($1.match, $1.key) }
  }
|
  token_const_check 
  {
    $$ = func(c Context) bool { return c.matchToken($1.key, $1.match) }
  }
|
  token_var_check
  {
    $$ = func(c Context) bool { return c.matchTokenAndVar($1.key, $1.match) }
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

