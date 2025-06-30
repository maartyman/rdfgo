package nquads

import (
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
	"io"
)

// Options holds the options for writing.
type Options struct {
}

// WriteNQuads writes the given stream to the output stream in N-Quads format.
func WriteNQuads(stream interfaces.IStream, outputStream io.Writer, options Options) (string, error) {
	if outputStream == nil {
		result := ""
		for quad := range stream {
			result += fmt.Sprintf("%s .\n", quad.ToString())
		}
		return result, nil
	}
	for quad := range stream {
		_, err := outputStream.Write([]byte(fmt.Sprintf("%s .\n", quad.ToString())))
		if err != nil {
			return "", fmt.Errorf("error writing quad: %w", err)
		}
	}

	return "", nil
}
