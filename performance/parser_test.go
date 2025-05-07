package performance

import (
	"github.com/maartyman/rdfgo"
	"github.com/maartyman/rdfgo/lib/parser/nquads"
	"github.com/maartyman/rdfgo/lib/parser/turtle"
	"os"
	"strings"
	"testing"
)

func BenchmarkParseTurtle(b *testing.B) {
	data, err := os.ReadFile("data/large.ttl")
	if err != nil {
		b.Error(err.Error())
	}

	b.Run("Parse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reader := strings.NewReader(string(data))
			quads, errChan := turtle.Parse(reader, turtle.Options{})
			for range quads {
			}
			if err, ok := <-errChan; ok {
				b.Error(err.Error())
			}
		}
	})

	b.Run("ParseAndStore", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reader := strings.NewReader(string(data))
			quads, errChan := turtle.Parse(reader, turtle.Options{})
			store := rdfgo.NewStore()
			store.Import(quads)
			if err, ok := <-errChan; ok {
				b.Error(err.Error())
			}
		}
	})
}

func BenchmarkParseNQuads(b *testing.B) {
	data, err := os.ReadFile("data/large.nq")
	if err != nil {
		b.Error(err.Error())
	}

	b.Run("Parse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reader := strings.NewReader(string(data))
			quads, errChan := nquads.Parse(reader, nquads.Options{})
			for range quads {
			}
			if err, ok := <-errChan; ok {
				b.Error(err.Error())
			}
		}
	})

	b.Run("ParseAndStore", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reader := strings.NewReader(string(data))
			quads, errChan := nquads.Parse(reader, nquads.Options{})
			store := rdfgo.NewStore()
			store.Import(quads)
			if err, ok := <-errChan; ok {
				b.Error(err.Error())
			}
		}
	})
}
