package rdfgo

import "rdfgo/interfaces"

type Stream interfaces.IStream

func NewStream(size ...int) Stream {
	if len(size) > 0 {
		return make(Stream, size[0])
	}
	return make(Stream)
}

func (s Stream) ToArray() []interfaces.IQuad {
	var quadArray []interfaces.IQuad
	for quad := range s {
		quadArray = append(quadArray, quad)
	}
	return quadArray
}

func (s Stream) ToStore() interfaces.IStore {
	store := NewStore()
	for quad := range s {
		store.AddQuad(quad)
	}
	return store
}

func (s Stream) Count() int {
	count := 0
	for range s {
		count++
	}
	return count
}

func (s Stream) Import(stream interfaces.IStream) {
	go func() {
		for quad := range stream {
			s <- quad
		}
		close(s)
	}()
}

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
