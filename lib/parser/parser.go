package rdfgo

import (
	"errors"
	"github.com/maartyman/rdfgo/interfaces"
	"github.com/maartyman/rdfgo/lib/parser/nquads"
	"github.com/maartyman/rdfgo/lib/parser/turtle"
	"io"
	"os"
	"strings"
)

// Options is used to configure the parser. It includes options for the base IRI, and the format type.
type Options struct {
	Format  string
	BaseIRI string
}

func getExtension(fileName string) string {
	extensionStart := -1
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			extensionStart = i
			break
		}
	}
	if extensionStart == -1 || extensionStart == 0 {
		return ""
	}
	return fileName[extensionStart:]
}

// ParseFile is a function that parses a file into a stream of quads. It accepts a file path and parser.Options as parameters, and returns a channel of quads and a channel of errors. It supports n-triples, n-quads, and turtle formats (.nt, .nq, .ttl).
func ParseFile(fileName string, options Options) (interfaces.IStream, chan error) {
	quads := make(interfaces.IStream)
	errChan := make(chan error, 1)

	file, err := os.Open(fileName)
	if err != nil {
		go func() {
			errChan <- err
			close(errChan)
			close(quads)
		}()
		return quads, errChan
	}
	switch getExtension(fileName) {
	case ".nt", ".nq":
		return nquads.Parse(file, nquads.Options{})
	case ".ttl":
		return turtle.Parse(file, turtle.Options{BaseIRI: options.BaseIRI})
	default:
		return Parse(file, options)
	}
}

// Parse is a function that parses a string into a stream of quads. It accepts an io.Reader and parser.Options as parameters, and returns a channel of quads and a channel of errors. If no format is specified, it defaults to turtle format.
func Parse(stream io.Reader, options Options) (interfaces.IStream, chan error) {
	if options.Format == "" {
		data, errChan := turtle.Parse(stream, turtle.Options{BaseIRI: options.BaseIRI})
		newErrChan := make(chan error, 1)
		go func() {
			if err := <-errChan; err != nil {
				newErrChan <- errors.New("unsupported format in options")
			}
			close(newErrChan)
		}()
		return data, newErrChan
	}
	options.Format = strings.ToLower(options.Format)
	switch options.Format {
	case "application/n-quads", "application/n-triples", "n-triples", "n-quads", "ntriples", "nquads":
		return nquads.Parse(stream, nquads.Options{})
	case "text/turtle", "turtle":
		return turtle.Parse(stream, turtle.Options{BaseIRI: options.BaseIRI})
	default:
		errChan := make(chan error)
		emptyChannel := make(interfaces.IStream)
		go func() {
			errChan <- errors.New("unsupported format in options")
			close(emptyChannel)
			close(errChan)
		}()
		return emptyChannel, errChan
	}
}
