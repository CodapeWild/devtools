%{
package calculate

import (
  "fmt"
  "math/big"
)
%}

%union {
  num *big.Rat
}

%type <num> expr expr1 expr2 expr3

%token ADD SUB MUL DIV LEFT_BRACE RIGHT_BRACE

%token <num> NUM

%%

line: expr
    {
      fmt.Println($1.String())
    }
    ;

expr: expr1
    | ADD expr
    {
      $$ = $2
    }
    | SUB expr
    {
      $$ = $2.Neg($2)
    }
    ;

expr1: expr2
     | expr1 ADD expr2
     {
       $$ = $1.Add($1, $3)
     }
     | expr1 SUB expr2
     {
       $$ = $1.Sub($1, $3)
     }
     ;

expr2: expr3
     | expr2 MUL expr3
     {
       $$ = $1.Mul($1, $3)
     }
     | expr2 DIV expr3
     {
       $$ = $1.Quo($1, $3)
     }
     ;

expr3: NUM
     | LEFT_BRACE expr RIGHT_BRACE
     {
       $$ = $2
     }
     ;

%%