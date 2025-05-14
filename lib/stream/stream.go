package rdfgo

import "github.com/maartyman/rdfgo/interfaces"

// Stream is an extension of the interfaces.IStream interface. It has various methods to manipulate the stream, like counting the quads, importing a stream, and converting to an array.
type Stream interfaces.IStream

// NewStream creates a new chanel of IQuad. It has some helper functions to combine streams, convert to an array, create a store, count the quads.
func NewStream(size ...int) Stream {
	if len(size) > 0 {
		return make(Stream, size[0])
	}
	return make(Stream)
}

// ToIStream converts a Stream to an IStream. This is useful for passing the stream to functions that expect an interfaces.IStream.
func (s Stream) ToIStream() interfaces.IStream {
	return interfaces.IStream(s)
}

// ToArray converts a Stream to an array of quads.
func (s Stream) ToArray() []interfaces.IQuad {
	var quadArray []interfaces.IQuad
	for quad := range s {
		quadArray = append(quadArray, quad)
	}
	return quadArray
}

// ToStore imports the Stream into a new Store and returns this store.
func (s Stream) ToStore() interfaces.IStore {
	store := NewStore()
	for quad := range s {
		store.AddQuad(quad)
	}
	return store
}

// Count returns the number of quads in the Stream.
func (s Stream) Count() int {
	count := 0
	for range s {
		count++
	}
	return count
}

// Import imports another stream into this one.
func (s Stream) Import(stream interfaces.IStream) {
	go func() {
		for quad := range stream {
			s <- quad
		}
		close(s)
	}()
}

// ArrayToStream converts an array of quads to a Stream. It accepts a slice of interfaces.IQuad and returns a Stream.
func ArrayToStream(quads []interfaces.IQuad) Stream {
	channel := make(interfaces.IStream)
	go func() {
		for _, quad := range quads {
			channel <- quad
		}
		close(channel)
	}()
	return Stream(channel)
}
