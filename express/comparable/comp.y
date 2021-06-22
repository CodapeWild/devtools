%{
package comparable

import (
  "fmt"
)
%}

%union {
  comp Comparable
  comps []Comparable
}

%type <comp> expr expr1 expr2 expr3 expr4 expr5 expr6 expr7

%type <comps> params

%type <comp> param

%token SEMICOLON LEFT_CURLY_BRACKET RIGHT_CULY_BRACKET COMMA AND OR GREAT LITTLE GREAT_EQUAL LITTLE_EQUAL EQUAL NOT_EQUAL IN NOT_IN MATCH LEFT_BRACE RIGHT_BRACE

%token <comp> COMPARABLE

%%

line: expr
    {
      if b,err := $1.boolean(); err != nil {
        fmt.Println(err.Error())
      }else {
        if b {
          fmt.Println("true")
        }else {
          fmt.Println("false")
        }
      }
    }
    ;

expr: expr1
    | SEMICOLON expr
    {
      $$ = $2
    }
    ;

expr1: expr2
     | expr1 SEMICOLON expr2
     {
       $$ = $1.Or($3)
     }
     ;

expr2: expr3
     | LEFT_CURLY_BRACKET expr2 RIGHT_CULY_BRACKET
     {
       $$ = $2
     }
     ;

expr3: expr4
     | COMMA expr3
     {
       $$ = $2
     }
     ;

expr4: expr5
     | expr4 COMMA expr5
     {
       $$ = $1.And($3)
     }
     ;

expr5: expr6
     | expr5 AND expr6
     {
       $$ = $1.And($3)
     }
     | expr5 OR expr6
     {
       $$ = $1.Or($3)
     }
     ;

expr6: expr7
     | expr6 GREAT expr7
     {
       $$ = $1.GreatThan($3)
     }
     | expr6 LITTLE expr7
     {
       $$ = $1.LittleThan($3)
     }
     | expr6 GREAT_EQUAL expr7
     {
       $$ = $1.GreatEqualThan($3)
     }
     | expr6 LITTLE_EQUAL expr7
     {
       $$ = $1.LittleEqualThan($3)
     }
     | expr6 EQUAL expr7
     {
       $$ = $1.Equal($3)
     }
     | expr6 NOT_EQUAL expr7
     {
       $$ = $1.NotEqual($3)
     }
     | expr6 IN LEFT_BRACE params RIGHT_BRACE
     {
       $$ = $1.In(params...)
     }
     | expr6 NOT_IN LEFT_BRACE params RIGHT_BRACE
     {
       $$ = $1.NotIn(params...)
     }
     | MATCH LEFT_BRACE expr6 COMMA expr7 RIGHT_BRACE
     {
       $$ = Match($3.string(), $5.string())
     }
     ;

params: param
      {
        $$ = []Comparable{ $1 }
      }
      | params COMMA param
      {
        $$ = append($$, $3)
      }
      ;

param: COMPARABLE
      {
        $$ = $1
      }
      ;

expr7: COMPARABLE
     | LEFT_BRACE expr3 RIGHT_BRACE
     {
       $$ = $2
     }
     ;

%%