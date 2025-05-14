# rdfgo

[![Coverage Status](https://coveralls.io/repos/github/maartyman/rdfgo/badge.svg)](https://coveralls.io/github/maartyman/rdfgo)
![CI status](https://github.com/maartyman/rdfgo/actions/workflows/CI.yml/badge.svg)

rdfgo proposes interfaces to work with rdf data based on the [rdfjs](https://rdf.js.org/) spec.
This can be found in the interfaces folder `github.com/maartyman/rdfgo/interfaces`.

Furthermore, rdfgo has implementations of most of these interfaces.
These implementations can be found in the lib folder `github.com/maartyman/rdfgo/lib`.

## Installation

```bash
go get github.com/maartyman/rdfgo
```

## Usage

### Data model

Terms can be created using the following methods:
```go
package main

import (
	"github.com/maartyman/rdfgo"
)

func main() {
	rdfgo.NewNamedNode("http://example.com/s")
	rdfgo.NewBlankNode("1")
	rdfgo.NewDefaultGraph()
	rdfgo.NewLiteral("string", "en", rdfgo.IRI.XSD.String)
	rdfgo.NewStringLiteral("string", "en")
	rdfgo.NewDecimalLiteral(0.1)
	rdfgo.NewBooleanLiteral(true)
	rdfgo.NewDoubleLiteral(0.1)
	rdfgo.NewIntegerLiteral(1)
	rdfgo.NewVariable("s")
    
    subject := rdfgo.NewNamedNode("http://example.com/s")
    predicate := rdfgo.NewNamedNode("http://example.com/p")
    object := rdfgo.NewStringLiteral("string", "en")
    quad, err := rdfgo.NewQuad(subject, predicate, object, nil)
    if err != nil {
        println(err)
    }
}
```

Terms have the following methods:
```go
namedNode := rdfgo.NewNamedNode("http://example.com/s")

namedNode.Equals(namedNode) // Check if two terms are equal
namedNode.GetType() // Get the term type, returns `TermType` enum
namedNode.GetValue() // Get the value of the term
namedNode.ToString() // Get the string representation of the term
```

TermType is an enum that has the following methods:
```go
termType := namedNode.GetType()
termType.String() // Get the string representation of the enum
termType.EnumIndex() // Get the index of the enum
```

Quad has the following methods to extract the subject, predicate, object and graph:
```go
quad := rdfgo.NewQuad(
    rdfgo.NewNamedNode("http://example.com/s"),
    rdfgo.NewNamedNode("http://example.com/p"),
    rdfgo.NewStringLiteral("string", "en"),
	nil
)
quad.GetSubject()
quad.GetPredicate()
quad.GetObject()
quad.GetGraph()
```

### Parser
The RDF parser can parse turtle, n-triples, and n-quads formats.
One can use the parser to parse a file or an io.reader.
```go
package main

import (
	"github.com/maartyman/rdfgo"
	"os"
	"strings"
)

func main() {
	// Parse a file
	quads, errChan := rdfgo.ParseFile("file.ttl", rdfgo.Options{
		// you can set the base IRI here
		BaseIRI: "http://example.com/",
	})
	...
	
	// Parse a reader (or string)
	data, err := os.ReadFile("file.ttl")
	if err != nil {
		panic(err)
	}
	reader := strings.NewReader(string(data))
	quads, errChan := rdfgo.Parse(reader, rdfgo.Options{
		// The format can be "turtle", "n-triples", "n-quads", "ntriples" or "nquads" (case-insensitive)
		Format: "turtle",
	})
	if err, ok := <-errChan; ok {
		panic(err)
	}
	
	// One can then for example import the quads into a store
	store := rdfgo.NewStore()
	store.Import(quads)
	
	// or count them
	rdfgo.Stream(quads).Count()
}
```

### Stream
The stream can be used to create a stream of quads and perform operations on them.
The stream is a channel of quads.
The Stream needs to be converted to an IStream interface for the store to import it.
```go
package main

import (
	"github.com/maartyman/rdfgo"
)

func main() {
	quad, _ := rdfgo.NewQuad(
		rdfgo.NewNamedNode("http://example.com/s"),
		rdfgo.NewNamedNode("http://example.com/p"),
		rdfgo.NewNamedNode("http://example.com/o"),
		nil,
	)

	stream := rdfgo.NewStream() // This will create a new Stream
	stream <- quad        // This will add a quad to the stream
	close(stream)         // This will close the stream
	stream.ToIStream()    // This will convert the stream to an IStream interface

	stream.Import(rdfgo.NewStream().ToIStream())           // This will import a stream to another stream
	stream = rdfgo.ArrayToStream([]rdfgo.IQuad{quad}) // This will import an array of quads to a stream

	stream.Count()   // This will return the amount of quads in the stream
	stream.ToStore() // This will return a store with the quads from the stream
	stream.ToArray() // This will return an array of quads from the stream
}
```

### Store
The store can be used to store quads and perform operations on them. 
Note that the store is implemented with set semantics, meaning that it will not store duplicate quads.
```go
package main

import (
	"github.com/maartyman/rdfgo"
)

func main() {
	store := rdfgo.NewStore() // This will create a new store

	s := rdfgo.NewNamedNode("http://example.com/s")
	p := rdfgo.NewNamedNode("http://example.com/p")
	o := rdfgo.NewNamedNode("http://example.com/o")
	addStream := rdfgo.NewStream()
	quad, _ := rdfgo.NewQuad(
		s,
		p,
		o,
		nil,
	)
	removeStream := rdfgo.ArrayToStream([]rdfgo.IQuad{
		quad,
	})
	go func() {
		for i := 0; i < 10; i++ {
			quad, _ := rdfgo.NewQuad(
				rdfgo.NewNamedNode("http://example.com/s"+string(rune(i))),
				rdfgo.NewNamedNode("http://example.com/p"+string(rune(i))),
				rdfgo.NewNamedNode("http://example.com/o"+string(rune(i))),
				nil,
			)
			addStream <- quad
		}
		close(addStream)
	}()
	store.Import(addStream.ToIStream())         // This will import all quads from the stream to the store
	store.Remove(removeStream.ToIStream())      // This will remove the quad from the store

	store.AddQuad(quad)                         // This will add the quad to the store
	store.AddQuadFromTerms(s, p, o, nil)        // This will create and add the quad to the store
	store.RemoveQuad(quad)                      // This will remove the quad from the store
	store.RemoveMatches(nil, nil, nil, nil)     // This will remove all quads from the store
	store.DeleteGraph(NewDefaultGraph())        // This will remove all quads from the default graph
	store.Match(nil, nil, nil, nil)             // This will return a stream with all quads in the store
	store.Has(quad)                             // This will return true if the quad is in the store
	store.Size()                                // This will return the number of quads in the store
	store.ForEach(func(quad interfaces.IQuad) { // This will loop over all quads in the store
		println(quad.ToString())
	})
}
```

## Future work
### package
- [ ] Improve tests
- [ ] Improve make RDF 1.2 compliant

### interfaces
- [ ] Add support for the Query rdfjs spec

### lib
- [ ] Add dataset support to the store
- [ ] Add json-ld parser
- [ ] Add n3 parser
- [ ] Add trig parser
- [ ] Add nquads writer
- [ ] Add turtle writer
- [ ] The parser should make use of the DataFactory, this way users could supply their own.

## Development
RDFgo has a makefile that can be used to run tests and build the package.
To start developing, first run the following command:
```bash
make setup-project
```
This will enable the git hooks and make sure the commit message follow the Conventional Commits format:
```
type(scope): description
```
Valid types: feat, fix, chore<br>
Example: `feat(parser): add ability to parse arrays` <br>
It will also run the tests and linter before committing.

For other commands see the makefile.
