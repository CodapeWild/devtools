%{
package comparable

import (
  "log"
)
%}

%union {
  comp Comparable
}

%type <comp> expr expr1 expr2 expr3 expr4 expr5 expr6

%token SEMICOLON LEFT_CURLY_BRACKET RIGHT_CULY_BRACKET AND OR GREAT LITTLE GREAT_EQUAL LITTLE_EQUAL EQUAL NOT_EQUAL IN NOT_INT LEFT_BRACE RIGHT_BRACE

%token <comp> COMPARABLE

%%



%%