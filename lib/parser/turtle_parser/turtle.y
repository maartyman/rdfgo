%{
package turtle_parser

import (
	"fmt"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"github.com/maartyman/rdfgo/interfaces"
	"strings"
)

type PredicateObject struct {
	Predicate interfaces.ITerm
	Object    interfaces.ITerm
}
%}

%union {
	term interfaces.ITerm
	literal interfaces.ILiteral
	str string
	quad interfaces.IQuad
	predObj struct {
		pairs []PredicateObject
	}
	blankNodePO struct {
		node interfaces.IBlankNode
		triples []interfaces.IQuad
	}
	objectList []interfaces.ITerm
}

%token <literal> NUMBER
%token <str> PNAME PVALUE BNODE NNODE LITERAL LANGTAG DATATYPE SPARQLPREFIX SPARQLBASE
%token DOT SEMICOLON COMMA PREFIX BASE ERROR A BRACKETOPEN BRACKETCLOSE TRUE FALSE LPAREN RPAREN

%type <term> subject object prefixed verb fullNNode
%type <predObj> predicateObjectList objectList predicateObjectTail predicateObjectPair
%type <blankNodePO> blankNodePropertyList collection
%type <literal> literal
%type <objectList> objectListItems

%%

turtleDoc:
	statements
	;

statements:
	/* empty */
	| statements statement
	;

statement:
	triples
	| prefixDecl
	;

prefixDecl:
	SPARQLPREFIX PNAME NNODE {
		yylex.(*Lexer).prefixes[$2] = $3
	}
	| PREFIX PNAME NNODE DOT {
		yylex.(*Lexer).prefixes[$2] = $3
	}
	| BASE NNODE DOT {
		yylex.(*Lexer).base = $2
	}
	| SPARQLBASE NNODE {
		yylex.(*Lexer).base = $2
	}
	;

triples:
    blankNodePropertyList DOT {
	for _, q := range $1.triples {
	    yylex.(*Lexer).output <- q
	}
    }
    | subject predicateObjectList DOT {
        for _, po := range $2.pairs {
            quad, err := NewQuad($1, po.Predicate, po.Object, NewDefaultGraph())
            if err != nil {
                yylex.(*Lexer).Error(fmt.Sprintf("error constructing quad: %v", err))
                return 1
            }
            yylex.(*Lexer).output <- quad
        }
    }
    ;

predicateObjectList:
    predicateObjectPair predicateObjectTail {
        $$ = $1
        $$.pairs = append($$.pairs, $2.pairs...)
    }
    ;

predicateObjectPair:
    verb objectList {
        $$ = $2
        for i := range $$.pairs {
            $$.pairs[i].Predicate = $1
        }
    }
    ;

predicateObjectTail:
    /* empty */ {
    	$$.pairs = []PredicateObject{}
    }
    | SEMICOLON predicateObjectPair predicateObjectTail {
        $$ = $2
        $$.pairs = append($$.pairs, $3.pairs...)
    }
    | SEMICOLON predicateObjectTail {
        $$ = $2  // handle extra trailing semicolons like ;;
    }
    ;

verb:
	fullNNode { $$ = $1 }
	| prefixed { $$ = $1 }
	| A { $$ = NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#type") }

objectList:
    object {
        $$.pairs = []PredicateObject{{Predicate: nil, Object: $1}}
    }
    | object COMMA objectList {
        $$.pairs = append([]PredicateObject{{Predicate: nil, Object: $1}}, $3.pairs...)
    }
    ;

collection:
	LPAREN RPAREN {
		$$.node = NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil")
		$$.triples = nil
	}
	| LPAREN objectListItems RPAREN {
		triples := []interfaces.IQuad{}
		var head, prev interfaces.IBlankNode
		prev = yylex.(*Lexer).newBlankNode()
		for i, item := range $2 {
			curr := prev
			quad, err := NewQuad(curr, NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), item, NewDefaultGraph())
			if err != nil {
				yylex.(*Lexer).Error(fmt.Sprintf("error constructing quad: %v", err))
				return 1
			}
			triples = append(triples, quad)

			var rest interfaces.ITerm
			if i == len($2)-1 {
				rest = NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil")
			} else {
				rest = yylex.(*Lexer).newBlankNode()
			}
			quad, err = NewQuad(curr, NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), rest, NewDefaultGraph())
			if err != nil {
				yylex.(*Lexer).Error(fmt.Sprintf("error constructing quad: %v", err))
				return 1
			}
			triples = append(triples, quad)

			if i == 0 {
				head = curr
			}
			prev = rest
		}
		$$.node = head
		$$.triples = triples
	}

objectListItems:
	object {
		$$ = []interfaces.ITerm{$1}
	}
	| object objectListItems {
		$$ = append([]interfaces.ITerm{$1}, $2...)
	}

subject:
	fullNNode { $$ = $1 }
	| BNODE { $$ = NewBlankNode($1) }
	| prefixed { $$ = $1 }
	| blankNodePropertyList {
                for _, q := range $1.triples {
			yylex.(*Lexer).output <- q
                }
                $$ = $1.node
            }
	| collection          {
		for _, q := range $1.triples {
		    yylex.(*Lexer).output <- q
		}
		$$ = $1.node
	    }
	;

object:
	fullNNode { $$ = $1 }
	| BNODE { $$ = NewBlankNode($1) }
	| prefixed { $$ = $1 }
	| literal { $$ = $1 }
	| blankNodePropertyList {
                for _, q := range $1.triples {
			yylex.(*Lexer).output <- q
                }
                $$ = $1.node
            }
	| collection {
		for _, q := range $1.triples {
			yylex.(*Lexer).output <- q
		}
		$$ = $1.node
	}
	;

literal:
	TRUE { $$ = NewLiteral("true", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#boolean")) }
	| FALSE { $$ = NewLiteral("false", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#boolean")) }
	| LITERAL { $$ = NewLiteral($1, "", nil) }
	| LITERAL LANGTAG { $$ = NewLiteral($1, $2, nil) }
	| LITERAL DATATYPE { $$ = NewLiteral($1, "", NewNamedNode($2)) }
	| NUMBER { $$ = $1 }
	;

blankNodePropertyList:
    BRACKETOPEN BRACKETCLOSE {
	$$.node = yylex.(*Lexer).newBlankNode()
	$$.triples = nil
    }
    | BRACKETOPEN predicateObjectList BRACKETCLOSE {
	bnode := yylex.(*Lexer).newBlankNode()
	var quads []interfaces.IQuad
	for _, po := range $2.pairs {
	    q, err := NewQuad(bnode, po.Predicate, po.Object, NewDefaultGraph())
	    if err != nil {
		yylex.(*Lexer).Error(fmt.Sprintf("blank node quad error: %v", err))
		return 1
	    }
	    quads = append(quads, q)
	}
	$$.node = bnode
	$$.triples = quads
    }

prefixed:
	PNAME PVALUE {
		prefix := yylex.(*Lexer).prefixes[$1]
		if prefix == "" {
			yylex.(*Lexer).Error(fmt.Sprintf("prefix not found: %s", $1))
			return 1
		}
		$$ = NewNamedNode(prefix + $2)
	}
	| PNAME {
		prefix := yylex.(*Lexer).prefixes[$1]
		if prefix == "" {
			yylex.(*Lexer).Error(fmt.Sprintf("prefix not found: %s", $1))
			return 1
		}
		$$ = NewNamedNode(prefix)
	}
	;

fullNNode:
	NNODE {
		if !strings.Contains($1, ":") && !strings.HasPrefix($1, "/") {
			base := yylex.(*Lexer).base
			if base == "" {
				yylex.(*Lexer).Error("@base not defined")
				return 1
			}
			$$ = NewNamedNode(base + $1)
		} else {
			$$ = NewNamedNode($1)
		}
	}

%%
