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
	. "github.com/maartyman/rdfgo/lib/data_model"
)

func main() {
    NewNamedNode("http://example.com/s")
    NewBlankNode("1")
    NewDefaultGraph()
    NewLiteral("string", "en", IRI.XSD.String)
    NewStringLiteral("string", "en")
    NewDecimalLiteral(0.1)
    NewBooleanLiteral(true)
    NewDoubleLiteral(0.1)
    NewIntegerLiteral(1)
    NewVariable("s")
    
    subject := NewNamedNode("http://example.com/s")
    predicate := NewNamedNode("http://example.com/p")
    object := NewStringLiteral("string", "en")
    quad, err := NewQuad(subject, predicate, object, nil)
    if err != nil {
        println(err)
    }
}
```

Terms have the following methods:
```go
namedNode := NewNamedNode("http://example.com/s")

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
quad := NewQuad(
	NewNamedNode("http://example.com/s"),
	NewNamedNode("http://example.com/p"),
	NewStringLiteral("string", "en"),
	nil
)
quad.GetSubject()
quad.GetPredicate()
quad.GetObject()
quad.GetGraph()
```

### Stream
The stream can be used to create a stream of quads and perform operations on them.
The stream is a channel of quads.
The Stream needs to be converted to an IStream interface for the store to import it.
```go
package main

import (
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	. "github.com/maartyman/rdfgo/lib/stream"
)

func main() {
	quad, _ := NewQuad(
		NewNamedNode("http://example.com/s"),
		NewNamedNode("http://example.com/p"),
		NewNamedNode("http://example.com/o"),
		nil,
	)

	stream := NewStream() // This will create a new Stream
	stream <- quad        // This will add a quad to the stream
	close(stream)         // This will close the stream
	stream.ToIStream()    // This will convert the stream to an IStream interface

	stream.Import(NewStream().ToIStream())           // This will import a stream to another stream
	stream = ArrayToStream([]interfaces.IQuad{quad}) // This will import an array of quads to a stream

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
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	. "github.com/maartyman/rdfgo/lib/stream"
)

func main() {
	store := NewStore() // This will create a new store

	s := NewNamedNode("http://example.com/s")
	p := NewNamedNode("http://example.com/p")
	o := NewNamedNode("http://example.com/o")
	addStream := NewStream()
	quad, _ := NewQuad(
		s,
		p,
		o,
		nil,
	)
	removeStream := ArrayToStream([]interfaces.IQuad{
		quad,
	})
	go func() {
		for i := 0; i < 10; i++ {
			quad, _ := NewQuad(
				NewNamedNode("http://example.com/s"+string(rune(i))),
				NewNamedNode("http://example.com/p"+string(rune(i))),
				NewNamedNode("http://example.com/o"+string(rune(i))),
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
- [ ] Add CI/CD to the package

### interfaces
- [ ] Add support for the Query rdfjs spec

### lib
- [ ] Add dataset support to the store
- [ ] Add a parser to the lib portion of the package

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
Valid types: feat, fix, chore, docs, style, refactor, test, perf, ci <br>
Example: `feat(parser): add ability to parse arrays` <br>
It will also run the tests and linter before committing.

For other commands see the makefile.
