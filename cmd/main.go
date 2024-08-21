package main

import (
	"fmt"
	"rdfgo/interfaces"
	. "rdfgo/lib/data_model"
	. "rdfgo/lib/stream"
)

func main() {
	stream := make(interfaces.IStream, 1)
	store := NewStore()
	df := NewDataFactory()

	go func() {
		for i := 0; i < 10; i++ {
			quad, _ := NewQuad(
				df.NamedNode("http://example.com/s"),
				df.NamedNode("http://example.com/p"),
				df.NamedNode("http://example.com/o"),
				nil,
			)
			stream <- quad
		}
		close(stream)
	}()

	fmt.Println("importing stream")
	store.Import(stream)
	fmt.Println("Stream finished")
	fmt.Println(store.Size())
}
