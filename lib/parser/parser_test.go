package rdfgo

import (
	"os"
	"strings"
	"testing"

	"github.com/maartyman/rdfgo/interfaces"
)

func TestGetExtension(t *testing.T) {
	tests := map[string]string{
		"file.nq":      ".nq",
		"file.nt":      ".nt",
		"file.rdf":     ".rdf",
		"file.tar.gz":  ".gz", // only last extension is taken
		".hiddenfile":  "",    // edge case
		"noextension":  "",    // edge case
		"noextension.": ".",   // edge case
	}

	for input, expected := range tests {
		if ext := getExtension(input); ext != expected {
			t.Errorf("getExtension(%q) = %q, want %q", input, ext, expected)
		}
	}
}

func TestParseFile_NQuadsAndNTriples(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		wantLen  int
	}{
		{
			name:     "Valid .nq file",
			filePath: "test_data/test.nq",
			wantLen:  1,
		},
		{
			name:     "Valid .nt file",
			filePath: "test_data/test.nt",
			wantLen:  1,
		},
		{
			name:     "Valid .ttl file",
			filePath: "test_data/test.ttl",
			wantLen:  1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			quads, errChan := ParseFile(tc.filePath, Options{})

			var got []interfaces.IQuad
			for q := range quads {
				got = append(got, q)
			}

			err := <-errChan
			if err != nil {
				t.Fatalf("Unexpected error parsing %s: %v", tc.filePath, err)
			}
			if len(got) != tc.wantLen {
				t.Errorf("Expected %d quads, got %d", tc.wantLen, len(got))
			}
		})
	}
}

func TestParseFileUnsupportedExtension(t *testing.T) {
	_, errChan := ParseFile("test_data/test.rdfed", Options{})
	err := <-errChan
	if err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("Expected unsupported file format error, got: %v", err)
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, errChan := ParseFile("nonexistent.nq", Options{})
	err := <-errChan
	if err == nil || !os.IsNotExist(err) {
		t.Errorf("Expected file not found error, got: %v", err)
	}
}

func TestParseUnsupportedFormat(t *testing.T) {
	_, errChan := Parse(strings.NewReader(""), Options{Format: "application/xml"})
	err := <-errChan
	if err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("Expected unsupported format error, got: %v", err)
	}
}

func TestParseWithEmptyFormat(t *testing.T) {
	data := `<http://example.com/s> <http://example.com/p> <http://example.com/o> .`
	quads, errChan := Parse(strings.NewReader(data), Options{})
	var got []interfaces.IQuad
	for q := range quads {
		got = append(got, q)
	}
	if err := <-errChan; err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("Expected 1 quad, got %d", len(got))
	}
}

func TestParseWithNQuadsFormat(t *testing.T) {
	data := `<s> <p> <o> .`
	quads, errChan := Parse(strings.NewReader(data), Options{Format: "application/n-quads"})
	var got []interfaces.IQuad
	for q := range quads {
		got = append(got, q)
	}
	if err := <-errChan; err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("Expected 1 quad, got %d", len(got))
	}
}

func TestParseWithTurtleFormat(t *testing.T) {
	data := `@base <http://example.com/>. <s> <p> <o> .`
	quads, errChan := Parse(strings.NewReader(data), Options{Format: "text/turtle"})
	var got []interfaces.IQuad
	for q := range quads {
		got = append(got, q)
	}
	if err := <-errChan; err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("Expected 1 quad, got %d", len(got))
	}
}
