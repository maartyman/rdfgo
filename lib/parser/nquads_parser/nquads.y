%{
package nquads_parser

import (
    "fmt"
    . "github.com/maartyman/rdfgo/lib/data_model"
    "github.com/maartyman/rdfgo/interfaces"
)
%}

%union {
    str   string
    term  interfaces.ITerm
    quad  interfaces.IQuad
}

%token <term> NNODE BNODE DATATYPE
%token <str> LITERALVALUE LANGTAG
%token DOT ERROR

%type <term> subject predicate object graph_opt literal
%type <quad> triple_line

%%

document:
      /* empty */
    | document triple_line
    ;

triple_line:
    subject predicate object graph_opt DOT {
        quad, err := NewQuad($1, $2, $3, $4)
        if err != nil {
            yylex.Error(fmt.Sprintf("Quad creation failed: %v", err))
            return 1
        }
        yylex.(*Lexer).output <- quad
        $$ = quad
    }
    ;

subject:
      NNODE  { $$ = $1 }
    | BNODE  { $$ = $1 }
    ;

predicate:
      NNODE  { $$ = $1 }
    ;

object:
      NNODE     { $$ = $1 }
    | BNODE     { $$ = $1 }
    | literal   { $$ = $1 }
    ;

graph_opt:
      /* empty */ { $$ = NewDefaultGraph() }
    | NNODE       { $$ = $1 }
    | BNODE       { $$ = $1 }
    ;

literal:
      LITERALVALUE                      { $$ = NewLiteral($1, "", nil) }
    | LITERALVALUE LANGTAG            { $$ = NewLiteral($1, $2, nil) }
    | LITERALVALUE DATATYPE         { $$ = NewLiteral($1, "", $2) }
    ;

%%
