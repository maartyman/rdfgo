package rdfgo

import (
	"bytes"
	datamodel "github.com/maartyman/rdfgo/lib/data_model"
	stream "github.com/maartyman/rdfgo/lib/stream"
	"testing"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		expectedOutput string
		expectedError  error
	}{
		{
			name:           "Valid N-Quads format",
			format:         "n-quads",
			expectedOutput: "<http://ex/s> <http://ex/p> <http://ex/o> .",
			expectedError:  nil,
		},
		{
			name:           "Valid N-Triples format",
			format:         "n-triples",
			expectedOutput: "<http://ex/s> <http://ex/p> <http://ex/o> .",
			expectedError:  nil,
		},
		{
			name:           "Invalid format defaults to N-Quads",
			format:         "invalid-format",
			expectedOutput: "<http://ex/s> <http://ex/p> <http://ex/o> .",
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a real stream with RDF quads
			store := stream.NewStore()
			store.AddQuadFromTerms(
				datamodel.NewNamedNode("http://ex/s"),
				datamodel.NewNamedNode("http://ex/p"),
				datamodel.NewNamedNode("http://ex/o"),
				nil,
			)
			matchStream := store.Match(nil, nil, nil, nil)

			outputStream := &bytes.Buffer{}
			options := WriterOptions{Format: tt.format}

			// Call the Write function
			_, err := Write(matchStream, outputStream, options)

			// Validate output and error
			if !bytes.Contains(outputStream.Bytes(), []byte(tt.expectedOutput)) {
				t.Errorf("Write() output = %v, want %v", outputStream.String(), tt.expectedOutput)
			}
			if err != tt.expectedError {
				t.Errorf("Write() error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}
