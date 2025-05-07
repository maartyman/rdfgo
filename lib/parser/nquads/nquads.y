%{
package nquads

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

%token <term> _NNODE _BNODE _DATATYPE
%token <str> _LITERALVALUE _LANGTAG
%token _DOT _ERROR

%type <term> subject predicate object graph_opt literal
%type <quad> triple_line

%%

document:
      /* empty */
    | document triple_line
    ;

triple_line:
    subject predicate object graph_opt _DOT {
        quad, err := NewQuad($1, $2, $3, $4)
        if err != nil {
            yylex.Error(fmt.Sprintf("Quad creation failed: %v", err))
            return 1
        }
        yylex.(*lexer).output <- quad
        $$ = quad
    }
    ;

subject:
      _NNODE  { $$ = $1 }
    | _BNODE  { $$ = $1 }
    ;

predicate:
      _NNODE  { $$ = $1 }
    ;

object:
      _NNODE     { $$ = $1 }
    | _BNODE     { $$ = $1 }
    | literal   { $$ = $1 }
    ;

graph_opt:
      /* empty */ { $$ = NewDefaultGraph() }
    | _NNODE       { $$ = $1 }
    | _BNODE       { $$ = $1 }
    ;

literal:
      _LITERALVALUE                      { $$ = NewLiteral($1, "", nil) }
    | _LITERALVALUE _LANGTAG            { $$ = NewLiteral($1, $2, nil) }
    | _LITERALVALUE _DATATYPE         { $$ = NewLiteral($1, "", $2) }
    ;

%%
