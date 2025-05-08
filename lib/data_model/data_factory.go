package rdfgo

import (
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
)

type dataFactory struct {
	blankNodeCounter int
}

// DataFactory is an extension of the interfaces.IDataFactory interface. It includes a SimpleLiteral method that creates a literal with a default datatype of xsd:string.
type DataFactory interface {
	interfaces.IDataFactory
	SimpleLiteral(value string) interfaces.ILiteral
}

// NewDataFactory is a constructor for creating a data factory. This returns an implementation of the DataFactory interface.
func NewDataFactory() DataFactory {
	return &dataFactory{0}
}

// NamedNode is a method that creates a named node. This returns an implementation of the interfaces.INamedNode interface.
func (df *dataFactory) NamedNode(value string) interfaces.INamedNode {
	return NewNamedNode(value)
}

// BlankNode is a method that creates a blank node. This returns an implementation of the interfaces.IBlankNode interface. if the value is empty, it generates a new blank node with an incremented counter (b0, b1, b2, etc.).
func (df *dataFactory) BlankNode(value string) interfaces.IBlankNode {
	if value == "" {
		value = fmt.Sprintf("b%d", df.blankNodeCounter)
		df.blankNodeCounter++
	}
	return NewBlankNode(value)
}

// SimpleLiteral is a method that creates a literal with a default datatype of xsd:string. This returns an implementation of the interfaces.ILiteral interface.
func (df *dataFactory) SimpleLiteral(value string) interfaces.ILiteral {
	return NewLiteral(value, "", IRI.XSD.String)
}

// Literal is a method that creates a literal. This returns an implementation of the interfaces.ILiteral interface. It accepts a value, language, and datatype.
func (df *dataFactory) Literal(value string, language string, datatype interfaces.INamedNode) interfaces.ILiteral {
	if datatype == nil {
		datatype = IRI.XSD.String
	}
	return NewLiteral(value, language, datatype)
}

// Variable is a method that creates a variable. This returns an implementation of the interfaces.IVariable interface.
func (df *dataFactory) Variable(value string) interfaces.IVariable {
	return NewVariable(value)
}

// DefaultGraph is a method that creates a default graph. This returns an implementation of the interfaces.IDefaultGraph interface.
func (df *dataFactory) DefaultGraph() interfaces.IDefaultGraph {
	return NewDefaultGraph()
}

// Quad is a method that creates a quad. This returns an implementation of the interfaces.IQuad interface.
func (df *dataFactory) Quad(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) (interfaces.IQuad, error) {
	return NewQuad(subject, predicate, object, graph)
}
