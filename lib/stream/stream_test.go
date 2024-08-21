package rdfgo

import (
	"rdfgo/interfaces"
	. "rdfgo/lib/data_model"
	"testing"
	"time"
)

func TestNewStream(t *testing.T) {
	stream := NewStream()
	if stream == nil {
		t.Error("Stream should not be nil")
	}
	if cap(stream) != 0 {
		t.Errorf("Buffer size should be 0, but got %d", cap(stream))
	}
	close(stream)

	stream = NewStream(10)
	if stream == nil {
		t.Error("Stream should not be nil")
	}
	if cap(stream) != 10 {
		t.Errorf("Buffer size should be 10, but got %d", cap(stream))
	}
	close(stream)

	stream = NewStream(0)
	if stream == nil {
		t.Error("Stream should not be nil")
	}
	if cap(stream) != 0 {
		t.Errorf("Buffer size should be 0, but got %d", cap(stream))
	}
	close(stream)

	stream = NewStream(0)
	done := make(chan bool)
	go func() {
		for i := 0; i < 1_000; i++ {
			stream <- nil
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Error("Stream did not behave as unbounded; sending blocked.")
	}

	close(stream)
}

func TestStream_ToArray(t *testing.T) {
	stream := NewStream()
	quad, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewDefaultGraph())
	go func() {
		stream <- quad
		close(stream)
	}()
	quadArray := stream.ToArray()
	if len(quadArray) != 1 {
		t.Errorf("Expected 1, got %d", len(quadArray))
	}
	if !quadArray[0].Equals(quad) {
		t.Errorf("Expected %v, got %v", quad, quadArray[0])
	}
}

func TestStream_ToStore(t *testing.T) {
	stream := NewStream()
	quad, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewDefaultGraph())
	go func() {
		stream <- quad
		close(stream)
	}()
	store := stream.ToStore()
	if Stream(store.Match(nil, nil, nil, nil)).Count() != 1 {
		t.Errorf("Expected true, got false")
	}
}

func TestStream_Count(t *testing.T) {
	stream := NewStream()
	quad, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewDefaultGraph())
	go func() {
		stream <- quad
		close(stream)
	}()
	if stream.Count() != 1 {
		t.Errorf("Expected 1, got %d", stream.Count())
	}
}

func TestStream_Import(t *testing.T) {
	stream := NewStream()
	quad, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewDefaultGraph())
	stream2 := NewStream()
	go func() {
		stream2 <- quad
		close(stream2)
	}()
	stream.Import(interfaces.IStream(stream2))
	if stream.Count() != 1 {
		t.Errorf("Expected 1, got %d", stream.Count())
	}
}

func TestNewArrayStreamWithNoElements(t *testing.T) {
	quads := []interfaces.IQuad{}
	stream := ArrayToStream(quads)
	if stream == nil {
		t.Error("Stream should not be nil")
	}
	quad, ok := <-stream
	if ok {
		t.Error("Stream should be closed")
	}
	if quad != nil {
		t.Error("Quad should be nil")
	}
}

func TestNewArrayStreamWithOneElement(t *testing.T) {
	quad, _ := NewQuad(
		NewNamedNode("http://example.com/s"),
		NewNamedNode("http://example.com/p"),
		NewNamedNode("http://example.com/o"),
		NewDefaultGraph(),
	)
	quads := []interfaces.IQuad{
		quad,
	}
	stream := ArrayToStream(quads)
	if stream == nil {
		t.Error("Stream should not be nil")
	}
	count := 0
	for quad := range stream {
		if quad == nil {
			t.Error("Quad should not be nil")
		} else {
			count++
		}
	}
	if count != 1 {
		t.Error("Stream should have only one quad but received", count)
	}
}

func TestNewArrayStreamWithOneNilElement(t *testing.T) {
	quads := []interfaces.IQuad{
		nil,
	}
	stream := ArrayToStream(quads)
	if stream == nil {
		t.Error("Stream should not be nil")
	}
	count := 0
	for quad := range stream {
		if quad != nil {
			t.Error("Quad should be nil")
			count++
		}
	}
	if count != 0 {
		t.Error("Stream should have only zero quad but received", count)
	}
}
