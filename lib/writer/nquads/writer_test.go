package nquads

import (
	"bytes"
	"errors"
	"github.com/maartyman/rdfgo/interfaces"
	datamodel "github.com/maartyman/rdfgo/lib/data_model"
	stream "github.com/maartyman/rdfgo/lib/stream"
	"io"
	"strings"
	"testing"
)

// mockWriter to simulate io.Writer errors
type mockWriter struct {
	failAfter int
	writes    int
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.writes++
	if m.writes > m.failAfter {
		return 0, errors.New("simulated write error")
	}
	return len(p), nil
}

func mkQuad(s, p string, o interfaces.ITerm, g interfaces.ITerm) interfaces.IQuad {
	quad, err := datamodel.NewQuad(
		datamodel.NewNamedNode(s),
		datamodel.NewNamedNode(p),
		o,
		g,
	)
	if err != nil {
		panic(err)
	}
	return quad
}

func qMust(q interfaces.IQuad, err error) interfaces.IQuad {
	if err != nil {
		panic(err)
	}
	return q
}

func TestWriteNQuads(t *testing.T) {
	tests := []struct {
		name      string
		quads     []interfaces.IQuad
		wantLines []string
		expectErr bool
		writer    io.Writer
	}{
		{
			name: "basic named nodes",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewNamedNode("http://ex/o"), nil),
			},
			wantLines: []string{
				"<http://ex/s> <http://ex/p> <http://ex/o> .",
			},
		},
		{
			name: "literal with lang tag",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewStringLiteral("hello", "en"), nil),
			},
			wantLines: []string{
				"<http://ex/s> <http://ex/p> \"hello\"@en .",
			},
		},
		{
			name: "quad with named graph",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewIntegerLiteral(42), datamodel.NewNamedNode("http://ex/g")),
			},
			wantLines: []string{
				"<http://ex/s> <http://ex/p> \"42\"^^<http://www.w3.org/2001/XMLSchema#integer> <http://ex/g> .",
			},
		},
		{
			name: "blank node object",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewBlankNode("b1"), nil),
			},
			wantLines: []string{
				"<http://ex/s> <http://ex/p> _:b1 .",
			},
		},
		{
			name: "default graph omission",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewStringLiteral("in default", ""), datamodel.NewDefaultGraph()),
			},
			wantLines: []string{
				"<http://ex/s> <http://ex/p> \"in default\" .",
			},
		},
		{
			name: "named and default graph mixed",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s1", "http://ex/p", datamodel.NewStringLiteral("default", ""), datamodel.NewDefaultGraph()),
				mkQuad("http://ex/s2", "http://ex/p", datamodel.NewStringLiteral("named", ""), datamodel.NewNamedNode("http://ex/g")),
			},
			wantLines: []string{
				"<http://ex/s1> <http://ex/p> \"default\" .",
				"<http://ex/s2> <http://ex/p> \"named\" <http://ex/g> .",
			},
		},
		{
			name: "graph omission vs named node",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewNamedNode("http://ex/o"), nil),
			},
			wantLines: []string{
				"<http://ex/s> <http://ex/p> <http://ex/o> .",
			},
		},
		{
			name: "literal with datatype and lang (lang wins)",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewLiteral("chat", "fr", datamodel.IRI.XSD.String), nil),
			},
			wantLines: []string{
				`<http://ex/s> <http://ex/p> "chat"@fr .`,
			},
		},
		{
			name: "empty literal string",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewStringLiteral("", ""), nil),
			},
			wantLines: []string{
				`<http://ex/s> <http://ex/p> "" .`,
			},
		},
		{
			name: "blank node as subject",
			quads: []interfaces.IQuad{
				qMust(datamodel.NewQuad(
					datamodel.NewBlankNode("abc"),
					datamodel.NewNamedNode("http://ex/p"),
					datamodel.NewNamedNode("http://ex/o"),
					nil,
				)),
			},
			wantLines: []string{
				`_:abc <http://ex/p> <http://ex/o> .`,
			},
		},
		{
			name: "variable term in quad (non-standard)",
			quads: []interfaces.IQuad{
				qMust(datamodel.NewQuad(
					datamodel.NewVariable("var"),
					datamodel.NewNamedNode("http://ex/p"),
					datamodel.NewNamedNode("http://ex/o"),
					nil,
				)),
			},
			wantLines: []string{
				`?var <http://ex/p> <http://ex/o> .`,
			},
		},
		{
			name: "literal with escaped quotes",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewStringLiteral(`He said "hi"`, ""), nil),
			},
			wantLines: []string{
				`<http://ex/s> <http://ex/p> "He said \"hi\"" .`,
			},
		},
		{
			name: "invalid characters in literals",
			quads: []interfaces.IQuad{
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewStringLiteral("Line\nBreak", ""), nil),
				mkQuad("http://ex/s", "http://ex/p", datamodel.NewStringLiteral("Tab\tCharacter", ""), nil),
			},
			wantLines: []string{
				`<http://ex/s> <http://ex/p> "Line\nBreak" .`,
				`<http://ex/s> <http://ex/p> "Tab\tCharacter" .`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataStream := stream.ArrayToStream(tt.quads)

			result, err := WriteNQuads(dataStream.ToIStream(), nil, Options{})
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("unexpected error: %v", err)
				}
			} else if tt.expectErr {
				t.Fatal("expected error but got none")
			}

			if !tt.expectErr {
				for _, want := range tt.wantLines {
					if !strings.Contains(result, want) {
						t.Errorf("missing expected line: %q\noutput:\n%s", want, result)
					}
				}
			}
		})
	}
}

func TestWriteNQuads_Errors(t *testing.T) {
	t.Run("writer returns error immediately", func(t *testing.T) {
		dataStream := stream.ArrayToStream([]interfaces.IQuad{
			mkQuad("http://ex/s", "http://ex/p", datamodel.NewNamedNode("http://ex/o"), nil),
		})

		writer := &mockWriter{failAfter: 0}
		_, err := WriteNQuads(dataStream.ToIStream(), writer, Options{})

		if err == nil {
			t.Error("expected an error from mockWriter but got none")
		}
	})

	t.Run("write error mid-stream", func(t *testing.T) {
		dataStream := stream.ArrayToStream([]interfaces.IQuad{
			mkQuad("http://ex/s1", "http://ex/p", datamodel.NewNamedNode("http://ex/o1"), nil),
			mkQuad("http://ex/s2", "http://ex/p", datamodel.NewNamedNode("http://ex/o2"), nil),
			mkQuad("http://ex/s3", "http://ex/p", datamodel.NewNamedNode("http://ex/o3"), nil),
		})

		writer := &mockWriter{failAfter: 2} // error on third write
		_, err := WriteNQuads(dataStream.ToIStream(), writer, Options{})

		if err == nil {
			t.Error("expected write error mid-stream but got none")
		}
	})

	t.Run("closed stream behaves cleanly", func(t *testing.T) {
		dataStream := stream.NewStream()
		close(dataStream)

		var buf bytes.Buffer
		_, err := WriteNQuads(dataStream.ToIStream(), &buf, Options{})

		if err != nil {
			t.Errorf("unexpected error from closed stream: %v", err)
		}

		if buf.Len() != 0 {
			t.Errorf("expected empty output, got: %s", buf.String())
		}
	})
}

func TestWriteMultipleStores(t *testing.T) {
	// Create two stores with different quads
	store1 := stream.NewStore()
	store2 := stream.NewStore()

	store1.AddQuadFromTerms(
		datamodel.NewNamedNode("http://ex/s1"),
		datamodel.NewNamedNode("http://ex/p1"),
		datamodel.NewNamedNode("http://ex/o1"),
		nil,
	)

	store2.AddQuadFromTerms(
		datamodel.NewNamedNode("http://ex/s2"),
		datamodel.NewNamedNode("http://ex/p2"),
		datamodel.NewNamedNode("http://ex/o2"),
		nil,
	)

	// Combine the stores into one writer
	var buf bytes.Buffer
	stores := []interfaces.IStore{store1, store2}

	for _, store := range stores {
		dataStream := store.Match(nil, nil, nil, nil)
		_, err := WriteNQuads(dataStream, &buf, Options{})
		if err != nil {
			t.Fatalf("unexpected error writing store: %v", err)
		}
	}

	// Verify the output contains quads from both stores
	output := buf.String()
	expectedLines := []string{
		"<http://ex/s1> <http://ex/p1> <http://ex/o1> .",
		"<http://ex/s2> <http://ex/p2> <http://ex/o2> .",
	}

	for _, line := range expectedLines {
		if !strings.Contains(output, line) {
			t.Errorf("missing expected line: %q\noutput:\n%s", line, output)
		}
	}
}
