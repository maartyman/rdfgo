package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"sync"
)

type Store struct {
	size    int
	entries map[string][]interfaces.IQuad
	mux     sync.RWMutex
}

type IStore interface {
	interfaces.IStore
	Size() int
	Has(interfaces.IQuad) bool
	AddQuadFromTerms(interfaces.ITerm, interfaces.ITerm, interfaces.ITerm, interfaces.ITerm) bool
	AddQuad(interfaces.IQuad) bool
	RemoveQuad(interfaces.IQuad)
	ForEach(func(interfaces.IQuad))
}

func NewStore() IStore {
	return &Store{
		size:    0,
		entries: make(map[string][]interfaces.IQuad),
	}
}

func (s *Store) Size() int {
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

func (s *Store) Has(quad interfaces.IQuad) bool {
	s.mux.Lock()
	_, exists := s.entries[quad.GetSubject().ToString()+","+quad.GetPredicate().ToString()+","+quad.GetObject().ToString()+","+quad.GetGraph().ToString()]
	s.mux.Unlock()
	return exists
}

func (s *Store) AddQuadFromTerms(
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

func (s *Store) AddQuad(quad interfaces.IQuad) bool {
	return s.AddQuadFromTerms(quad.GetSubject(), quad.GetPredicate(), quad.GetObject(), quad.GetGraph())
}

func (s *Store) RemoveQuad(quad interfaces.IQuad) {
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

func (s *Store) RemoveMatches(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) {
	for quad := range s.Match(subject, predicate, object, graph) {
		s.RemoveQuad(quad)
	}
}

func (s *Store) Remove(stream interfaces.IStream) {
	for quad := range stream {
		if quad != nil {
			s.RemoveQuad(quad)
		}
	}
}

func (s *Store) DeleteGraph(graph interfaces.ITerm) {
	s.RemoveMatches(nil, nil, nil, graph)
}

func (s *Store) matchSubject(subject interfaces.ITerm) []interfaces.IQuad {
	return s.entries[subject.ToString()+",,,"]
}

func (s *Store) matchPredicate(predicate interfaces.ITerm) []interfaces.IQuad {
	return s.entries[","+predicate.ToString()+",,"]
}

func (s *Store) matchObject(object interfaces.ITerm) []interfaces.IQuad {
	return s.entries[",,"+object.ToString()+","]
}

func (s *Store) matchGraph(graph interfaces.ITerm) []interfaces.IQuad {
	if graph.GetType() == interfaces.DefaultGraphType {
		return s.entries[",,,"+DefaultGraphValue]
	}
	return s.entries[",,,"+graph.ToString()]
}

func (s *Store) Match(
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

func (s *Store) Import(quadStream interfaces.IStream) {
	for quad := range quadStream {
		if quad != nil {
			s.AddQuad(quad)
		}
	}
}

func (s *Store) ForEach(callback func(interfaces.IQuad)) {
	for key, value := range s.entries {
		if key[0] != ',' && key[len(key)-1] != ',' {
			for _, quad := range value {
				callback(quad)
			}
		}
	}
}
