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

// Options holds the options for parsing.
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

// ParseFile parses the input file based on its extension. It supports .nt, .nq, and .ttl formats.
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

// Parse parses the input stream based on the provided format type. If no format type is provided, it will assume turtle.
func Parse(stream io.Reader, options Options) (interfaces.IStream, chan error) {
	if options.Format == "" {
		data, errChan := turtle.Parse(stream, turtle.Options{BaseIRI: options.BaseIRI})
		newErrChan := make(chan error, 1)
		go func() {
			println("Parsing turtle format, waiting if parsering error")
			if err := <-errChan; err != nil {
				println("Error parsing turtle format")
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
