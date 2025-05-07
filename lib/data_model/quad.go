package rdfgo

import (
	"errors"
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
)

const QuadValue = ""

var SubjectTermTypeError = errors.New("subject needs to be a namedNode, blankNode, quad or variable")
var PredicateTermTypeError = errors.New("predicate needs to be a namedNode or variable")
var ObjectTermTypeError = errors.New("object needs to be a namedNode, blankNode, literal, or variable")
var GraphTermTypeError = errors.New("graph needs to be a namedNode, blankNode, defaultGraph, or variable")

type quad struct {
	subject   interfaces.ITerm
	predicate interfaces.ITerm
	object    interfaces.ITerm
	graph     interfaces.ITerm
}

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

func (q *quad) GetType() interfaces.TermType {
	return interfaces.QuadType
}

func (q *quad) GetValue() string {
	return QuadValue
}

func (q *quad) GetSubject() interfaces.ITerm {
	return q.subject
}

func (q *quad) GetPredicate() interfaces.ITerm {
	return q.predicate
}

func (q *quad) GetObject() interfaces.ITerm {
	return q.object
}

func (q *quad) GetGraph() interfaces.ITerm {
	return q.graph
}

func (q *quad) ToString() string {
	return fmt.Sprintf(
		"%s %s %s %s",
		q.subject.ToString(),
		q.predicate.ToString(),
		q.object.ToString(),
		q.graph.ToString())
}
