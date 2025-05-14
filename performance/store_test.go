package performance

import (
	"github.com/maartyman/rdfgo"
	"os"
	"strings"
	"testing"
)

func BenchmarkStore(b *testing.B) {
	data, err := os.ReadFile("data/large.ttl")
	if err != nil {
		b.Error(err.Error())
	}
	reader := strings.NewReader(string(data))
	quads, errChan := rdfgo.Parse(reader, rdfgo.Options{
		Format:  "text/turtle",
		BaseIRI: "http://example.com/",
	})
	store := rdfgo.NewStore()
	store.Import(quads)
	if err, ok := <-errChan; ok {
		b.Error(err.Error())
	}

	b.Run("Match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stream := store.Match(nil, nil, rdfgo.NewNamedNode("http://www.semanticweb.org/ontologies/2015/trainbenchmark#TrackElement"), nil)
			for range stream {
			}
		}
	})
}
