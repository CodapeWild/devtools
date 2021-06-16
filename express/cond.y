%{

package express

import (
  "fmt"
  "math/rat"
)

%}

%union{
  num *big.Rat
}



%type <num> expr expr1 expr2 expr3

%token '>' '<' '=' '&' '|' '(' ')'

%token <num> NUM

%%

cond: expr
    {

    }
    ;

expr: expr1
    | expr '>' expr1
    {

    }
    | expr '<' expr1
    {

    }
    | expr '=' '=' expr1
    {

    }
    | expr '&' '&' expr1
    {

    }
    | expr '|' '|' expr1
    {

    }
    ;

expr1: NUM
     | '(' expr ')'
     {
       $$ = $2
     }
     ;

%%