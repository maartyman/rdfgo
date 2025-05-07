package rdfgo

import (
	"errors"
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
)

// QuadValue is a constant that represents the quad value. This is never used as the value of a quad is the combination of the subject, predicate, object, and graph.
const QuadValue = ""

// SubjectTermTypeError is an error type that is thrown when a quad's subject is not of the correct type.
var SubjectTermTypeError = errors.New("subject needs to be a namedNode, blankNode, quad or variable")

// PredicateTermTypeError is an error type that is thrown when a quad's predicate is not of the correct type.
var PredicateTermTypeError = errors.New("predicate needs to be a namedNode or variable")

// ObjectTermTypeError is an error type that is thrown when a quad's object is not of the correct type.
var ObjectTermTypeError = errors.New("object needs to be a namedNode, blankNode, literal, or variable")

// GraphTermTypeError is an error type that is thrown when a quad's graph is not of the correct type.
var GraphTermTypeError = errors.New("graph needs to be a namedNode, blankNode, defaultGraph, or variable")

type quad struct {
	subject   interfaces.ITerm
	predicate interfaces.ITerm
	object    interfaces.ITerm
	graph     interfaces.ITerm
}

// NewQuad is a constructor for creating a quad. This returns an implementation of the interfaces.IQuad interface.
func NewQuad(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) (interfaces.IQuad, error) {
	if subject == nil || (subject.GetType() != interfaces.NamedNodeType &&
		subject.GetType() != interfaces.BlankNodeType &&
		subject.GetType() != interfaces.QuadType &&
		subject.GetType() != interfaces.VariableType) {
		return nil, SubjectTermTypeError
	}
	if predicate == nil || (predicate.GetType() != interfaces.NamedNodeType &&
		predicate.GetType() != interfaces.VariableType) {
		return nil, PredicateTermTypeError
	}
	if object == nil || (object.GetType() != interfaces.NamedNodeType &&
		object.GetType() != interfaces.BlankNodeType &&
		object.GetType() != interfaces.QuadType &&
		object.GetType() != interfaces.LiteralType &&
		object.GetType() != interfaces.VariableType) {
		return nil, ObjectTermTypeError
	}
	if graph == nil {
		graph = NewDefaultGraph()
	} else if graph.GetType() != interfaces.NamedNodeType &&
		graph.GetType() != interfaces.BlankNodeType &&
		graph.GetType() != interfaces.DefaultGraphType &&
		graph.GetType() != interfaces.VariableType {
		return nil, GraphTermTypeError
	}
	return &quad{
		subject:   subject,
		predicate: predicate,
		object:    object,
		graph:     graph,
	}, nil
}

// Equals is a method that checks if two quads are equal. This returns true if the other term is also a quad and has the same subject, predicate, object, and graph.
func (q *quad) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	if q == other {
		return true
	}
	quad, ok := other.(*quad)
	if !ok || interfaces.QuadType != other.GetType() {
		return false
	}
	return QuadValue == quad.GetValue() &&
		q.subject.Equals(quad.GetSubject()) &&
		q.object.Equals(quad.GetObject()) &&
		q.predicate.Equals(quad.GetPredicate()) &&
		q.graph.Equals(quad.GetGraph())
}

// GetType is a method that returns the (integer) type of the quad.
func (q *quad) GetType() interfaces.TermType {
	return interfaces.QuadType
}

// GetValue is a method that returns the value of the quad. This is never used as the value of a quad is the combination of the subject, predicate, object, and graph.
func (q *quad) GetValue() string {
	return QuadValue
}

// GetSubject is a method that returns the subject of the quad.
func (q *quad) GetSubject() interfaces.ITerm {
	return q.subject
}

// GetPredicate is a method that returns the predicate of the quad.
func (q *quad) GetPredicate() interfaces.ITerm {
	return q.predicate
}

// GetObject is a method that returns the object of the quad.
func (q *quad) GetObject() interfaces.ITerm {
	return q.object
}

// GetGraph is a method that returns the graph of the quad.
func (q *quad) GetGraph() interfaces.ITerm {
	return q.graph
}

// ToString is a method that returns the string representation of the quad (<s> <p> <o> <g>).
func (q *quad) ToString() string {
	return fmt.Sprintf(
		"%s %s %s %s",
		q.subject.ToString(),
		q.predicate.ToString(),
		q.object.ToString(),
		q.graph.ToString())
}
