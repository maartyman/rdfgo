package performance

import (
	"github.com/maartyman/rdfgo/lib/parser/turtle_parser"
	. "github.com/maartyman/rdfgo/lib/stream"
	"os"
	"strings"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	data, err := os.ReadFile("data/large.ttl")
	if err != nil {
		b.Error(err.Error())
	}

	b.Run("Parse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reader := strings.NewReader(string(data))
			quads, errChan := turtle_parser.ParseTurtle(reader, nil)
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
			quads, errChan := turtle_parser.ParseTurtle(reader, nil)
			store := NewStore()
			store.Import(quads)
			if err, ok := <-errChan; ok {
				b.Error(err.Error())
			}
		}
	})
}
