package turtle_parser

import (
	"errors"
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

type TurtleTestCase struct {
	name        string
	input       string
	expectError bool
	expectQuads []interfaces.IQuad
}

func quadsListToString(quads []interfaces.IQuad) string {
	var sb strings.Builder
	sb.WriteString("[\n")
	for _, quad := range quads {
		sb.WriteString(fmt.Sprintf("\t%s,\n", quad.ToString()))
	}
	sb.WriteString("]")
	return sb.String()
}

func q(s interfaces.ITerm, p interfaces.ITerm, o interfaces.ITerm, g interfaces.ITerm) interfaces.IQuad {
	quad, _ := NewQuad(
		s,
		p,
		o,
		g,
	)
	return quad
}

func nn(val string) interfaces.INamedNode {
	namedNode := NewNamedNode("http://example.org/" + val)
	return namedNode
}

func bn(val string) interfaces.IBlankNode {
	blankNode := NewBlankNode(val)
	return blankNode
}

func l(val string) interfaces.ILiteral {
	literal := NewLiteral(val, "", nil)
	return literal
}

func lit(val string, lang string, dt interfaces.INamedNode) interfaces.ILiteral {
	literal := NewLiteral(val, lang, dt)
	return literal
}

func TestOutput(t *testing.T) {
	tests := []TurtleTestCase{
		{
			name:  "Basic triple",
			input: `<http://example.org/s> <http://example.org/p> <http://example.org/o> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("o"), nil),
			},
		},
		{
			name:  "Triple with literal object",
			input: `<http://example.org/s> <http://example.org/p> "hello" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("hello"), nil),
			},
		},
		{
			name:  "Blank node subject",
			input: `_:b0 <http://example.org/p> "blank" .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), nn("p"), l("blank"), nil),
			},
		},
		{
			name:        "Invalid input",
			input:       `this is not valid turtle`,
			expectError: true,
		},
		{
			name: "Prefix usage in triple",
			input: `
@prefix ex: <http://example.org/> .
ex<http://example.org/s> ex<http://example.org/p> ex:o .
`,
			expectError: true,
		},
		{
			name: "Default prefix usage in triple",
			input: `
@prefix : <http://example.org/> .
<http://example.org/s> <http://example.org/p> :o .
`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("o"), nil),
			},
		},
		{
			name: "Default prefix usage in triple",
			input: `
@prefix : <http://example.org/> .
<http://example.org/s> <http://example.org/p> :s .
`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("s"), nil),
			},
		},
		{
			name:  "supports a as shorthand for rdf:type",
			input: "<http://example.org/s> a <http://example.org/o> .",
			expectQuads: []interfaces.IQuad{
				q(nn("s"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"), nn("o"), nil),
			},
		},
		{
			name:  "Multiple predicates with ;",
			input: `<http://example.org/s> <http://example.org/p1> "v1" ; <http://example.org/p2> "v2" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p1"), l("v1"), nil),
				q(nn("s"), nn("p2"), l("v2"), nil),
			},
		},
		{
			name:  "Multiple objects with ,",
			input: `<http://example.org/s> <http://example.org/p> "v1", "v2", "v3" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("v1"), nil),
				q(nn("s"), nn("p"), l("v2"), nil),
				q(nn("s"), nn("p"), l("v3"), nil),
			},
		},
		{
			name:  "Mixed ; and ,",
			input: `<http://example.org/s> <http://example.org/p1> "v1", "v2" ; <http://example.org/p2> "v3", "v4" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p1"), l("v1"), nil),
				q(nn("s"), nn("p1"), l("v2"), nil),
				q(nn("s"), nn("p2"), l("v3"), nil),
				q(nn("s"), nn("p2"), l("v4"), nil),
			},
		},
		{
			name:  "Single predicate, one object",
			input: `<http://example.org/s> <http://example.org/p> "v1" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("v1"), nil),
			},
		},
		{
			name: "Base IRI expands subject",
			input: `
			@base <http://example.org/> .
			<foo> <http://example.org/p> <http://example.org/o> .
		`,
			expectQuads: []interfaces.IQuad{
				q(nn("foo"), nn("p"), nn("o"), NewDefaultGraph()),
			},
		},
		{
			name: "Base IRI expands predicate",
			input: `
			@base <http://example.org/> .
			<http://example.org/s> <bar> <http://example.org/o> .
		`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("bar"), nn("o"), NewDefaultGraph()),
			},
		},
		{
			name: "Base IRI expands object",
			input: `
			@base <http://example.org/> .
			<http://example.org/s> <http://example.org/p> <baz> .
		`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("baz"), NewDefaultGraph()),
			},
		},
		{
			name: "Base IRI with nested path",
			input: `
			@base <http://example.org/base/> .
			<thing> <http://example.org/p> <target> .
		`,
			expectQuads: []interfaces.IQuad{
				q(NewNamedNode("http://example.org/base/thing"), nn("p"), NewNamedNode("http://example.org/base/target"), NewDefaultGraph()),
			},
		},
		{
			name: "Relative and absolute mix",
			input: `
			@base <http://example.org/base/> .
			<local> <http://example.org/p> <http://example.org/o> .
		`,
			expectQuads: []interfaces.IQuad{
				q(NewNamedNode("http://example.org/base/local"), nn("p"), nn("o"), NewDefaultGraph()),
			},
		},
		{
			name:  "Simple blank node property list",
			input: `@base <http://example.org/> . <s> <p> [ <q> "val" ] .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), nn("q"), l("val"), nil),
				q(nn("s"), nn("p"), bn("b0"), nil),
			},
		},
		{
			name:  "Blank node property list with multiple predicates",
			input: `@base <http://example.org/> . <s> <p> [ <q1> "v1" ; <q2> "v2" ] .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), nn("q1"), l("v1"), nil),
				q(bn("b0"), nn("q2"), l("v2"), nil),
				q(nn("s"), nn("p"), bn("b0"), nil),
			},
		},
		{
			name:  "Nested blank node property list",
			input: `@base <http://example.org/> . <s> <p> [ <q> [ <r> "val" ] ] .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), nn("r"), l("val"), nil),
				q(bn("b1"), nn("q"), bn("b0"), nil),
				q(nn("s"), nn("p"), bn("b1"), nil),
			},
		},
		{
			name:  "Nested blank node property list",
			input: `@base <http://example.org/> . <s><p>[<q>[<r>"val"]].`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), nn("r"), l("val"), nil),
				q(bn("b1"), nn("q"), bn("b0"), nil),
				q(nn("s"), nn("p"), bn("b1"), nil),
			},
		},
		{
			name:  "Multiple blank node lists as objects",
			input: `@base <http://example.org/> . <s> <p> [ <q> "v1" ], [ <q> "v2" ] .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), nn("q"), l("v1"), nil),
				q(bn("b1"), nn("q"), l("v2"), nil),
				q(nn("s"), nn("p"), bn("b0"), nil),
				q(nn("s"), nn("p"), bn("b1"), nil),
			},
		},
		{
			name:        "Invalid: unterminated blank node list",
			input:       `<http://example.org/s> <http://example.org/p> [ <http://example.org/q> "v1" .`,
			expectError: true,
		},
		{
			name: "Multiline literal with triple double quotes",
			input: `<http://example.org/s> <http://example.org/p> """This is
a multiline
string.""" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("This is\na multiline\nstring."), nil),
			},
		},
		{
			name: "Multiline literal with triple single quotes",
			input: `<http://example.org/s> <http://example.org/p> '''Another
multiline
example.''' .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("Another\nmultiline\nexample."), nil),
			},
		},
		{
			name:  "Literal double quotes triple with escape",
			input: `<http://example.org/s> <http://example.org/p> """\"""" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("\""), nil),
			},
		},
		{
			name:  "Literal single quotes triple with escape",
			input: `<http://example.org/s> <http://example.org/p> '''\'''' .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("'"), nil),
			},
		},
		{
			name: "Multiline literal double quotes triple with escape",
			input: `<http://example.org/s> <http://example.org/p> """
\"
""" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("\n\"\n"), nil),
			},
		},
		{
			name: "Multiline literal single quotes triple with escape",
			input: `<http://example.org/s> <http://example.org/p> '''
\'
''' .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("\n'\n"), nil),
			},
		},
		{
			name: "Multiline literal with a comment",
			input: `<http://example.org/s> <http://example.org/p> """This looks like # a comment
but it is not. # this isn't either""" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("This looks like # a comment\nbut it is not. # this isn't either"), nil),
			},
		},
		{
			name:  "Two multiline literals on one line",
			input: `<http://example.org/s> <http://example.org/p> """literal1""" . <http://example.org/s> <http://example.org/p> """literal2""" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("literal1"), nil),
				q(nn("s"), nn("p"), l("literal2"), nil),
			},
		},
		{
			name:  "Integer literal",
			input: `<http://example.org/s> <http://example.org/age> 42 .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("age"), lit("42", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#integer")), nil),
			},
		},
		{
			name:  "Decimal literal",
			input: `<http://example.org/s> <http://example.org/price> 19.99 .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("price"), lit("19.99", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#decimal")), nil),
			},
		},
		{
			name:  "Double literal (exponent)",
			input: `<http://example.org/s> <http://example.org/measurement> 6.022e23 .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("measurement"), lit("6.022e23", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#double")), nil),
			},
		},
		{
			name:  "Boolean literal true",
			input: `<http://example.org/s> <http://example.org/enabled> true .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("enabled"), lit("true", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#boolean")), nil),
			},
		},
		{
			name:  "Boolean literal false",
			input: `<http://example.org/s> <http://example.org/enabled> false .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("enabled"), lit("false", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#boolean")), nil),
			},
		},
		{
			name:  "Plain string literal",
			input: `<http://example.org/s> <http://example.org/p> "hello world" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("hello world"), nil),
			},
		},
		{
			name:  "Literal with language tag",
			input: `<http://example.org/s> <http://example.org/p> "hello"@en .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), lit("hello", "en", nil), nil),
			},
		},
		{
			name:  "Literal with language tag (regional)",
			input: `<http://example.org/s> <http://example.org/p> "bonjour"@fr-FR .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), lit("bonjour", "fr-FR", nil), nil),
			},
		},
		{
			name:  "Literal with xsd<http://example.org/s>tring",
			input: `<http://example.org/s> <http://example.org/p> "typed"^^<http://www.w3.org/2001/XMLSchema#string> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), lit("typed", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#string")), nil),
			},
		},
		{
			name:  "Literal with integer datatype",
			input: `<http://example.org/s> <http://example.org/p> "123"^^<http://www.w3.org/2001/XMLSchema#integer> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), lit("123", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#integer")), nil),
			},
		},
		{
			name:  "Literal with decimal datatype",
			input: `<http://example.org/s> <http://example.org/p> "3.14"^^<http://www.w3.org/2001/XMLSchema#decimal> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), lit("3.14", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#decimal")), nil),
			},
		},
		{
			name:  "Literal with boolean true",
			input: `<http://example.org/s> <http://example.org/p> "true"^^<http://www.w3.org/2001/XMLSchema#boolean> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), lit("true", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#boolean")), nil),
			},
		},
		{
			name:  "Literal with boolean false",
			input: `<http://example.org/s> <http://example.org/p> "false"^^<http://www.w3.org/2001/XMLSchema#boolean> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), lit("false", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#boolean")), nil),
			},
		},
		{
			name:  "Literal with escaped characters",
			input: `<http://example.org/s> <http://example.org/p> "Line1\nLine2\tTabbed" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("Line1\nLine2\tTabbed"), nil),
			},
		},
		{
			name:  "Literal with unicode escape",
			input: `<http://example.org/s> <http://example.org/p> "snowman: \u2603" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("snowman: ☃"), nil),
			},
		},
		{
			name:  "Empty collection",
			input: `<http://example.org/s> <http://example.org/p> () .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"), nil),
			},
		},
		{
			name:  "Single item collection",
			input: `<http://example.org/s> <http://example.org/p> ("a") .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), l("a"), nil),
				q(bn("b0"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"), nil),
				q(nn("s"), nn("p"), bn("b0"), nil),
			},
		},
		{
			name:  "Three item collection",
			input: `<http://example.org/s> <http://example.org/p> ("a" "b" "c") .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), l("a"), nil),
				q(bn("b0"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), bn("b1"), nil),
				q(bn("b1"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), l("b"), nil),
				q(bn("b1"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), bn("b2"), nil),
				q(bn("b2"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), l("c"), nil),
				q(bn("b2"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"), nil),
				q(nn("s"), nn("p"), bn("b0"), nil),
			},
		},
		{
			name: "Simple collection of literals",
			input: `
@prefix : <http://example.org/stuff/1.0/> .
:a :b ( "apple" "banana" ) .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), l("apple"), nil),
				q(bn("b0"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), bn("b1"), nil),
				q(bn("b1"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), l("banana"), nil),
				q(bn("b1"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"), nil),
				q(NewNamedNode("http://example.org/stuff/1.0/a"), NewNamedNode("http://example.org/stuff/1.0/b"), bn("b0"), nil),
			},
		},
		{
			name: "Multiline and escaped literal",
			input: `
@prefix : <http://example.org/stuff/1.0/> .

:a :b "The first line\nThe second line\n  more" .
:a :b """The first line
The second line
  more""" .`,
			expectQuads: []interfaces.IQuad{
				q(NewNamedNode("http://example.org/stuff/1.0/a"), NewNamedNode("http://example.org/stuff/1.0/b"), l("The first line\nThe second line\n  more"), nil),
				q(NewNamedNode("http://example.org/stuff/1.0/a"), NewNamedNode("http://example.org/stuff/1.0/b"), l("The first line\nThe second line\n  more"), nil),
			},
		},
		{
			name:  "Collection as subject",
			input: `(1 2.0 3E1) <http://example.org/p> "w" .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), lit("1", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#integer")), nil),
				q(bn("b0"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), bn("b1"), nil),
				q(bn("b1"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), lit("2.0", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#decimal")), nil),
				q(bn("b1"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), bn("b2"), nil),
				q(bn("b2"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), lit("3E1", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#double")), nil),
				q(bn("b2"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"), nil),
				q(bn("b0"), nn("p"), l("w"), nil),
			},
		},
		{
			name:  "Nested collection and blank node",
			input: `(1 [<http://example.org/p> <http://example.org/q>] (2)) <http://example.org/p2> <http://example.org/q2> .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), nn("p"), nn("q"), nil),
				q(bn("b1"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), lit("2", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#integer")), nil),
				q(bn("b1"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"), nil),
				q(bn("b2"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), lit("1", "", NewNamedNode("http://www.w3.org/2001/XMLSchema#integer")), nil),
				q(bn("b2"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), bn("b3"), nil),
				q(bn("b3"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), bn("b0"), nil),
				q(bn("b3"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), bn("b4"), nil),
				q(bn("b4"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"), bn("b1"), nil),
				q(bn("b4"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"), NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"), nil),
				q(bn("b2"), nn("p2"), nn("q2"), nil),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outChan, errChan := ParseTurtle(strings.NewReader(tc.input), nil)

			var got []interfaces.IQuad
			for quad := range outChan {
				got = append(got, quad)
			}

			err, ok := <-errChan
			if tc.expectError {
				if err == nil && ok {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(got) != len(tc.expectQuads) {
				t.Errorf("Expected %d quads, got %d", len(tc.expectQuads), len(got))
				t.Errorf("Expected: %s", quadsListToString(tc.expectQuads))
				t.Errorf("Got: %s", quadsListToString(got))
				return
			}

			for i, want := range tc.expectQuads {
				if i >= len(got) {
					t.Errorf("Missing quad: %s", want)
					continue
				}
				if !want.Equals(got[i]) {
					t.Errorf(
						"Mismatch at quad %d:\n  expected: %s\n       got: %s",
						i,
						want.ToString(),
						got[i].ToString(),
					)
				}
			}
		})
	}
}

func TestSpec(t *testing.T) {
	tests := []TurtleTestCase{
		{
			name:        "IRI_subject",
			input:       "spec_tests/IRI_subject.ttl",
			expectError: false,
		},
		{
			name:        "IRI_with_four_digit_numeric_escape",
			input:       "spec_tests/IRI_with_four_digit_numeric_escape.ttl",
			expectError: false,
		},
		{
			name:        "IRI_with_eight_digit_numeric_escape",
			input:       "spec_tests/IRI_with_eight_digit_numeric_escape.ttl",
			expectError: false,
		},
		{
			name:        "IRI_with_all_punctuation",
			input:       "spec_tests/IRI_with_all_punctuation.ttl",
			expectError: false,
		},
		{
			name:        "bareword_a_predicate",
			input:       "spec_tests/bareword_a_predicate.ttl",
			expectError: false,
		},
		{
			name:        "old_style_prefix",
			input:       "spec_tests/old_style_prefix.ttl",
			expectError: false,
		},
		{
			name:        "SPARQL_style_prefix",
			input:       "spec_tests/SPARQL_style_prefix.ttl",
			expectError: false,
		},
		{
			name:        "prefixed_IRI_predicate",
			input:       "spec_tests/prefixed_IRI_predicate.ttl",
			expectError: false,
		},
		{
			name:        "prefixed_IRI_object",
			input:       "spec_tests/prefixed_IRI_object.ttl",
			expectError: false,
		},
		{
			name:        "prefix_only_IRI",
			input:       "spec_tests/prefix_only_IRI.ttl",
			expectError: false,
		},
		{
			name:        "prefix_with_PN_CHARS_BASE_character_boundaries",
			input:       "spec_tests/prefix_with_PN_CHARS_BASE_character_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "prefix_with_non_leading_extras",
			input:       "spec_tests/prefix_with_non_leading_extras.ttl",
			expectError: false,
		},
		{
			name:        "default_namespace_IRI",
			input:       "spec_tests/default_namespace_IRI.ttl",
			expectError: false,
		},
		{
			name:        "prefix_reassigned_and_used",
			input:       "spec_tests/prefix_reassigned_and_used.ttl",
			expectError: false,
		},
		{
			name:        "reserved_escaped_localName",
			input:       "spec_tests/reserved_escaped_localName.ttl",
			expectError: false,
		},
		{
			name:        "percent_escaped_localName",
			input:       "spec_tests/percent_escaped_localName.ttl",
			expectError: false,
		},
		{
			name:        "HYPHEN_MINUS_in_localName",
			input:       "spec_tests/HYPHEN_MINUS_in_localName.ttl",
			expectError: false,
		},
		{
			name:        "underscore_in_localName",
			input:       "spec_tests/underscore_in_localName.ttl",
			expectError: false,
		},
		{
			name:        "localname_with_COLON",
			input:       "spec_tests/localname_with_COLON.ttl",
			expectError: false,
		},
		{
			name:        "localName_with_assigned_nfc_bmp_PN_CHARS_BASE_character_boundaries",
			input:       "spec_tests/localName_with_assigned_nfc_bmp_PN_CHARS_BASE_character_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "localName_with_assigned_nfc_PN_CHARS_BASE_character_boundaries",
			input:       "spec_tests/localName_with_assigned_nfc_PN_CHARS_BASE_character_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "localName_with_nfc_PN_CHARS_BASE_character_boundaries",
			input:       "spec_tests/localName_with_nfc_PN_CHARS_BASE_character_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "localName_with_leading_underscore",
			input:       "spec_tests/localName_with_leading_underscore.ttl",
			expectError: false,
		},
		{
			name:        "localName_with_leading_digit",
			input:       "spec_tests/localName_with_leading_digit.ttl",
			expectError: false,
		},
		{
			name:        "localName_with_non_leading_extras",
			input:       "spec_tests/localName_with_non_leading_extras.ttl",
			expectError: false,
		},
		{
			name:        "old_style_base",
			input:       "spec_tests/old_style_base.ttl",
			expectError: false,
		},
		{
			name:        "SPARQL_style_base",
			input:       "spec_tests/SPARQL_style_base.ttl",
			expectError: false,
		},
		{
			name:        "labeled_blank_node_subject",
			input:       "spec_tests/labeled_blank_node_subject.ttl",
			expectError: false,
		},
		{
			name:        "labeled_blank_node_object",
			input:       "spec_tests/labeled_blank_node_object.ttl",
			expectError: false,
		},
		{
			name:        "labeled_blank_node_with_PN_CHARS_BASE_character_boundaries",
			input:       "spec_tests/labeled_blank_node_with_PN_CHARS_BASE_character_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "labeled_blank_node_with_leading_underscore",
			input:       "spec_tests/labeled_blank_node_with_leading_underscore.ttl",
			expectError: false,
		},
		{
			name:        "labeled_blank_node_with_leading_digit",
			input:       "spec_tests/labeled_blank_node_with_leading_digit.ttl",
			expectError: false,
		},
		{
			name:        "labeled_blank_node_with_non_leading_extras",
			input:       "spec_tests/labeled_blank_node_with_non_leading_extras.ttl",
			expectError: false,
		},
		{
			name:        "anonymous_blank_node_subject",
			input:       "spec_tests/anonymous_blank_node_subject.ttl",
			expectError: false,
		},
		{
			name:        "anonymous_blank_node_object",
			input:       "spec_tests/anonymous_blank_node_object.ttl",
			expectError: false,
		},
		{
			name:        "sole_blankNodePropertyList",
			input:       "spec_tests/sole_blankNodePropertyList.ttl",
			expectError: false,
		},
		{
			name:        "blankNodePropertyList_as_subject",
			input:       "spec_tests/blankNodePropertyList_as_subject.ttl",
			expectError: false,
		},
		{
			name:        "blankNodePropertyList_as_object",
			input:       "spec_tests/blankNodePropertyList_as_object.ttl",
			expectError: false,
		},
		{
			name:        "blankNodePropertyList_as_object_containing_objectList",
			input:       "spec_tests/blankNodePropertyList_as_object_containing_objectList.ttl",
			expectError: false,
		},
		{
			name:        "blankNodePropertyList_as_object_containing_objectList_of_two_objects",
			input:       "spec_tests/blankNodePropertyList_as_object_containing_objectList_of_two_objects.ttl",
			expectError: false,
		},
		{
			name:        "blankNodePropertyList_with_multiple_triples",
			input:       "spec_tests/blankNodePropertyList_with_multiple_triples.ttl",
			expectError: false,
		},
		{
			name:        "nested_blankNodePropertyLists",
			input:       "spec_tests/nested_blankNodePropertyLists.ttl",
			expectError: false,
		},
		{
			name:        "blankNodePropertyList_containing_collection",
			input:       "spec_tests/blankNodePropertyList_containing_collection.ttl",
			expectError: false,
		},
		{
			name:        "collection_subject",
			input:       "spec_tests/collection_subject.ttl",
			expectError: false,
		},
		{
			name:        "collection_object",
			input:       "spec_tests/collection_object.ttl",
			expectError: false,
		},
		{
			name:        "empty_collection",
			input:       "spec_tests/empty_collection.ttl",
			expectError: false,
		},
		{
			name:        "nested_collection",
			input:       "spec_tests/nested_collection.ttl",
			expectError: false,
		},
		{
			name:        "first",
			input:       "spec_tests/first.ttl",
			expectError: false,
		},
		{
			name:        "last",
			input:       "spec_tests/last.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL1",
			input:       "spec_tests/LITERAL1.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL1_ascii_boundaries",
			input:       "spec_tests/LITERAL1_ascii_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL1_with_UTF8_boundaries",
			input:       "spec_tests/LITERAL1_with_UTF8_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL1_all_controls",
			input:       "spec_tests/LITERAL1_all_controls.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL1_all_punctuation",
			input:       "spec_tests/LITERAL1_all_punctuation.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG1",
			input:       "spec_tests/LITERAL_LONG1.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG1_ascii_boundaries",
			input:       "spec_tests/LITERAL_LONG1_ascii_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG1_with_UTF8_boundaries",
			input:       "spec_tests/LITERAL_LONG1_with_UTF8_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG1_with_1_squote",
			input:       "spec_tests/LITERAL_LONG1_with_1_squote.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG1_with_2_squotes",
			input:       "spec_tests/LITERAL_LONG1_with_2_squotes.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL2",
			input:       "spec_tests/LITERAL2.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL2_ascii_boundaries",
			input:       "spec_tests/LITERAL2_ascii_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL2_with_UTF8_boundaries",
			input:       "spec_tests/LITERAL2_with_UTF8_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG2",
			input:       "spec_tests/LITERAL_LONG2.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG2_ascii_boundaries",
			input:       "spec_tests/LITERAL_LONG2_ascii_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG2_with_UTF8_boundaries",
			input:       "spec_tests/LITERAL_LONG2_with_UTF8_boundaries.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG2_with_1_squote",
			input:       "spec_tests/LITERAL_LONG2_with_1_squote.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG2_with_2_squotes",
			input:       "spec_tests/LITERAL_LONG2_with_2_squotes.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_CHARACTER_TABULATION",
			input:       "spec_tests/literal_with_CHARACTER_TABULATION.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_BACKSPACE",
			input:       "spec_tests/literal_with_BACKSPACE.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_LINE_FEED",
			input:       "spec_tests/literal_with_LINE_FEED.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_CARRIAGE_RETURN",
			input:       "spec_tests/literal_with_CARRIAGE_RETURN.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_FORM_FEED",
			input:       "spec_tests/literal_with_FORM_FEED.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_REVERSE_SOLIDUS",
			input:       "spec_tests/literal_with_REVERSE_SOLIDUS.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_escaped_CHARACTER_TABULATION",
			input:       "spec_tests/literal_with_escaped_CHARACTER_TABULATION.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_escaped_BACKSPACE",
			input:       "spec_tests/literal_with_escaped_BACKSPACE.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_escaped_LINE_FEED",
			input:       "spec_tests/literal_with_escaped_LINE_FEED.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_escaped_CARRIAGE_RETURN",
			input:       "spec_tests/literal_with_escaped_CARRIAGE_RETURN.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_escaped_FORM_FEED",
			input:       "spec_tests/literal_with_escaped_FORM_FEED.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_numeric_escape4",
			input:       "spec_tests/literal_with_numeric_escape4.ttl",
			expectError: false,
		},
		{
			name:        "literal_with_numeric_escape8",
			input:       "spec_tests/literal_with_numeric_escape8.ttl",
			expectError: false,
		},
		{
			name:        "IRIREF_datatype",
			input:       "spec_tests/IRIREF_datatype.ttl",
			expectError: false,
		},
		{
			name:        "prefixed_name_datatype",
			input:       "spec_tests/prefixed_name_datatype.ttl",
			expectError: false,
		},
		{
			name:        "bareword_integer",
			input:       "spec_tests/bareword_integer.ttl",
			expectError: false,
		},
		{
			name:        "bareword_decimal",
			input:       "spec_tests/bareword_decimal.ttl",
			expectError: false,
		},
		{
			name:        "bareword_double",
			input:       "spec_tests/bareword_double.ttl",
			expectError: false,
		},
		{
			name:        "double_lower_case_e",
			input:       "spec_tests/double_lower_case_e.ttl",
			expectError: false,
		},
		{
			name:        "negative_numeric",
			input:       "spec_tests/negative_numeric.ttl",
			expectError: false,
		},
		{
			name:        "positive_numeric",
			input:       "spec_tests/positive_numeric.ttl",
			expectError: false,
		},
		{
			name:        "numeric_with_leading_0",
			input:       "spec_tests/numeric_with_leading_0.ttl",
			expectError: false,
		},
		{
			name:        "literal_true",
			input:       "spec_tests/literal_true.ttl",
			expectError: false,
		},
		{
			name:        "literal_false",
			input:       "spec_tests/literal_false.ttl",
			expectError: false,
		},
		{
			name:        "langtagged_non_LONG",
			input:       "spec_tests/langtagged_non_LONG.ttl",
			expectError: false,
		},
		{
			name:        "langtagged_LONG",
			input:       "spec_tests/langtagged_LONG.ttl",
			expectError: false,
		},
		{
			name:        "lantag_with_subtag",
			input:       "spec_tests/lantag_with_subtag.ttl",
			expectError: false,
		},
		{
			name:        "objectList_with_two_objects",
			input:       "spec_tests/objectList_with_two_objects.ttl",
			expectError: false,
		},
		{
			name:        "predicateObjectList_with_two_objectLists",
			input:       "spec_tests/predicateObjectList_with_two_objectLists.ttl",
			expectError: false,
		},
		{
			name:        "predicateObjectList_with_blankNodePropertyList_as_object",
			input:       "spec_tests/predicateObjectList_with_blankNodePropertyList_as_object.ttl",
			expectError: false,
		},
		{
			name:        "repeated_semis_at_end",
			input:       "spec_tests/repeated_semis_at_end.ttl",
			expectError: false,
		},
		{
			name:        "repeated_semis_not_at_end",
			input:       "spec_tests/repeated_semis_not_at_end.ttl",
			expectError: false,
		},
		{
			name:        "comment_following_localName",
			input:       "spec_tests/comment_following_localName.ttl",
			expectError: false,
		},
		{
			name:        "number_sign_following_localName",
			input:       "spec_tests/number_sign_following_localName.ttl",
			expectError: false,
		},
		{
			name:        "comment_following_PNAME_NS",
			input:       "spec_tests/comment_following_PNAME_NS.ttl",
			expectError: false,
		},
		{
			name:        "number_sign_following_PNAME_NS",
			input:       "spec_tests/number_sign_following_PNAME_NS.ttl",
			expectError: false,
		},
		{
			name:        "LITERAL_LONG2_with_REVERSE_SOLIDUS",
			input:       "spec_tests/LITERAL_LONG2_with_REVERSE_SOLIDUS.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bad-num-05",
			input:       "spec_tests/turtle-syntax-bad-LITERAL2_with_langtag_and_datatype.ttl",
			expectError: true,
		},
		{
			name:        "two_LITERAL_LONG2s",
			input:       "spec_tests/two_LITERAL_LONG2s.ttl",
			expectError: false,
		},
		{
			name:        "langtagged_LONG_with_subtag",
			input:       "spec_tests/langtagged_LONG_with_subtag.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-file-01",
			input:       "spec_tests/turtle-syntax-file-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-file-02",
			input:       "spec_tests/turtle-syntax-file-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-file-03",
			input:       "spec_tests/turtle-syntax-file-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-uri-01",
			input:       "spec_tests/turtle-syntax-uri-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-uri-02",
			input:       "spec_tests/turtle-syntax-uri-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-uri-03",
			input:       "spec_tests/turtle-syntax-uri-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-uri-04",
			input:       "spec_tests/turtle-syntax-uri-04.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-base-01",
			input:       "spec_tests/turtle-syntax-base-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-base-02",
			input:       "spec_tests/turtle-syntax-base-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-base-03",
			input:       "spec_tests/turtle-syntax-base-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-base-04",
			input:       "spec_tests/turtle-syntax-base-04.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-prefix-01",
			input:       "spec_tests/turtle-syntax-prefix-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-prefix-02",
			input:       "spec_tests/turtle-syntax-prefix-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-prefix-03",
			input:       "spec_tests/turtle-syntax-prefix-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-prefix-04",
			input:       "spec_tests/turtle-syntax-prefix-04.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-prefix-05",
			input:       "spec_tests/turtle-syntax-prefix-05.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-prefix-06",
			input:       "spec_tests/turtle-syntax-prefix-06.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-prefix-07",
			input:       "spec_tests/turtle-syntax-prefix-07.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-prefix-08",
			input:       "spec_tests/turtle-syntax-prefix-08.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-prefix-09",
			input:       "spec_tests/turtle-syntax-prefix-09.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-01",
			input:       "spec_tests/turtle-syntax-string-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-02",
			input:       "spec_tests/turtle-syntax-string-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-03",
			input:       "spec_tests/turtle-syntax-string-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-04",
			input:       "spec_tests/turtle-syntax-string-04.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-05",
			input:       "spec_tests/turtle-syntax-string-05.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-06",
			input:       "spec_tests/turtle-syntax-string-06.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-07",
			input:       "spec_tests/turtle-syntax-string-07.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-08",
			input:       "spec_tests/turtle-syntax-string-08.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-09",
			input:       "spec_tests/turtle-syntax-string-09.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-10",
			input:       "spec_tests/turtle-syntax-string-10.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-string-11",
			input:       "spec_tests/turtle-syntax-string-11.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-str-esc-01",
			input:       "spec_tests/turtle-syntax-str-esc-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-str-esc-02",
			input:       "spec_tests/turtle-syntax-str-esc-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-str-esc-03",
			input:       "spec_tests/turtle-syntax-str-esc-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-pname-esc-01",
			input:       "spec_tests/turtle-syntax-pname-esc-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-pname-esc-02",
			input:       "spec_tests/turtle-syntax-pname-esc-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-pname-esc-03",
			input:       "spec_tests/turtle-syntax-pname-esc-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-01",
			input:       "spec_tests/turtle-syntax-bnode-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-02",
			input:       "spec_tests/turtle-syntax-bnode-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-03",
			input:       "spec_tests/turtle-syntax-bnode-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-04",
			input:       "spec_tests/turtle-syntax-bnode-04.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-05",
			input:       "spec_tests/turtle-syntax-bnode-05.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-06",
			input:       "spec_tests/turtle-syntax-bnode-06.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-07",
			input:       "spec_tests/turtle-syntax-bnode-07.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-08",
			input:       "spec_tests/turtle-syntax-bnode-08.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-09",
			input:       "spec_tests/turtle-syntax-bnode-09.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bnode-10",
			input:       "spec_tests/turtle-syntax-bnode-10.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-01",
			input:       "spec_tests/turtle-syntax-number-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-02",
			input:       "spec_tests/turtle-syntax-number-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-03",
			input:       "spec_tests/turtle-syntax-number-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-04",
			input:       "spec_tests/turtle-syntax-number-04.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-05",
			input:       "spec_tests/turtle-syntax-number-05.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-06",
			input:       "spec_tests/turtle-syntax-number-06.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-07",
			input:       "spec_tests/turtle-syntax-number-07.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-08",
			input:       "spec_tests/turtle-syntax-number-08.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-09",
			input:       "spec_tests/turtle-syntax-number-09.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-10",
			input:       "spec_tests/turtle-syntax-number-10.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-number-11",
			input:       "spec_tests/turtle-syntax-number-11.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-datatypes-01",
			input:       "spec_tests/turtle-syntax-datatypes-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-datatypes-02",
			input:       "spec_tests/turtle-syntax-datatypes-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-kw-01",
			input:       "spec_tests/turtle-syntax-kw-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-kw-02",
			input:       "spec_tests/turtle-syntax-kw-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-kw-03",
			input:       "spec_tests/turtle-syntax-kw-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-struct-01",
			input:       "spec_tests/turtle-syntax-struct-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-struct-02",
			input:       "spec_tests/turtle-syntax-struct-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-struct-03",
			input:       "spec_tests/turtle-syntax-struct-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-struct-04",
			input:       "spec_tests/turtle-syntax-struct-04.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-struct-05",
			input:       "spec_tests/turtle-syntax-struct-05.ttl",
			expectError: false,
		},
		{
			name:        "turtle-eval-lists-01",
			input:       "spec_tests/turtle-eval-lists-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-eval-lists-02",
			input:       "spec_tests/turtle-eval-lists-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-eval-lists-03",
			input:       "spec_tests/turtle-eval-lists-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-eval-lists-04",
			input:       "spec_tests/turtle-eval-lists-04.ttl",
			expectError: false,
		},
		{
			name:        "turtle-eval-lists-05",
			input:       "spec_tests/turtle-eval-lists-05.ttl",
			expectError: false,
		},
		{
			name:        "turtle-eval-lists-06",
			input:       "spec_tests/turtle-eval-lists-06.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-bad-uri-01",
			input:       "spec_tests/turtle-syntax-bad-uri-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-uri-02",
			input:       "spec_tests/turtle-syntax-bad-uri-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-uri-03",
			input:       "spec_tests/turtle-syntax-bad-uri-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-uri-04",
			input:       "spec_tests/turtle-syntax-bad-uri-04.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-uri-05",
			input:       "spec_tests/turtle-syntax-bad-uri-05.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-prefix-01",
			input:       "spec_tests/turtle-syntax-bad-prefix-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-prefix-02",
			input:       "spec_tests/turtle-syntax-bad-prefix-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-prefix-03",
			input:       "spec_tests/turtle-syntax-bad-prefix-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-prefix-04",
			input:       "spec_tests/turtle-syntax-bad-prefix-04.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-prefix-05",
			input:       "spec_tests/turtle-syntax-bad-prefix-05.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-base-01",
			input:       "spec_tests/turtle-syntax-bad-base-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-base-02",
			input:       "spec_tests/turtle-syntax-bad-base-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-base-03",
			input:       "spec_tests/turtle-syntax-bad-base-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-bnode-01",
			input:       "spec_tests/turtle-syntax-bad-bnode-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-bnode-02",
			input:       "spec_tests/turtle-syntax-bad-bnode-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-01",
			input:       "spec_tests/turtle-syntax-bad-struct-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-02",
			input:       "spec_tests/turtle-syntax-bad-struct-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-03",
			input:       "spec_tests/turtle-syntax-bad-struct-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-04",
			input:       "spec_tests/turtle-syntax-bad-struct-04.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-05",
			input:       "spec_tests/turtle-syntax-bad-struct-05.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-06",
			input:       "spec_tests/turtle-syntax-bad-struct-06.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-07",
			input:       "spec_tests/turtle-syntax-bad-struct-07.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-kw-01",
			input:       "spec_tests/turtle-syntax-bad-kw-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-kw-02",
			input:       "spec_tests/turtle-syntax-bad-kw-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-kw-03",
			input:       "spec_tests/turtle-syntax-bad-kw-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-kw-04",
			input:       "spec_tests/turtle-syntax-bad-kw-04.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-kw-05",
			input:       "spec_tests/turtle-syntax-bad-kw-05.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-01",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-02",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-03",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-04",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-04.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-05",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-05.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-06",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-06.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-07",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-07.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-08",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-08.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-09",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-09.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-10",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-10.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-11",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-11.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-12",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-12.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-n3-extras-13",
			input:       "spec_tests/turtle-syntax-bad-n3-extras-13.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-08",
			input:       "spec_tests/turtle-syntax-bad-struct-08.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-09",
			input:       "spec_tests/turtle-syntax-bad-struct-09.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-10",
			input:       "spec_tests/turtle-syntax-bad-struct-10.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-11",
			input:       "spec_tests/turtle-syntax-bad-struct-11.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-12",
			input:       "spec_tests/turtle-syntax-bad-struct-12.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-13",
			input:       "spec_tests/turtle-syntax-bad-struct-13.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-14",
			input:       "spec_tests/turtle-syntax-bad-struct-14.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-15",
			input:       "spec_tests/turtle-syntax-bad-struct-15.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-16",
			input:       "spec_tests/turtle-syntax-bad-struct-16.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-struct-17",
			input:       "spec_tests/turtle-syntax-bad-struct-17.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-lang-01",
			input:       "spec_tests/turtle-syntax-bad-lang-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-esc-01",
			input:       "spec_tests/turtle-syntax-bad-esc-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-esc-02",
			input:       "spec_tests/turtle-syntax-bad-esc-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-esc-03",
			input:       "spec_tests/turtle-syntax-bad-esc-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-esc-04",
			input:       "spec_tests/turtle-syntax-bad-esc-04.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-pname-01",
			input:       "spec_tests/turtle-syntax-bad-pname-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-pname-02",
			input:       "spec_tests/turtle-syntax-bad-pname-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-pname-03",
			input:       "spec_tests/turtle-syntax-bad-pname-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-string-01",
			input:       "spec_tests/turtle-syntax-bad-string-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-string-02",
			input:       "spec_tests/turtle-syntax-bad-string-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-string-03",
			input:       "spec_tests/turtle-syntax-bad-string-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-string-04",
			input:       "spec_tests/turtle-syntax-bad-string-04.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-string-05",
			input:       "spec_tests/turtle-syntax-bad-string-05.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-string-06",
			input:       "spec_tests/turtle-syntax-bad-string-06.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-string-07",
			input:       "spec_tests/turtle-syntax-bad-string-07.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-num-01",
			input:       "spec_tests/turtle-syntax-bad-num-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-num-02",
			input:       "spec_tests/turtle-syntax-bad-num-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-num-03",
			input:       "spec_tests/turtle-syntax-bad-num-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-num-04",
			input:       "spec_tests/turtle-syntax-bad-num-04.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-num-05",
			input:       "spec_tests/turtle-syntax-bad-num-05.ttl",
			expectError: true,
		},
		{
			name:        "turtle-eval-struct-01",
			input:       "spec_tests/turtle-eval-struct-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-eval-struct-02",
			input:       "spec_tests/turtle-eval-struct-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-01",
			input:       "spec_tests/turtle-subm-01.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-02",
			input:       "spec_tests/turtle-subm-02.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-03",
			input:       "spec_tests/turtle-subm-03.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-04",
			input:       "spec_tests/turtle-subm-04.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-05",
			input:       "spec_tests/turtle-subm-05.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-06",
			input:       "spec_tests/turtle-subm-06.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-07",
			input:       "spec_tests/turtle-subm-07.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-08",
			input:       "spec_tests/turtle-subm-08.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-09",
			input:       "spec_tests/turtle-subm-09.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-10",
			input:       "spec_tests/turtle-subm-10.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-11",
			input:       "spec_tests/turtle-subm-11.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-12",
			input:       "spec_tests/turtle-subm-12.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-13",
			input:       "spec_tests/turtle-subm-13.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-14",
			input:       "spec_tests/turtle-subm-14.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-15",
			input:       "spec_tests/turtle-subm-15.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-16",
			input:       "spec_tests/turtle-subm-16.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-17",
			input:       "spec_tests/turtle-subm-17.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-18",
			input:       "spec_tests/turtle-subm-18.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-19",
			input:       "spec_tests/turtle-subm-19.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-20",
			input:       "spec_tests/turtle-subm-20.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-21",
			input:       "spec_tests/turtle-subm-21.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-22",
			input:       "spec_tests/turtle-subm-22.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-23",
			input:       "spec_tests/turtle-subm-23.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-24",
			input:       "spec_tests/turtle-subm-24.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-25",
			input:       "spec_tests/turtle-subm-25.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-26",
			input:       "spec_tests/turtle-subm-26.ttl",
			expectError: false,
		},
		{
			name:        "turtle-subm-27",
			input:       "spec_tests/turtle-subm-27.ttl",
			expectError: false,
		},
		{
			name:        "turtle-eval-bad-01",
			input:       "spec_tests/turtle-eval-bad-01.ttl",
			expectError: true,
		},
		{
			name:        "turtle-eval-bad-02",
			input:       "spec_tests/turtle-eval-bad-02.ttl",
			expectError: true,
		},
		{
			name:        "turtle-eval-bad-03",
			input:       "spec_tests/turtle-eval-bad-03.ttl",
			expectError: true,
		},
		{
			name:        "turtle-eval-bad-04",
			input:       "spec_tests/turtle-eval-bad-04.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-blank-label-dot-end",
			input:       "spec_tests/turtle-syntax-bad-blank-label-dot-end.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-ln-dash-start",
			input:       "spec_tests/turtle-syntax-bad-ln-dash-start.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-ln-escape-start",
			input:       "spec_tests/turtle-syntax-bad-ln-escape-start.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-ln-escape",
			input:       "spec_tests/turtle-syntax-bad-ln-escape.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-missing-ns-dot-end",
			input:       "spec_tests/turtle-syntax-bad-missing-ns-dot-end.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-missing-ns-dot-start",
			input:       "spec_tests/turtle-syntax-bad-missing-ns-dot-start.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-ns-dot-end",
			input:       "spec_tests/turtle-syntax-bad-ns-dot-end.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-ns-dot-start",
			input:       "spec_tests/turtle-syntax-bad-ns-dot-start.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-bad-number-dot-in-anon",
			input:       "spec_tests/turtle-syntax-bad-number-dot-in-anon.ttl",
			expectError: true,
		},
		{
			name:        "turtle-syntax-blank-label",
			input:       "spec_tests/turtle-syntax-blank-label.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-ln-colons",
			input:       "spec_tests/turtle-syntax-ln-colons.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-ln-dots",
			input:       "spec_tests/turtle-syntax-ln-dots.ttl",
			expectError: false,
		},
		{
			name:        "turtle-syntax-ns-dots",
			input:       "spec_tests/turtle-syntax-ns-dots.ttl",
			expectError: false,
		},
		{
			name:        "IRI-resolution-01",
			input:       "spec_tests/IRI-resolution-01.ttl",
			expectError: false,
		},
		{
			name:        "IRI-resolution-02",
			input:       "spec_tests/IRI-resolution-02.ttl",
			expectError: false,
		},
		{
			name:        "IRI-resolution-07",
			input:       "spec_tests/IRI-resolution-07.ttl",
			expectError: false,
		},
		{
			name:        "IRI-resolution-08",
			input:       "spec_tests/IRI-resolution-08.ttl",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.Open(tc.input)
			if err != nil {
				t.Fatalf("Failed to open file %s: %v", tc.input, err)
			}
			defer data.Close()

			outChan, errChan := ParseTurtle(data, &TurtleParserOptions{Base: "http://example.org/"})

			for range outChan {
				// We don't validate quads here, just parsing behavior
			}

			err, ok := <-errChan
			if tc.expectError {
				if err == nil && ok {
					t.Errorf("Expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	}
}

func TestStreaming(t *testing.T) {
	t.Run("Chunked input - across lines", func(t *testing.T) {
		chunked := []string{
			"<http://example.org/s> <http://example.org/p> ",
			"<http://example.org/o> .\n",
			"<http://example.org/s> <http://example.org/p> \"lit",
			"eral\" .\n",
			"_:b1 <http://example.org/p> _:b2 .",
		}

		r, w := io.Pipe()
		go func() {
			for _, chunk := range chunked {
				_, err := w.Write([]byte(chunk))
				if err != nil {
					t.Errorf("Write error: %v", err)
					return
				}
			}
			w.Close()
		}()

		outChan, errChan := ParseTurtle(r, nil)

		var got []interfaces.IQuad
		for quad := range outChan {
			got = append(got, quad)
		}

		if err, ok := <-errChan; err != nil && ok {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(got) != 3 {
			t.Errorf("Expected 3 quads, got %d", len(got))
		}
	})

	t.Run("Simulate scanner read error", func(t *testing.T) {
		badReader := iotest.ErrReader(errors.New("simulated read error"))
		_, errChan := ParseTurtle(badReader, nil)

		err, ok := <-errChan
		if !ok || err == nil || !strings.Contains(err.Error(), "simulated read error") {
			t.Errorf("Expected read error, got: %v", err)
		}
	})

	t.Run("Input ends mid-token", func(t *testing.T) {
		in := `<http://example.org/s> <http://example.org/p> "unterminated`
		outChan, errChan := ParseTurtle(strings.NewReader(in), nil)

		var got []interfaces.IQuad
		for q := range outChan {
			got = append(got, q)
		}
		if err, ok := <-errChan; err == nil || !ok {
			t.Errorf("Expected parse error for unterminated literal, got: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("Expected 0 quads, got %d", len(got))
		}
	})

	t.Run("Empty input", func(t *testing.T) {
		outChan, errChan := ParseTurtle(strings.NewReader(""), nil)

		var got []interfaces.IQuad
		for q := range outChan {
			got = append(got, q)
		}
		if len(got) != 0 {
			t.Errorf("Expected 0 quads for empty input, got %d", len(got))
		}
		if err := <-errChan; err != nil {
			t.Errorf("Unexpected error for empty input: %v", err)
		}
	})

	t.Run("Delayed input", func(t *testing.T) {
		pr, pw := io.Pipe()

		go func() {
			_, err := pw.Write([]byte("<http://example.org/s> "))
			if err != nil {
				t.Errorf("Write error: %v", err)
				return
			}
			time.Sleep(50 * time.Millisecond)
			_, err = pw.Write([]byte("<http://example.org/p> "))
			if err != nil {
				t.Errorf("Write error: %v", err)
				return
			}
			time.Sleep(50 * time.Millisecond)
			_, err = pw.Write([]byte("<http://example.org/o> .\n"))
			if err != nil {
				t.Errorf("Write error: %v", err)
				return
			}
			pw.Close()
		}()

		outChan, errChan := ParseTurtle(pr, nil)
		var got []interfaces.IQuad
		for q := range outChan {
			got = append(got, q)
		}
		if err := <-errChan; err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("Expected 1 quad, got %d", len(got))
		}
	})
}

func TestErrors(t *testing.T) {
	t.Run("Unterminated literal", func(t *testing.T) {
		input := `<http://example.org/s> <http://example.org/p> "unterminated`
		out, errChan := ParseTurtle(strings.NewReader(input), nil)

		for range out {
			// drain
		}
		err, ok := <-errChan
		if !ok || err == nil {
			t.Errorf("Expected error for unterminated literal, got: %v", err)
		}
	})

	t.Run("Illegal token", func(t *testing.T) {
		input := `<http://example.org/s> <http://example.org/p> ?? .`
		out, errChan := ParseTurtle(strings.NewReader(input), nil)

		for range out {
			// drain
		}
		err, ok := <-errChan
		if !ok || err == nil || !strings.Contains(err.Error(), "Syntax error") {
			t.Errorf("Expected syntax error for illegal token, got: %v", err)
		}
	})

	t.Run("Bad prefix usage", func(t *testing.T) {
		input := `foo:bar <http://example.org/p> <http://example.org/o> .`
		out, errChan := ParseTurtle(strings.NewReader(input), nil)

		for range out {
		}
		err, ok := <-errChan
		if !ok || err == nil || !strings.Contains(err.Error(), "prefix not found") {
			t.Errorf("Expected error for undefined prefix, got: %v", err)
		}
	})

	t.Run("Base-relative IRI without @base", func(t *testing.T) {
		input := `<thing> <http://example.org/p> <http://example.org/o> .`
		out, errChan := ParseTurtle(strings.NewReader(input), nil)

		for range out {
		}
		err, ok := <-errChan
		if !ok || err == nil || !strings.Contains(err.Error(), "@base not defined") {
			t.Errorf("Expected base error for relative IRI, got: %v", err)
		}
	})

	t.Run("Scanner read error", func(t *testing.T) {
		badReader := iotest.ErrReader(errors.New("fake read error"))
		_, errChan := ParseTurtle(badReader, nil)

		err, ok := <-errChan
		if !ok || err == nil || !strings.Contains(err.Error(), "fake read error") {
			t.Errorf("Expected scanner error, got: %v", err)
		}
	})

	t.Run("Parser failure triggers fallback message", func(t *testing.T) {
		// This input should parse but will trigger a deliberate syntax error if your parser can't continue after object
		input := `<http://example.org/s> <http://example.org/p> .`
		out, errChan := ParseTurtle(strings.NewReader(input), nil)

		for range out {
		}
		err, ok := <-errChan
		if !ok || err == nil || !strings.Contains(err.Error(), "Syntax error near token: \".\" (syntax error: unexpected DOT)") {
			t.Errorf("Expected parser failure fallback error, got: %v", err)
		}
	})
}
