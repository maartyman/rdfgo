package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"sync"
)

type store struct {
	size    int
	entries map[string][]interfaces.IQuad
	mux     sync.RWMutex
}

// Store is an extension of the interfaces.IStore interface. It has various methods to manipulate the store, like adding and removing quads, checking if a quad exists, and iterating over the quads.
type Store interface {
	interfaces.IStore
	Size() int
	Has(interfaces.IQuad) bool
	AddQuadFromTerms(interfaces.ITerm, interfaces.ITerm, interfaces.ITerm, interfaces.ITerm) bool
	AddQuad(interfaces.IQuad) bool
	RemoveQuad(interfaces.IQuad)
	ForEach(func(interfaces.IQuad))
}

// NewStore creates a new store of quads. A store indexes the quads by their subject, predicate, object, and graph. Note that this store is set semantics, meaning that it does not allow duplicate quads.
func NewStore() Store {
	return &store{
		size:    0,
		entries: make(map[string][]interfaces.IQuad),
	}
}

// Size returns the number of quads in the store.
func (s *store) Size() int {
	return s.size
}

func getHashes(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) []string {
	//TODO change to multiple return values
	if graph.GetType() == interfaces.DefaultGraphType {
		return []string{
			subject.ToString() + ",,,",
			"," + predicate.ToString() + ",,",
			",," + object.ToString() + ",",
			",,," + DefaultGraphValue,
			subject.ToString() + "," + predicate.ToString() + "," + object.ToString() + "," + DefaultGraphValue,
		}
	}
	return []string{
		subject.ToString() + ",,,",
		"," + predicate.ToString() + ",,",
		",," + object.ToString() + ",",
		",,," + graph.ToString(),
		subject.ToString() + "," + predicate.ToString() + "," + object.ToString() + "," + graph.ToString(),
	}
}

func convertVariablesToNil(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) (interfaces.ITerm, interfaces.ITerm, interfaces.ITerm, interfaces.ITerm) {
	if subject != nil && subject.GetType() == interfaces.VariableType {
		subject = nil
	}
	if predicate != nil && predicate.GetType() == interfaces.VariableType {
		predicate = nil
	}
	if object != nil && object.GetType() == interfaces.VariableType {
		object = nil
	}
	if graph != nil && graph.GetType() == interfaces.VariableType {
		graph = nil
	}
	return subject, predicate, object, graph
}

// Has checks if the store contains a quad. It returns true if the quad exists, false otherwise.
func (s *store) Has(quad interfaces.IQuad) bool {
	s.mux.Lock()
	_, exists := s.entries[quad.GetSubject().ToString()+","+quad.GetPredicate().ToString()+","+quad.GetObject().ToString()+","+quad.GetGraph().ToString()]
	s.mux.Unlock()
	return exists
}

// AddQuadFromTerms adds a quad to the store. It takes four terms: subject, predicate, object, and graph. If any of the terms are nil, it returns false. If the quad already exists in the store, it returns false.
func (s *store) AddQuadFromTerms(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) bool {
	if subject == nil || predicate == nil || object == nil {
		return false
	}
	if graph == nil {
		graph = NewDefaultGraph()
	}
	if subject.GetType() == interfaces.VariableType ||
		predicate.GetType() == interfaces.VariableType ||
		object.GetType() == interfaces.VariableType ||
		graph.GetType() == interfaces.VariableType {
		return false
	}

	quad, err := NewQuad(subject, predicate, object, graph)
	if err != nil {
		return false
	}
	if s.Has(quad) {
		return false
	}
	s.mux.Lock()
	for _, hash := range getHashes(subject, predicate, object, graph) {
		quadArray, exists := s.entries[hash]
		if !exists {
			s.entries[hash] = []interfaces.IQuad{quad}
		} else {
			s.entries[hash] = append(quadArray, quad)
		}
	}

	s.size++
	s.mux.Unlock()
	return true
}

// AddQuad adds a quad to the store. It takes a quad as an argument. If the quad already exists in the store, it returns false.
func (s *store) AddQuad(quad interfaces.IQuad) bool {
	return s.AddQuadFromTerms(quad.GetSubject(), quad.GetPredicate(), quad.GetObject(), quad.GetGraph())
}

// RemoveQuad removes a quad from the store. It takes a quad as an argument. If the quad does not exist in the store, it does nothing.
func (s *store) RemoveQuad(quad interfaces.IQuad) {
	if !s.Has(quad) {
		return
	}

	for _, hash := range getHashes(quad.GetSubject(), quad.GetPredicate(), quad.GetObject(), quad.GetGraph()) {
		s.mux.Lock()
		quadArray := s.entries[hash]
		for i := 0; i < len(quadArray); i++ {
			q := quadArray[i]
			if q.GetSubject().Equals(quad.GetSubject()) &&
				q.GetPredicate().Equals(quad.GetPredicate()) &&
				q.GetObject().Equals(quad.GetObject()) &&
				q.GetGraph().Equals(quad.GetGraph()) {
				if len(quadArray) == 1 {
					delete(s.entries, hash)
					break
				}
				s.entries[hash] = append(quadArray[:i], quadArray[i+1:]...)
				break
			}
		}
		s.mux.Unlock()
	}

	s.mux.Lock()
	s.size--
	s.mux.Unlock()
}

// RemoveMatches removes all quads that match the given subject, predicate, object, and graph. It takes four terms as arguments.
func (s *store) RemoveMatches(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) {
	for quad := range s.Match(subject, predicate, object, graph) {
		s.RemoveQuad(quad)
	}
}

// Remove removes all quads from the store that are in the given stream. It takes a stream as an argument.
func (s *store) Remove(stream interfaces.IStream) {
	for quad := range stream {
		if quad != nil {
			s.RemoveQuad(quad)
		}
	}
}

// DeleteGraph removes all quads from the store that are in the given graph. It takes a graph as an argument.
func (s *store) DeleteGraph(graph interfaces.ITerm) {
	s.RemoveMatches(nil, nil, nil, graph)
}

func (s *store) matchSubject(subject interfaces.ITerm) []interfaces.IQuad {
	return s.entries[subject.ToString()+",,,"]
}

func (s *store) matchPredicate(predicate interfaces.ITerm) []interfaces.IQuad {
	return s.entries[","+predicate.ToString()+",,"]
}

func (s *store) matchObject(object interfaces.ITerm) []interfaces.IQuad {
	return s.entries[",,"+object.ToString()+","]
}

func (s *store) matchGraph(graph interfaces.ITerm) []interfaces.IQuad {
	if graph.GetType() == interfaces.DefaultGraphType {
		return s.entries[",,,"+DefaultGraphValue]
	}
	return s.entries[",,,"+graph.ToString()]
}

// Match returns a stream of quads that match the given subject, predicate, object, and graph. It takes four terms as arguments. If all terms are nil, it returns all quads in the store.
func (s *store) Match(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) interfaces.IStream {
	subject, predicate, object, graph = convertVariablesToNil(subject, predicate, object, graph)
	quadStream := make(interfaces.IStream, 10)

	if subject == nil && predicate == nil && object == nil && graph == nil {
		go func() {
			s.mux.Lock()
			for key, value := range s.entries {
				s.mux.Unlock()
				if key[0] != ',' && key[len(key)-1] != ',' {
					for _, quad := range value {
						quadStream <- quad
					}
				}
				s.mux.Lock()
			}
			s.mux.Unlock()
			close(quadStream)
		}()
	} else if subject != nil && predicate != nil && object != nil && graph != nil {
		quad, err := NewQuad(subject, predicate, object, graph)
		if err != nil {
			close(quadStream)
			return quadStream
		}
		if s.Has(quad) {
			quadStream <- quad
		}
		close(quadStream)
	} else {
		smallest := [2]int{0, 0} // [index, size]
		var subjectMatches []interfaces.IQuad = nil
		if subject != nil {
			subjectMatches = s.matchSubject(subject)
			if subjectMatches != nil {
				smallest[1] = len(subjectMatches)
			}
		}
		var predicateMatches []interfaces.IQuad = nil
		if predicate != nil {
			predicateMatches = s.matchPredicate(predicate)
			if predicateMatches != nil && (smallest[1] == 0 || len(predicateMatches) < smallest[1]) {
				smallest[0] = 1
				smallest[1] = len(predicateMatches)
			}
		}
		var objectMatches []interfaces.IQuad = nil
		if object != nil {
			objectMatches = s.matchObject(object)
			if objectMatches != nil && (smallest[1] == 0 || len(objectMatches) < smallest[1]) {
				smallest[0] = 2
				smallest[1] = len(objectMatches)
			}
		}
		var graphMatches []interfaces.IQuad = nil
		if graph != nil {
			graphMatches = s.matchGraph(graph)
			if graphMatches != nil && (smallest[1] == 0 || len(graphMatches) < smallest[1]) {
				smallest[0] = 3
				smallest[1] = len(graphMatches)
			}
		}

		switch smallest[0] {
		case 0:
			go func() {
				for _, quad := range subjectMatches {
					if (predicate == nil || quad.GetPredicate().Equals(predicate)) &&
						(object == nil || quad.GetObject().Equals(object)) &&
						(graph == nil || quad.GetGraph().Equals(graph)) {
						quadStream <- quad
					}
				}
				close(quadStream)
			}()
		case 1:
			go func() {
				for _, quad := range predicateMatches {
					if (subject == nil || quad.GetSubject().Equals(subject)) &&
						(object == nil || quad.GetObject().Equals(object)) &&
						(graph == nil || quad.GetGraph().Equals(graph)) {
						quadStream <- quad
					}
				}
				close(quadStream)
			}()
		case 2:
			go func() {
				for _, quad := range objectMatches {
					if (subject == nil || quad.GetSubject().Equals(subject)) &&
						(predicate == nil || quad.GetPredicate().Equals(predicate)) &&
						(graph == nil || quad.GetGraph().Equals(graph)) {
						quadStream <- quad
					}
				}
				close(quadStream)
			}()
		case 3:
			go func() {
				for _, quad := range graphMatches {
					if (subject == nil || quad.GetSubject().Equals(subject)) &&
						(predicate == nil || quad.GetPredicate().Equals(predicate)) &&
						(object == nil || quad.GetObject().Equals(object)) {
						quadStream <- quad
					}
				}
				close(quadStream)
			}()
		}
	}
	return quadStream
}

// Import imports a stream of quads into the store. It takes a stream as an argument.
func (s *store) Import(quadStream interfaces.IStream) {
	for quad := range quadStream {
		if quad != nil {
			s.AddQuad(quad)
		}
	}
}

// ForEach iterates over all quads in the store and applies the given callback function to each quad. It takes a callback function as an argument.
func (s *store) ForEach(callback func(interfaces.IQuad)) {
	for key, value := range s.entries {
		if key[0] != ',' && key[len(key)-1] != ',' {
			for _, quad := range value {
				callback(quad)
			}
		}
	}
}
