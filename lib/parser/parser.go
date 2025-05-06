package rdfgo

import (
	"errors"
	"github.com/maartyman/rdfgo/interfaces"
	"github.com/maartyman/rdfgo/lib/parser/nquads_parser"
	"github.com/maartyman/rdfgo/lib/parser/turtle_parser"
	"io"
	"os"
)

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

func ParseFile(fileName string) (chan interfaces.IQuad, chan error) {
	quads := make(chan interfaces.IQuad)
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
		return nquads_parser.ParseNQuads(file)
	case ".ttl":
		return turtle_parser.ParseTurtle(file, nil)
	default:
		go func() {
			errChan <- errors.New("unsupported file format")
			close(errChan)
			close(quads)
		}()
		return quads, errChan
	}
}

// Parse parses the input stream based on the provided MIME type. If no MIME type is provided, it tries all supported formats.
func Parse(stream io.Reader, mime string) (chan interfaces.IQuad, chan error) {
	if mime == "" {
		return nquads_parser.ParseNQuads(stream)
	}
	switch mime {
	case "application/n-quads", "application/n-triples":
		return nquads_parser.ParseNQuads(stream)
	case "text/turtle":
		return turtle_parser.ParseTurtle(stream, nil)
	default:
		errChan := make(chan error)
		emptyChannel := make(chan interfaces.IQuad)
		go func() {
			errChan <- errors.New("unsupported MIME type")
			close(emptyChannel)
			close(errChan)
		}()
		return emptyChannel, errChan
	}
}
