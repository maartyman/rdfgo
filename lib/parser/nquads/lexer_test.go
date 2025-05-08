package nquads

import (
	"errors"
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

type QuadTestCase struct {
	name        string
	input       string
	expectError bool
	expectQuads []interfaces.IQuad
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
	namedNode := NewNamedNode(val)
	return namedNode
}

func bn(val string) interfaces.IBlankNode {
	blankNode := NewBlankNode(val)
	return blankNode
}

func l(val string) interfaces.ILiteral {
	literal := NewStringLiteral(val, "")
	return literal
}

func lit(val string, lang string, dt interfaces.INamedNode) interfaces.ILiteral {
	literal := NewLiteral(val, lang, dt)
	return literal
}

func TestNQuadsOutput(t *testing.T) {
	tests := []QuadTestCase{
		{
			name:  "Simple quad with graph",
			input: `<s> <p> <o> <g> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("o"), nn("g")),
			},
		},
		{
			name:  "literal with lang",
			input: `<s> <p> "hello" @en .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), lit("hello", "en", nil), nil),
			},
		},
		{
			name:  "literal with datatype",
			input: `<s> <p> "42"^^<http://www.w3.org/2001/XMLSchema#integer> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), lit("42", "", nn("http://www.w3.org/2001/XMLSchema#integer")), NewDefaultGraph()),
			},
		},
		{
			name:  "Blank node subject",
			input: `_:b1 <p> <o> .`,
			expectQuads: []interfaces.IQuad{
				q(bn("b0"), nn("p"), nn("o"), NewDefaultGraph()),
			},
		},
		{
			name:  "Blank node object and graph",
			input: `<s> <p> _:obj _:graph .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), bn("b0"), bn("b1")),
			},
		},
		{
			name:  "Escaped quote in literal",
			input: `<s> <p> "he said \"hello\"" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l(`he said "hello"`), NewDefaultGraph()),
			},
		},
		{
			name:        "Invalid quad (too few terms)",
			input:       `<s> <p> .`,
			expectError: true,
		},
		{
			name:        "Invalid IRI syntax",
			input:       `<s> <p> <invalid .`,
			expectError: true,
		},
		{
			name:  "IRI with fragment",
			input: `<http://example.org/s#frag> <p> <o> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("http://example.org/s#frag"), nn("p"), nn("o"), NewDefaultGraph()),
			},
		},
		{
			name:  "Triple with comment at end",
			input: `<s> <p> <o> . # this is a comment`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("o"), NewDefaultGraph()),
			},
		},
		{
			name:  "Triple with comment in graph",
			input: `<s> <p> <o> <g#1> .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("o"), nn("g#1")),
			},
		},
		{
			name:  "Whitespace variations",
			input: `   <s>   <p>   <o>   .   `,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("o"), NewDefaultGraph()),
			},
		},
		{
			name: "Empty line is skipped",
			input: `

<s> <p> <o> .

`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("o"), NewDefaultGraph()),
			},
		},
		{
			name:  "literal is just dot character",
			input: `<s> <p> "." .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("."), NewDefaultGraph()),
			},
		},
		{
			name:  "Unicode escape literal",
			input: `<s> <p> "hello\u0020world" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("hello world"), NewDefaultGraph()),
			},
		},
		{
			name:  "Multiple quads",
			input: "<s> <p> <o> .\n<s> <p> \"literal\" .",
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("o"), NewDefaultGraph()),
				q(nn("s"), nn("p"), l("literal"), NewDefaultGraph()),
			},
		},
		{
			name:  "Multiple quads split with a tab",
			input: "<s> <p> <o> .\t<s> <p> \"literal\" .",
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), nn("o"), NewDefaultGraph()),
				q(nn("s"), nn("p"), l("literal"), NewDefaultGraph()),
			},
		},
		{
			name:        "literal with both lang and datatype (invalid)",
			input:       `<s> <p> "hello"@en^^<http://example.org/type> .`,
			expectError: true,
		},
		{
			name:        "Unterminated literal",
			input:       `<s> <p> "unterminated .`,
			expectError: true,
		},
		{
			name:  "Escaped quote",
			input: `<s> <p> "He said \"hi\"" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l(`He said "hi"`), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped backslash",
			input: `<s> <p> "Path: C:\\Windows\\" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l(`Path: C:\Windows\`), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped newline",
			input: `<s> <p> "line1\nline2" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("line1\nline2"), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped tab",
			input: `<s> <p> "column1\tcolumn2" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("column1\tcolumn2"), NewDefaultGraph()),
			},
		},
		{
			name:  "Unicode escape: basic",
			input: `<s> <p> "heart: \u2665" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("heart: ♥"), NewDefaultGraph()),
			},
		},
		{
			name:  "Unicode escape: emoji",
			input: `<s> <p> "emoji: \U0001F600" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("emoji: 😀"), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped backslash and quote",
			input: `<s> <p> "She said: \\\"cool\\\"" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l(`She said: \"cool\"`), NewDefaultGraph()),
			},
		},
		{
			name:        "Illegal escape sequence",
			input:       `<s> <p> "bad\qescape" .`,
			expectError: true,
		},
		{
			name:        "Broken unicode escape",
			input:       `<s> <p> "oops\u00G1" .`,
			expectError: true,
		},
		{
			name:  "Escaped newline \\n",
			input: `<s> <p> "line1\nline2" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("line1\nline2"), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped tab \\t",
			input: `<s> <p> "col1\tcol2" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("col1\tcol2"), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped carriage return \\r",
			input: `<s> <p> "line1\rline2" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("line1\rline2"), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped backspace \\b",
			input: `<s> <p> "a\bback" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("a\bback"), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped form feed \\f",
			input: `<s> <p> "a\fform" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l("a\fform"), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped quote \\\"",
			input: `<s> <p> "She said: \"hi\"" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l(`She said: "hi"`), NewDefaultGraph()),
			},
		},
		{
			name:  "Escaped backslash \\\\",
			input: `<s> <p> "C:\\\\path" .`,
			expectQuads: []interfaces.IQuad{
				q(nn("s"), nn("p"), l(`C:\\path`), NewDefaultGraph()),
			},
		},
		{
			name:        "Invalid token: unexpected character",
			input:       `<s> <p> $invalid .`,
			expectError: true,
		},
		{
			name:        "Invalid token: random junk",
			input:       `<s> <p> ?!? .`,
			expectError: true,
		},
		{
			name:        "Invalid token: single quote literal (N-Quads does not support it)",
			input:       `<s> <p> 'single-quoted' .`,
			expectError: true,
		},
		{
			name:        "Invalid token: numeric literal (not supported)",
			input:       `<s> <p> 123 .`,
			expectError: true,
		},
		{
			name:        "Invalid token: bare identifier",
			input:       `<s> <p> bareword .`,
			expectError: true,
		},
		{
			name:        "Invalid token: lone comment token",
			input:       `#justacomment`,
			expectQuads: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outChan, errChan := Parse(strings.NewReader(tc.input), Options{})

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

func TestNQuadsSpec(t *testing.T) {
	tests := []QuadTestCase{
		{
			name:        "nt-syntax-bad-bnode-01",
			input:       "spec_tests/nt-syntax-bad-bnode-01.nq",
			expectError: true,
		},
		{
			name:        "lantag_with_subtag",
			input:       "spec_tests/lantag_with_subtag.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-bnode-02",
			input:       "spec_tests/nq-syntax-bnode-02.nq",
			expectError: false,
		},
		{
			name:        "literal_with_CHARACTER_TABULATION",
			input:       "spec_tests/literal_with_CHARACTER_TABULATION.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-bnode-05",
			input:       "spec_tests/nq-syntax-bnode-05.nq",
			expectError: false,
		},
		{
			name:        "literal_with_2_squotes",
			input:       "spec_tests/literal_with_2_squotes.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-uri-02",
			input:       "spec_tests/nt-syntax-uri-02.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-subm-01",
			input:       "spec_tests/nt-syntax-subm-01.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-file-03",
			input:       "spec_tests/nt-syntax-file-03.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-base-01",
			input:       "spec_tests/nt-syntax-bad-base-01.nq",
			expectError: true,
		},
		{
			name:        "nq-syntax-bnode-03",
			input:       "spec_tests/nq-syntax-bnode-03.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-bnode-01",
			input:       "spec_tests/nq-syntax-bnode-01.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-uri-01",
			input:       "spec_tests/nq-syntax-uri-01.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-bad-literal-03",
			input:       "spec_tests/nq-syntax-bad-literal-03.nq",
			expectError: true,
		},
		{
			name:        "nq-syntax-bad-quint-01",
			input:       "spec_tests/nq-syntax-bad-quint-01.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-string-05",
			input:       "spec_tests/nt-syntax-bad-string-05.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-string-06",
			input:       "spec_tests/nt-syntax-bad-string-06.nq",
			expectError: true,
		},
		{
			name:        "literal_with_numeric_escape4",
			input:       "spec_tests/literal_with_numeric_escape4.nq",
			expectError: false,
		},
		{
			name:        "langtagged_string",
			input:       "spec_tests/langtagged_string.nq",
			expectError: false,
		},
		{
			name:        "literal_with_squote",
			input:       "spec_tests/literal_with_squote.nq",
			expectError: false,
		},
		{
			name:        "minimal_whitespace",
			input:       "spec_tests/minimal_whitespace.nq",
			expectError: false,
		},
		{
			name:        "literal_with_2_dquotes",
			input:       "spec_tests/literal_with_2_dquotes.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-uri-02",
			input:       "spec_tests/nq-syntax-uri-02.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-uri-04",
			input:       "spec_tests/nt-syntax-bad-uri-04.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-string-04",
			input:       "spec_tests/nt-syntax-bad-string-04.nq",
			expectError: true,
		},
		{
			name:        "literal_with_UTF8_boundaries",
			input:       "spec_tests/literal_with_UTF8_boundaries.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-uri-03",
			input:       "spec_tests/nq-syntax-uri-03.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-uri-03",
			input:       "spec_tests/nt-syntax-uri-03.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-datatypes-01",
			input:       "spec_tests/nt-syntax-datatypes-01.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-num-03",
			input:       "spec_tests/nt-syntax-bad-num-03.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-uri-09",
			input:       "spec_tests/nt-syntax-bad-uri-09.nq",
			expectError: true,
		},
		{
			name:        "nq-syntax-bad-literal-02",
			input:       "spec_tests/nq-syntax-bad-literal-02.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-lang-01",
			input:       "spec_tests/nt-syntax-bad-lang-01.nq",
			expectError: true,
		},
		{
			name:        "literal",
			input:       "spec_tests/literal.nq",
			expectError: false,
		},
		{
			name:        "literal_ascii_boundaries",
			input:       "spec_tests/literal_ascii_boundaries.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-esc-02",
			input:       "spec_tests/nt-syntax-bad-esc-02.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-uri-02",
			input:       "spec_tests/nt-syntax-bad-uri-02.nq",
			expectError: true,
		},
		{
			name:        "literal_with_BACKSPACE",
			input:       "spec_tests/literal_with_BACKSPACE.nq",
			expectError: false,
		},
		{
			name:        "literal_with_REVERSE_SOLIDUS2",
			input:       "spec_tests/literal_with_REVERSE_SOLIDUS2.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-bnode-06",
			input:       "spec_tests/nq-syntax-bnode-06.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-str-esc-03",
			input:       "spec_tests/nt-syntax-str-esc-03.nq",
			expectError: false,
		},
		{
			name:        "literal_with_CARRIAGE_RETURN",
			input:       "spec_tests/literal_with_CARRIAGE_RETURN.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-bad-uri-01",
			input:       "spec_tests/nq-syntax-bad-uri-01.nq",
			expectError: true,
		},
		{
			name:        "literal_with_dquote",
			input:       "spec_tests/literal_with_dquote.nq",
			expectError: false,
		},
		{
			name:        "literal_all_controls",
			input:       "spec_tests/literal_all_controls.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-uri-06",
			input:       "spec_tests/nq-syntax-uri-06.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-string-02",
			input:       "spec_tests/nt-syntax-string-02.nq",
			expectError: false,
		},
		{
			name:        "literal_with_FORM_FEED",
			input:       "spec_tests/literal_with_FORM_FEED.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-bnode-04",
			input:       "spec_tests/nq-syntax-bnode-04.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-esc-03",
			input:       "spec_tests/nt-syntax-bad-esc-03.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-string-07",
			input:       "spec_tests/nt-syntax-bad-string-07.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-str-esc-02",
			input:       "spec_tests/nt-syntax-str-esc-02.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-esc-01",
			input:       "spec_tests/nt-syntax-bad-esc-01.nq",
			expectError: true,
		},
		{
			name:        "nq-syntax-uri-04",
			input:       "spec_tests/nq-syntax-uri-04.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-num-01",
			input:       "spec_tests/nt-syntax-bad-num-01.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-num-02",
			input:       "spec_tests/nt-syntax-bad-num-02.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-string-01",
			input:       "spec_tests/nt-syntax-bad-string-01.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-struct-02",
			input:       "spec_tests/nt-syntax-bad-struct-02.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-uri-01",
			input:       "spec_tests/nt-syntax-bad-uri-01.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-str-esc-01",
			input:       "spec_tests/nt-syntax-str-esc-01.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bnode-02",
			input:       "spec_tests/nt-syntax-bnode-02.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-string-03",
			input:       "spec_tests/nt-syntax-string-03.nq",
			expectError: false,
		},
		{
			name:        "literal_with_numeric_escape8",
			input:       "spec_tests/literal_with_numeric_escape8.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-uri-04",
			input:       "spec_tests/nt-syntax-uri-04.nq",
			expectError: false,
		},
		{
			name:        "literal_all_punctuation",
			input:       "spec_tests/literal_all_punctuation.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-file-02",
			input:       "spec_tests/nt-syntax-file-02.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-uri-07",
			input:       "spec_tests/nt-syntax-bad-uri-07.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-struct-01",
			input:       "spec_tests/nt-syntax-bad-struct-01.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-uri-05",
			input:       "spec_tests/nt-syntax-bad-uri-05.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-uri-06",
			input:       "spec_tests/nt-syntax-bad-uri-06.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-string-02",
			input:       "spec_tests/nt-syntax-bad-string-02.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-uri-08",
			input:       "spec_tests/nt-syntax-bad-uri-08.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-bnode-02",
			input:       "spec_tests/nt-syntax-bad-bnode-02.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-uri-01",
			input:       "spec_tests/nt-syntax-uri-01.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-uri-03",
			input:       "spec_tests/nt-syntax-bad-uri-03.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-file-01",
			input:       "spec_tests/nt-syntax-file-01.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bnode-01",
			input:       "spec_tests/nt-syntax-bnode-01.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bad-prefix-01",
			input:       "spec_tests/nt-syntax-bad-prefix-01.nq",
			expectError: true,
		},
		{
			name:        "nt-syntax-bad-string-03",
			input:       "spec_tests/nt-syntax-bad-string-03.nq",
			expectError: true,
		},
		{
			name:        "comment_following_triple",
			input:       "spec_tests/comment_following_triple.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-datatypes-02",
			input:       "spec_tests/nt-syntax-datatypes-02.nq",
			expectError: false,
		},
		{
			name:        "nq-syntax-bad-literal-01",
			input:       "spec_tests/nq-syntax-bad-literal-01.nq",
			expectError: true,
		},
		{
			name:        "nq-syntax-uri-05",
			input:       "spec_tests/nq-syntax-uri-05.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-bnode-03",
			input:       "spec_tests/nt-syntax-bnode-03.nq",
			expectError: false,
		},
		{
			name:        "literal_with_REVERSE_SOLIDUS",
			input:       "spec_tests/literal_with_REVERSE_SOLIDUS.nq",
			expectError: false,
		},
		{
			name:        "literal_with_LINE_FEED",
			input:       "spec_tests/literal_with_LINE_FEED.nq",
			expectError: false,
		},
		{
			name:        "nt-syntax-string-01",
			input:       "spec_tests/nt-syntax-string-01.nq",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.Open(tc.input)
			if err != nil {
				panic(err)
			}
			defer data.Close()

			outChan, errChan := Parse(data, Options{})

			for range outChan {
			}

			err, ok := <-errChan
			if tc.expectError {
				if err == nil && ok {
					t.Errorf("Expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}
		})
	}
}

func TestNQuadsStreaming(t *testing.T) {
	t.Run("Chunked input - across lines", func(t *testing.T) {
		chunked := []string{
			"<", "s> <", "p> ", "<o> .\n",
			"<s> <p>", " \"l", "iteral\" .\n",
			"_:b1 <p> _:b2", " .",
		}

		r, w := io.Pipe()
		go func() {
			for _, chunk := range chunked {
				_, err := w.Write([]byte(chunk))
				if err != nil {
					t.Errorf("Error in test: %v", err)
					return
				}
			}
			err := w.Close()
			if err != nil {
				t.Errorf("Error in test: %v", err)
				return
			}
		}()

		outChan, errChan := Parse(r, Options{})

		var got []interfaces.IQuad
		for quad := range outChan {
			got = append(got, quad)
		}
		err, ok := <-errChan
		if err != nil && ok {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(got) != 3 {
			t.Errorf("Expected 3 quads, got %d", len(got))
		}
	})

	t.Run("Simulate scanner read error", func(t *testing.T) {
		badReader := iotest.ErrReader(errors.New("simulated read error"))
		_, errChan := Parse(badReader, Options{})

		err, ok := <-errChan
		if !ok || err == nil || !strings.Contains(err.Error(), "simulated read error") {
			t.Errorf("Expected read error, got: %v", err)
		}
	})

	t.Run("Input ends mid-token", func(t *testing.T) {
		in := `<s> <p> "unterminated`
		outChan, errChan := Parse(strings.NewReader(in), Options{})

		var got []interfaces.IQuad
		for q := range outChan {
			got = append(got, q)
		}
		err, ok := <-errChan
		if err == nil || !ok {
			t.Errorf("Expected parse error for unterminated literal, got: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("Expected 0 quads, got %d", len(got))
		}
	})

	t.Run("Empty input", func(t *testing.T) {
		outChan, errChan := Parse(strings.NewReader(""), Options{})

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
			_, err := pw.Write([]byte("<s> "))
			if err != nil {
				t.Errorf("Error in test: %v", err)
				return
			}
			_, err = pw.Write([]byte("<p> "))
			if err != nil {
				t.Errorf("Error in test: %v", err)
				return
			}
			_, err = pw.Write([]byte("<o> .\n"))
			if err != nil {
				t.Errorf("Error in test: %v", err)
				return
			}
			time.Sleep(100 * time.Millisecond)
			err = pw.Close()
			if err != nil {
				t.Errorf("Error in test: %v", err)
				return
			}
		}()

		outChan, errChan := Parse(pr, Options{})
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
