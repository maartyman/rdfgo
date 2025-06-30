package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
	"github.com/maartyman/rdfgo/lib/writer/nquads"
	"github.com/maartyman/rdfgo/lib/writer/ntriples"
	"io"
)

// WriterOptions is used to configure the parser. It includes options for the base IRI, and the format type.
type WriterOptions struct {
	Format   string
	Prefixes map[string]string
}

func Write(stream interfaces.IStream, outputStream io.Writer, options WriterOptions) (string, error) {
	switch options.Format {
	case "application/n-quads", "n-quads", "nquads":
		return nquads.WriteNQuads(stream, outputStream, nquads.Options{})
	case "application/n-triples", "n-triples", "ntriples":
		return ntriples.WriteNTriples(stream, outputStream, ntriples.Options{})
	default:
		return nquads.WriteNQuads(stream, outputStream, nquads.Options{})
	}
}
