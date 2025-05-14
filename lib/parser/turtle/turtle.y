%{
package turtle

import (
	"fmt"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"github.com/maartyman/rdfgo/interfaces"
	"strings"
)

type predicateObject struct {
	Predicate interfaces.ITerm
	Object	  interfaces.ITerm
}
%}

%union {
	term interfaces.ITerm
	literal interfaces.ILiteral
	str string
	quad interfaces.IQuad
	predObj struct {
		pairs []predicateObject
	}
	blankNodePO struct {
		node interfaces.IBlankNode
		triples []interfaces.IQuad
	}
	objectList []interfaces.ITerm
}

%token <literal> _NUMBER
%token <str> _PNAME _PVALUE _BNODE _NNODE _LITERAL _LANGTAG _SPARQLPREFIX _SPARQLBASE
%token _DOT _SEMICOLON _COMMA _PREFIX _BASE _ERROR _A _BRACKETOPEN _BRACKETCLOSE _TRUE _FALSE _LPAREN _RPAREN _DATATYPE

%type <term> subject object prefixed verb fullNNode
%type <predObj> predicateObjectList objectList predicateObjectTail predicateObjectPair
%type <blankNodePO> blankNodePropertyList collection
%type <literal> literal
%type <objectList> objectListItems

%%

turtleDoc:
	/* empty */
	| turtleDoc statement
	;

statement:
	triples
	| prefixDecl
	;

prefixDecl:
	_SPARQLPREFIX _PNAME fullNNode {
		yylex.(*lexer).prefixes[$2] = $3.GetValue()
	}
	| _PREFIX _PNAME fullNNode _DOT {
		yylex.(*lexer).prefixes[$2] = $3.GetValue()
	}
	| _BASE fullNNode _DOT {
		yylex.(*lexer).base = $2.GetValue()
	}
	| _SPARQLBASE fullNNode {
		yylex.(*lexer).base = $2.GetValue()
	}
	;

triples:
	blankNodePropertyList _DOT {
		for _, q := range $1.triples {
			yylex.(*lexer).output <- q
		}
	}
	| subject predicateObjectList _DOT {
		for _, po := range $2.pairs {
			quad, err := yylex.(*lexer).dataFactory.Quad($1, po.Predicate, po.Object, yylex.(*lexer).dataFactory.DefaultGraph())
			if err != nil {
				yylex.(*lexer).Error(fmt.Sprintf("error constructing quad: %v", err))
				return 1
			}
			yylex.(*lexer).output <- quad
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
		$$.pairs = []predicateObject{}
	}
	| _SEMICOLON predicateObjectPair predicateObjectTail {
		$$ = $2
		$$.pairs = append($$.pairs, $3.pairs...)
	}
	| _SEMICOLON predicateObjectTail {
		$$ = $2
	}
	;

verb:
	fullNNode { $$ = $1 }
	| prefixed { $$ = $1 }
	| _A { $$ = IRI.RDF.Type }

objectList:
	object {
		$$.pairs = []predicateObject{{Predicate: nil, Object: $1}}
	}
	| object _COMMA objectList {
		$$.pairs = append([]predicateObject{{Predicate: nil, Object: $1}}, $3.pairs...)
	}
	;

collection:
	_LPAREN _RPAREN {
		$$.node = IRI.RDF.Nil
		$$.triples = nil
	}
	| _LPAREN objectListItems _RPAREN {
		triples := []interfaces.IQuad{}
		var head, prev interfaces.IBlankNode
		prev = yylex.(*lexer).dataFactory.BlankNode("")
		for i, item := range $2 {
			curr := prev
			quad, err := yylex.(*lexer).dataFactory.Quad(curr, IRI.RDF.First, item, yylex.(*lexer).dataFactory.DefaultGraph())
			if err != nil {
				yylex.(*lexer).Error(fmt.Sprintf("error constructing quad: %v", err))
				return 1
			}
			triples = append(triples, quad)

			var rest interfaces.ITerm
			if i == len($2)-1 {
				rest = IRI.RDF.Nil
			} else {
				rest = yylex.(*lexer).dataFactory.BlankNode("")
			}
			quad, err = yylex.(*lexer).dataFactory.Quad(curr, IRI.RDF.Rest, rest, yylex.(*lexer).dataFactory.DefaultGraph())
			if err != nil {
				yylex.(*lexer).Error(fmt.Sprintf("error constructing quad: %v", err))
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
	| _BNODE { $$ = yylex.(*lexer).dataFactory.BlankNode("") }
	| prefixed { $$ = $1 }
	| blankNodePropertyList {
		for _, q := range $1.triples {
			yylex.(*lexer).output <- q
		}
		$$ = $1.node
	}
	| collection {
		for _, q := range $1.triples {
			yylex.(*lexer).output <- q
		}
		$$ = $1.node
	}
	;

object:
	fullNNode { $$ = $1 }
	| _BNODE { $$ = yylex.(*lexer).dataFactory.BlankNode("") }
	| prefixed { $$ = $1 }
	| literal { $$ = $1 }
	| blankNodePropertyList {
		for _, q := range $1.triples {
			yylex.(*lexer).output <- q
		}
		$$ = $1.node
	}
	| collection {
		for _, q := range $1.triples {
			yylex.(*lexer).output <- q
		}
		$$ = $1.node
	}
	;

literal:
	_TRUE { $$ = yylex.(*lexer).dataFactory.Literal("true", "", IRI.XSD.Boolean) }
	| _FALSE { $$ = yylex.(*lexer).dataFactory.Literal("false", "", IRI.XSD.Boolean) }
	| _LITERAL { $$ = yylex.(*lexer).dataFactory.Literal($1, "", nil) }
	| _LITERAL _LANGTAG { $$ = yylex.(*lexer).dataFactory.Literal($1, $2, nil) }
	| _LITERAL _DATATYPE fullNNode { $$ = yylex.(*lexer).dataFactory.Literal($1, "", $3) }
	| _LITERAL _DATATYPE prefixed { $$ = yylex.(*lexer).dataFactory.Literal($1, "", $3) }
	| _NUMBER { $$ = $1 }
	;

blankNodePropertyList:
	_BRACKETOPEN _BRACKETCLOSE {
		$$.node = yylex.(*lexer).dataFactory.BlankNode("")
		$$.triples = nil
	}
	| _BRACKETOPEN predicateObjectList _BRACKETCLOSE {
		bnode := yylex.(*lexer).dataFactory.BlankNode("")
		var quads []interfaces.IQuad
		for _, po := range $2.pairs {
			q, err := yylex.(*lexer).dataFactory.Quad(bnode, po.Predicate, po.Object, yylex.(*lexer).dataFactory.DefaultGraph())
			if err != nil {
				yylex.(*lexer).Error(fmt.Sprintf("blank node quad error: %v", err))
				return 1
			}
			quads = append(quads, q)
		}
		$$.node = bnode
		$$.triples = quads
	}

prefixed:
	_PNAME _PVALUE {
		prefix := yylex.(*lexer).prefixes[$1]
		if prefix == "" {
			yylex.(*lexer).Error(fmt.Sprintf("prefix not found: %s", $1))
			return 1
		}
		$$ = yylex.(*lexer).dataFactory.NamedNode(prefix + $2)
	}
	| _PNAME {
		prefix := yylex.(*lexer).prefixes[$1]
		if prefix == "" {
			yylex.(*lexer).Error(fmt.Sprintf("prefix not found: %s", $1))
			return 1
		}
		$$ = yylex.(*lexer).dataFactory.NamedNode(prefix)
	}
	;

fullNNode:
	_NNODE {
		if !strings.Contains($1, ":") && !strings.HasPrefix($1, "/") {
			base := yylex.(*lexer).base
			if base == "" {
				yylex.(*lexer).Error("@base not defined")
				return 1
			}
			$$ = yylex.(*lexer).dataFactory.NamedNode(base + $1)
		} else {
			$$ = yylex.(*lexer).dataFactory.NamedNode($1)
		}
	}

%%
